package middleware

import (
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) int64 {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return id.(int64)
}

func GetUsername(c *gin.Context) string {
	name, exists := c.Get("username")
	if !exists {
		return ""
	}
	return name.(string)
}
