package handlers

import (
	"net/http"
	"strconv"

	"caddy_server/database"
	"caddy_server/middleware"

	"github.com/gin-gonic/gin"
)

func GetLogsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "system_log.html", gin.H{
		"title":       "系统日志",
		"HeaderTitle": "系统日志",
		"Active":      "logs",
		"Username":    middleware.GetUsername(c),
	})
}

func GetLogsAPI(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	level := c.Query("level")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var total int64
	countSQL := "SELECT COUNT(*) FROM t_caddy_log WHERE 1=1"
	countArgs := []interface{}{}
	if level != "" && level != "all" {
		countSQL += " AND log_level = ?"
		countArgs = append(countArgs, level)
	}
	if search != "" {
		countSQL += " AND log_message LIKE ?"
		countArgs = append(countArgs, "%"+search+"%")
	}
	database.DB.QueryRow(countSQL, countArgs...).Scan(&total)

	querySQL := "SELECT id, log_level, log_message, log_source, created_time FROM t_caddy_log WHERE 1=1"
	queryArgs := make([]interface{}, len(countArgs))
	copy(queryArgs, countArgs)
	if level != "" && level != "all" {
		querySQL += " AND log_level = ?"
	}
	if search != "" {
		querySQL += " AND log_message LIKE ?"
	}
	querySQL += " ORDER BY id DESC LIMIT ? OFFSET ?"
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := database.DB.Query(querySQL, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	defer rows.Close()

	var logs []database.CaddyLog
	for rows.Next() {
		var l database.CaddyLog
		rows.Scan(&l.ID, &l.LogLevel, &l.LogMessage, &l.LogSource, &l.CreatedTime)
		logs = append(logs, l)
	}
	if logs == nil {
		logs = []database.CaddyLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func ClearLogs(c *gin.Context) {
	database.DB.Exec("DELETE FROM t_caddy_log")
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "日志已清空"})
}
