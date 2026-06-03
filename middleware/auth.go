package middleware

import (
	"net/http"
	"strings"
	"time"

	"caddy_server/database"

	"github.com/gin-gonic/gin"
)

var openRoutes = map[string]bool{
	"/init":  true,
	"/login": true,
}

func AuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if openRoutes[path] {
			c.Next()
			return
		}

		if strings.HasPrefix(path, "/static/") {
			c.Next()
			return
		}

		if c.Request.Method == "POST" && path == "/init" {
			c.Next()
			return
		}

		tokenStr := c.GetHeader("Token")
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}
		if tokenStr == "" {
			tokenStr, _ = c.Cookie("auth_token")
		}

		if tokenStr == "" {
			unauthorized(c)
			return
		}

		var userID int64
		var username string
		var expiresAt int64

		err := database.DB.QueryRow(`
			SELECT t_token.user_id, t_user.username, t_token.expires_at
			FROM t_token
			INNER JOIN t_user ON t_user.id = t_token.user_id
			WHERE t_token.token = ?
			AND t_token.delete_time = 0
			AND t_user.delete_time = 0
		`, tokenStr).Scan(&userID, &username, &expiresAt)

		if err != nil {
			unauthorized(c)
			return
		}

		if expiresAt < time.Now().UnixMilli() {
			database.DB.Exec("DELETE FROM t_token WHERE token = ?", tokenStr)
			unauthorized(c)
			return
		}

		c.Set("user_id", userID)
		c.Set("username", username)

		c.Next()
	}
}

func unauthorized(c *gin.Context) {
	if isAPIRequest(c) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "登录已过期，请重新登录",
			"data":    nil,
		})
	} else {
		redirectURL := "/login?redirect=" + c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			redirectURL = "/login?redirect=" + c.Request.URL.Path + "?" + c.Request.URL.RawQuery
		}
		c.Redirect(http.StatusFound, redirectURL)
	}
	c.Abort()
}

func isAPIRequest(c *gin.Context) bool {
	if strings.HasPrefix(c.Request.URL.Path, "/api/") {
		return true
	}
	accept := c.GetHeader("Accept")
	return strings.Contains(accept, "application/json")
}
