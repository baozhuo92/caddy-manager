package handlers

import (
	"net/http"
	"os"
	"time"

	"caddy_server/config"
	"caddy_server/database"
	"caddy_server/middleware"
	"caddy_server/services"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title":       "控制台",
		"HeaderTitle": "控制台",
		"Active":      "dashboard",
		"Username":    middleware.GetUsername(c),
	})
}

func GetStats(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		caddyStatus := "stopped"
		if services.IsRunning() {
			caddyStatus = "running"
		}

		binaryExists := services.IsBinaryExists(cfg.Caddy.BinaryPath)

		var totalSites, enabledSites, disabledSites int
		database.DB.QueryRow("SELECT COUNT(*) FROM t_site WHERE delete_time = 0").Scan(&totalSites)
		database.DB.QueryRow("SELECT COUNT(*) FROM t_site WHERE delete_time = 0 AND status = 'enabled'").Scan(&enabledSites)
		database.DB.QueryRow("SELECT COUNT(*) FROM t_site WHERE delete_time = 0 AND status = 'disabled'").Scan(&disabledSites)

		var dbSize int64
		if info, err := os.Stat(cfg.Database.Path); err == nil {
			dbSize = info.Size()
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"caddy_status":        caddyStatus,
				"caddy_binary_exists": binaryExists,
				"total_sites":         totalSites,
				"enabled_sites":       enabledSites,
				"disabled_sites":      disabledSites,
				"config": gin.H{
					"server_domain":       database.GetConfig("server_domain"),
					"caddyfile_path":      database.GetConfig("caddyfile_path"),
					"caddyfile_sites_dir": database.GetConfig("caddyfile_sites_dir"),
				},
				"db_size": formatFileSize(dbSize),
				"uptime":  formatUptime(),
			},
		})
	}
}

var startTime = time.Now()

func formatFileSize(size int64) string {
	if size < 1024 {
		return "1 KB"
	}
	return formatSize(size/1024) + " KB"
}

func formatSize(n int64) string {
	if n < 1024 {
		return toString(n)
	}
	return toString(n/1024) + " MB"
}

func toString(n int64) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

func formatUptime() string {
	d := time.Since(startTime)
	if d < time.Minute {
		return "不到1分钟"
	}
	if d < time.Hour {
		return toString(int64(d.Minutes())) + "分钟"
	}
	days := int64(d.Hours()) / 24
	hours := int64(d.Hours()) % 24
	if days > 0 {
		return toString(days) + "天" + toString(hours) + "小时"
	}
	return toString(hours) + "小时"
}

func GetCaddyVersion(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		binaryPath := database.GetConfig("caddy_binary_path")
		if binaryPath == "" {
			binaryPath = cfg.Caddy.BinaryPath
		}

		installed := services.GetBinaryVersion(binaryPath)
		latest := services.GetLatestVersion()

		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"installed":     installed,
				"latest":        latest,
				"has_update":    installed != "" && latest != "" && installed != latest,
				"not_installed": installed == "",
			},
		})
	}
}

func PostCaddyInstall(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		binaryPath := database.GetConfig("caddy_binary_path")
		if binaryPath == "" {
			binaryPath = cfg.Caddy.BinaryPath
		}

		state := services.GetInstallState()
		if state.Status == "downloading" || state.Status == "extracting" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "安装正在进行中"})
			return
		}

		go services.DownloadAndInstall(binaryPath)

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "开始安装"})
	}
}

func GetCaddyInstallProgress(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		state := services.GetInstallState()
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data":    state,
		})
	}
}
