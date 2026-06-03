package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"caddy_server/config"
	"caddy_server/database"
	"caddy_server/handlers"
	"caddy_server/middleware"
	"caddy_server/services"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbDir := filepath.Dir(cfg.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	if err := database.Init(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.DB.Close()

	database.CleanExpiredTokens()
	database.CleanOldLogs()

	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			database.CleanExpiredTokens()
			database.CleanOldLogs()
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.SetFuncMap(map[string]interface{}{})

	r.LoadHTMLGlob("templates/*.html")

	r.Static("/static", "./static")

	r.Use(middleware.InitCheck())
	r.Use(middleware.AuthToken())

	r.GET("/init", handlers.GetInitPage)
	r.POST("/init", handlers.PostInit)
	r.GET("/login", handlers.GetLoginPage)
	r.POST("/login", handlers.PostLogin)
	r.GET("/logout", handlers.GetLogout)

	r.GET("/", handlers.GetDashboard)
	r.GET("/api/stats", handlers.GetStats(cfg))
	r.GET("/api/caddy/version", handlers.GetCaddyVersion(cfg))
	r.POST("/api/caddy/install", handlers.PostCaddyInstall(cfg))
	r.GET("/api/caddy/install/progress", handlers.GetCaddyInstallProgress(cfg))

	r.GET("/sites", handlers.GetSites)
	r.GET("/api/sites", handlers.GetSitesAPI)
	r.GET("/sites/new", handlers.GetSiteNew)
	r.POST("/api/sites", handlers.PostSite)
	r.GET("/sites/:id/edit", handlers.GetSiteEdit)
	r.PUT("/api/sites/:id", handlers.UpdateSite)
	r.DELETE("/api/sites/:id", handlers.DeleteSite)
	r.PUT("/api/sites/:id/toggle", handlers.ToggleSite)

	r.GET("/caddy/manage", handlers.GetCaddyManage)
	r.POST("/caddy/start", handlers.PostCaddyStart(cfg))
	r.POST("/caddy/stop", handlers.PostCaddyStop(cfg))
	r.POST("/caddy/reload", handlers.PostCaddyReload(cfg))
	r.POST("/caddy/config", handlers.PostCaddyConfig)

	r.GET("/logs", handlers.GetLogsPage)
	r.GET("/api/logs", handlers.GetLogsAPI)
	r.DELETE("/api/logs", handlers.ClearLogs)

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title":   "404",
			"code":    "404",
			"message": "页面不存在",
		})
		c.Abort()
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		services.StopCaddy()
		os.Exit(0)
	}()

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
