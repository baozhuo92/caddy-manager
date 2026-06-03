package services

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"syscall"
	"time"

	"caddy_server/database"
)

var (
	CaddyCmd     *exec.Cmd
	CaddyRunning bool
	CaddyMu      sync.Mutex
	logChan      = make(chan database.CaddyLog, 100)
	LogStopChan  = make(chan struct{})
)

func StartCaddy(binaryPath, caddyfilePath string) error {
	CaddyMu.Lock()
	defer CaddyMu.Unlock()

	if CaddyRunning {
		return fmt.Errorf("Caddy 已在运行中")
	}

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("Caddy 二进制文件不存在: %s", binaryPath)
	}

	CaddyCmd = exec.Command(binaryPath, "run", "--config", caddyfilePath)
	CaddyCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := CaddyCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("无法获取 stdout: %w", err)
	}
	stderr, err := CaddyCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("无法获取 stderr: %w", err)
	}

	if err := CaddyCmd.Start(); err != nil {
		return fmt.Errorf("启动 Caddy 失败: %w", err)
	}

	CaddyRunning = true

	go collectLogs(stdout, "stdout")
	go collectLogs(stderr, "stderr")
	go batchWriteLogs()

	go func() {
		CaddyCmd.Wait()
		CaddyMu.Lock()
		CaddyRunning = false
		CaddyMu.Unlock()
		close(LogStopChan)
	}()

	return nil
}

func StopCaddy() error {
	CaddyMu.Lock()
	defer CaddyMu.Unlock()

	if !CaddyRunning || CaddyCmd == nil {
		return fmt.Errorf("Caddy 未在运行")
	}

	if err := CaddyCmd.Process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("发送停止信号失败: %w", err)
	}

	done := make(chan struct{})
	go func() {
		CaddyCmd.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		CaddyCmd.Process.Kill()
	}

	CaddyRunning = false
	return nil
}

func ReloadCaddy(binaryPath, caddyfilePath string) error {
	if !CaddyRunning {
		return fmt.Errorf("Caddy 未运行，无法刷新配置")
	}

	cmd := exec.Command(binaryPath, "reload", "--config", caddyfilePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("刷新配置失败: %s", string(output))
	}
	return nil
}

func IsRunning() bool {
	CaddyMu.Lock()
	defer CaddyMu.Unlock()
	return CaddyRunning
}

func IsBinaryExists(binaryPath string) bool {
	_, err := os.Stat(binaryPath)
	return err == nil
}

func collectLogs(pipe io.ReadCloser, source string) {
	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		text := scanner.Text()
		level := "info"
		if containsAny(text, []string{"error", "Error", "ERROR", "fail", "Failure"}) {
			level = "error"
		} else if containsAny(text, []string{"warn", "Warn", "WARNING", "warning"}) {
			level = "warn"
		}

		logChan <- database.CaddyLog{
			LogLevel:    level,
			LogMessage:  text,
			LogSource:   source,
			CreatedTime: time.Now().UnixMilli(),
		}
	}
}

func batchWriteLogs() {
	batch := make([]database.CaddyLog, 0, 50)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case log := <-logChan:
			batch = append(batch, log)
			if len(batch) >= 50 {
				flushLogBatch(batch)
				batch = batch[:0]
			}
		case <-ticker.C:
			if len(batch) > 0 {
				flushLogBatch(batch)
				batch = batch[:0]
			}
		case <-LogStopChan:
			if len(batch) > 0 {
				flushLogBatch(batch)
			}
			return
		}
	}
}

func flushLogBatch(batch []database.CaddyLog) {
	tx, err := database.DB.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare("INSERT INTO t_caddy_log (log_level, log_message, log_source, created_time) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return
	}
	defer stmt.Close()

	for _, l := range batch {
		stmt.Exec(l.LogLevel, l.LogMessage, l.LogSource, l.CreatedTime)
	}

	tx.Commit()
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

type InstallState struct {
	Status     string `json:"status"`
	Progress   int    `json:"progress"`
	Message    string `json:"message"`
	VersionNew string `json:"version_new"`
}

var CurrentInstall = InstallState{Status: "idle"}
var InstallMu sync.Mutex

func GetInstallState() InstallState {
	InstallMu.Lock()
	defer InstallMu.Unlock()
	return CurrentInstall
}

func GetBinaryVersion(binaryPath string) string {
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return ""
	}
	cmd := exec.Command(binaryPath, "version")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return extractVersion(string(out))
}

func extractVersion(output string) string {
	lines := splitLines(output)
	for _, line := range lines {
		if hasPrefix(line, "v") {
			return line
		}
	}
	return ""
}

func GetLatestVersion() string {
	resp, err := defaultHTTPClient().Get("https://api.github.com/repos/caddyserver/caddy/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := jsonDecode(resp.Body, &release); err != nil {
		return ""
	}
	return release.TagName
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func DownloadAndInstall(targetPath string) {
	InstallMu.Lock()
	CurrentInstall = InstallState{Status: "downloading", Progress: 0, Message: "正在获取最新版本..."}
	InstallMu.Unlock()

	platform := runtimeGOOS() + "_" + runtimeGOARCH()
	targetDir := targetPath[:lastSlash(targetPath)]

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "创建目录失败: " + err.Error()}
		InstallMu.Unlock()
		return
	}

	resp, err := defaultHTTPClient().Get("https://api.github.com/repos/caddyserver/caddy/releases/latest")
	if err != nil {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "获取版本信息失败: " + err.Error()}
		InstallMu.Unlock()
		return
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := jsonDecode(resp.Body, &release); err != nil {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "解析版本信息失败"}
		InstallMu.Unlock()
		return
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if containsStr(asset.Name, platform) && !containsStr(asset.Name, ".sig") && !containsStr(asset.Name, ".json") {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "未找到适配当前平台的安装包 (" + platform + ")"}
		InstallMu.Unlock()
		return
	}

	InstallMu.Lock()
	CurrentInstall = InstallState{Status: "downloading", Progress: 30, Message: "正在下载 Caddy " + release.TagName + " ...", VersionNew: release.TagName}
	InstallMu.Unlock()

	dlResp, err := defaultHTTPClient().Get(downloadURL)
	if err != nil {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "下载失败: " + err.Error()}
		InstallMu.Unlock()
		return
	}
	defer dlResp.Body.Close()

	InstallMu.Lock()
	CurrentInstall = InstallState{Status: "extracting", Progress: 60, Message: "正在解压..."}
	InstallMu.Unlock()

	if containsStr(downloadURL, ".tar.gz") || containsStr(downloadURL, ".tgz") {
		if err := extractTarGz(dlResp.Body, targetDir); err != nil {
			InstallMu.Lock()
			CurrentInstall = InstallState{Status: "error", Message: "解压失败: " + err.Error()}
			InstallMu.Unlock()
			return
		}
	} else if containsStr(downloadURL, ".zip") {
		if err := extractZip(dlResp.Body, targetDir); err != nil {
			InstallMu.Lock()
			CurrentInstall = InstallState{Status: "error", Message: "解压失败: " + err.Error()}
			InstallMu.Unlock()
			return
		}
	} else {
		if err := copyToFile(dlResp.Body, targetPath); err != nil {
			InstallMu.Lock()
			CurrentInstall = InstallState{Status: "error", Message: "写入文件失败: " + err.Error()}
			InstallMu.Unlock()
			return
		}
	}

	if err := os.Chmod(targetPath, 0755); err != nil {
		InstallMu.Lock()
		CurrentInstall = InstallState{Status: "error", Message: "设置执行权限失败"}
		InstallMu.Unlock()
		return
	}

	InstallMu.Lock()
	CurrentInstall = InstallState{Status: "done", Progress: 100, Message: "Caddy " + release.TagName + " 安装完成", VersionNew: release.TagName}
	InstallMu.Unlock()
}

func lastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return i
		}
	}
	return 0
}

func runtimeGOOS() string { return runtime.GOOS }

func runtimeGOARCH() string {
	a := runtime.GOARCH
	if a == "amd64" {
		return "amd64"
	}
	if a == "arm64" {
		return "arm64"
	}
	return a
}

func defaultHTTPClient() *http.Client {
	return &http.Client{Timeout: 30 * time.Minute}
}

func jsonDecode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && indexSubstr(s, sub) >= 0
}

func indexSubstr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func extractTarGz(r io.Reader, destDir string) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := destDir + "/" + hdr.Name
		if hdr.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		os.MkdirAll(destDir, 0755)
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, tr); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

func extractZip(r io.Reader, destDir string) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range zr.File {
		target := destDir + "/" + f.Name
		if f.FileInfo().IsDir() {
			os.MkdirAll(target, 0755)
			continue
		}

		os.MkdirAll(destDir, 0755)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		io.Copy(out, rc)
		out.Close()
		rc.Close()
	}
	return nil
}

func copyToFile(r io.Reader, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
