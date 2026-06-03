package middleware

import (
	"net/http"
	"strings"

	"caddy_server/database"

	"github.com/gin-gonic/gin"
)

var allowedWithoutAdmin = []string{"/init", "/login", "/static/"}

func InitCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if database.HasAdmin() {
			c.Next()
			return
		}

		path := c.Request.URL.Path

		for _, prefix := range allowedWithoutAdmin {
			if strings.HasPrefix(path, prefix) {
				c.Next()
				return
			}
		}

		if c.Request.Method == "POST" && path == "/init" {
			c.Next()
			return
		}

		c.Redirect(http.StatusFound, "/init")
		c.Abort()
	}
}
