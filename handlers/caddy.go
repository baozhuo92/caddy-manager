package handlers

import (
	"net/http"

	"caddy_server/config"
	"caddy_server/database"
	"caddy_server/middleware"
	"caddy_server/services"

	"github.com/gin-gonic/gin"
)

func GetCaddyManage(c *gin.Context) {
	c.HTML(http.StatusOK, "caddy_manage.html", gin.H{
		"title":       "Caddy 管理",
		"HeaderTitle": "Caddy 管理",
		"Active":      "caddy_manage",
		"Username":    middleware.GetUsername(c),
	})
}

func PostCaddyStart(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		binaryPath := database.GetConfig("caddy_binary_path")
		if binaryPath == "" {
			binaryPath = cfg.Caddy.BinaryPath
		}

		caddyfilePath := database.GetConfig("caddyfile_path")
		if caddyfilePath == "" {
			caddyfilePath = cfg.Caddy.CaddyfilePath
		}

		sitesDir := database.GetConfig("caddyfile_sites_dir")
		if sitesDir == "" {
			sitesDir = cfg.Caddy.CaddyfileSitesDir
		}

		serverDomain := database.GetConfig("server_domain")
		if serverDomain == "" {
			serverDomain = cfg.Caddy.ServerDomain
		}

		if err := services.RegenerateAllSites(sitesDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成站点配置失败: " + err.Error()})
			return
		}

		services.WriteMainCaddyfile(caddyfilePath, serverDomain, sitesDir)

		if err := services.StartCaddy(binaryPath, caddyfilePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Caddy 已启动"})
	}
}

func PostCaddyStop(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := services.StopCaddy(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "Caddy 已停止"})
	}
}

func PostCaddyReload(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		binaryPath := database.GetConfig("caddy_binary_path")
		if binaryPath == "" {
			binaryPath = cfg.Caddy.BinaryPath
		}

		caddyfilePath := database.GetConfig("caddyfile_path")
		if caddyfilePath == "" {
			caddyfilePath = cfg.Caddy.CaddyfilePath
		}

		sitesDir := database.GetConfig("caddyfile_sites_dir")
		if sitesDir == "" {
			sitesDir = cfg.Caddy.CaddyfileSitesDir
		}

		serverDomain := database.GetConfig("server_domain")
		if serverDomain == "" {
			serverDomain = cfg.Caddy.ServerDomain
		}

		if err := services.RegenerateAllSites(sitesDir); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成站点配置失败"})
			return
		}
		services.WriteMainCaddyfile(caddyfilePath, serverDomain, sitesDir)

		if err := services.ReloadCaddy(binaryPath, caddyfilePath); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "配置刷新成功"})
	}
}

func PostCaddyConfig(c *gin.Context) {
	var req struct {
		CaddyfilePath     string `json:"caddyfile_path"`
		CaddyfileSitesDir string `json:"caddyfile_sites_dir"`
		ServerDomain      string `json:"server_domain"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.CaddyfilePath != "" {
		database.SetConfig("caddyfile_path", req.CaddyfilePath, "Caddyfile 主文件路径")
	}
	if req.CaddyfileSitesDir != "" {
		database.SetConfig("caddyfile_sites_dir", req.CaddyfileSitesDir, "Caddyfile 子配置文件目录")
	}
	if req.ServerDomain != "" {
		database.SetConfig("server_domain", req.ServerDomain, "管理工具域名")
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "配置已保存"})
}
