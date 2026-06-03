package handlers

import (
	"net/http"
	"strconv"
	"time"

	"caddy_server/database"
	"caddy_server/middleware"
	"caddy_server/services"

	"github.com/gin-gonic/gin"
)

func GetSites(c *gin.Context) {
	c.HTML(http.StatusOK, "site_list.html", gin.H{
		"title":       "网站管理",
		"HeaderTitle": "网站管理",
		"Active":      "sites",
		"Username":    middleware.GetUsername(c),
	})
}

func GetSitesAPI(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	search := c.Query("search")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var total int64
	countSQL := "SELECT COUNT(*) FROM t_site WHERE delete_time = 0"
	args := []interface{}{}
	if search != "" {
		countSQL += " AND (domain LIKE ? OR name LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}
	database.DB.QueryRow(countSQL, args...).Scan(&total)

	querySQL := "SELECT id, name, domain, scheme, proxy_to, compression, request_headers, path_routes, extra_config, status, created_time, update_time FROM t_site WHERE delete_time = 0"
	queryArgs := make([]interface{}, len(args))
	copy(queryArgs, args)
	if search != "" {
		querySQL += " AND (domain LIKE ? OR name LIKE ?)"
	}
	querySQL += " ORDER BY id DESC LIMIT ? OFFSET ?"
	queryArgs = append(queryArgs, pageSize, offset)

	rows, err := database.DB.Query(querySQL, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "查询失败"})
		return
	}
	defer rows.Close()

	var sites []database.Site
	for rows.Next() {
		var s database.Site
		rows.Scan(&s.ID, &s.Name, &s.Domain, &s.Scheme, &s.ProxyTo, &s.Compression,
			&s.RequestHeaders, &s.PathRoutes, &s.ExtraConfig, &s.Status, &s.CreatedTime, &s.UpdateTime)
		sites = append(sites, s)
	}

	if sites == nil {
		sites = []database.Site{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      sites,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func GetSiteNew(c *gin.Context) {
	c.HTML(http.StatusOK, "site_form.html", gin.H{
		"title":       "新增网站",
		"HeaderTitle": "新增网站",
		"Active":      "site_new",
		"Username":    middleware.GetUsername(c),
		"Site":        database.Site{Scheme: "https", Status: "enabled"},
		"IsEdit":      false,
	})
}

func PostSite(c *gin.Context) {
	var req struct {
		Name           string `json:"name"`
		Domain         string `json:"domain"`
		Scheme         string `json:"scheme"`
		ProxyTo        string `json:"proxy_to"`
		Compression    bool   `json:"compression"`
		RequestHeaders string `json:"request_headers"`
		PathRoutes     string `json:"path_routes"`
		ExtraConfig    string `json:"extra_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Name == "" || req.Domain == "" || req.ProxyTo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "网站名称、域名和代理地址为必填项"})
		return
	}

	scheme := req.Scheme
	if scheme == "" {
		scheme = "https"
	}

	compression := 0
	if req.Compression {
		compression = 1
	}

	now := time.Now().UnixMilli()
	userID := middleware.GetUserID(c)

	result, err := database.DB.Exec(`
		INSERT INTO t_site (name, domain, scheme, proxy_to, compression, request_headers, path_routes, extra_config, status, created_by, created_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'enabled', ?, ?)
	`, req.Name, req.Domain, scheme, req.ProxyTo, compression, req.RequestHeaders, req.PathRoutes, req.ExtraConfig, userID, now)
	if err != nil {
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该域名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建失败"})
		return
	}

	siteID, _ := result.LastInsertId()

	site := database.Site{
		ID:             siteID,
		Name:           req.Name,
		Domain:         req.Domain,
		Scheme:         scheme,
		ProxyTo:        req.ProxyTo,
		Compression:    compression,
		RequestHeaders: req.RequestHeaders,
		PathRoutes:     req.PathRoutes,
		ExtraConfig:    req.ExtraConfig,
		Status:         "enabled",
	}

	sitesDir := database.GetConfig("caddyfile_sites_dir")
	services.WriteSiteConfig(site, sitesDir)

	caddyfilePath := database.GetConfig("caddyfile_path")
	serverDomain := database.GetConfig("server_domain")
	services.WriteMainCaddyfile(caddyfilePath, serverDomain, sitesDir)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功"})
}

func GetSiteEdit(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var site database.Site
	err := database.DB.QueryRow(`
		SELECT id, name, domain, scheme, proxy_to, compression, request_headers, path_routes, extra_config, status
		FROM t_site WHERE id = ? AND delete_time = 0
	`, id).Scan(&site.ID, &site.Name, &site.Domain, &site.Scheme, &site.ProxyTo,
		&site.Compression, &site.RequestHeaders, &site.PathRoutes, &site.ExtraConfig, &site.Status)
	if err != nil {
		c.Redirect(http.StatusFound, "/sites")
		return
	}

	c.HTML(http.StatusOK, "site_form.html", gin.H{
		"title":       "编辑网站",
		"HeaderTitle": "编辑网站",
		"Active":      "sites",
		"Username":    middleware.GetUsername(c),
		"Site":        site,
		"IsEdit":      true,
	})
}

func UpdateSite(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var req struct {
		Name           string `json:"name"`
		Domain         string `json:"domain"`
		Scheme         string `json:"scheme"`
		ProxyTo        string `json:"proxy_to"`
		Compression    bool   `json:"compression"`
		RequestHeaders string `json:"request_headers"`
		PathRoutes     string `json:"path_routes"`
		ExtraConfig    string `json:"extra_config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if req.Name == "" || req.Domain == "" || req.ProxyTo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "网站名称、域名和代理地址为必填项"})
		return
	}

	scheme := req.Scheme
	if scheme == "" {
		scheme = "https"
	}

	compression := 0
	if req.Compression {
		compression = 1
	}

	now := time.Now().UnixMilli()
	userID := middleware.GetUserID(c)

	_, err := database.DB.Exec(`
		UPDATE t_site SET name=?, domain=?, scheme=?, proxy_to=?, compression=?, request_headers=?, path_routes=?, extra_config=?, update_by=?, update_time=?
		WHERE id=? AND delete_time=0
	`, req.Name, req.Domain, scheme, req.ProxyTo, compression, req.RequestHeaders, req.PathRoutes, req.ExtraConfig, userID, now, id)
	if err != nil {
		if isUniqueConstraintError(err) {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "该域名已存在"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "更新失败"})
		return
	}

	sitesDir := database.GetConfig("caddyfile_sites_dir")

	var status string
	database.DB.QueryRow("SELECT status FROM t_site WHERE id = ?", id).Scan(&status)

	if status == "enabled" {
		site := database.Site{
			ID:             id,
			Name:           req.Name,
			Domain:         req.Domain,
			Scheme:         scheme,
			ProxyTo:        req.ProxyTo,
			Compression:    compression,
			RequestHeaders: req.RequestHeaders,
			PathRoutes:     req.PathRoutes,
			ExtraConfig:    req.ExtraConfig,
			Status:         status,
		}
		services.WriteSiteConfig(site, sitesDir)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功"})
}

func DeleteSite(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var domain string
	database.DB.QueryRow("SELECT domain FROM t_site WHERE id = ? AND delete_time = 0", id).Scan(&domain)

	now := time.Now().UnixMilli()
	_, err := database.DB.Exec("UPDATE t_site SET delete_time = ? WHERE id = ?", now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	sitesDir := database.GetConfig("caddyfile_sites_dir")
	services.DeleteSiteConfig(domain, sitesDir)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

func ToggleSite(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)

	var currentStatus string
	var site database.Site
	err := database.DB.QueryRow(`
		SELECT id, name, domain, scheme, proxy_to, compression, request_headers, path_routes, extra_config, status
		FROM t_site WHERE id = ? AND delete_time = 0
	`, id).Scan(&site.ID, &site.Name, &site.Domain, &site.Scheme, &site.ProxyTo,
		&site.Compression, &site.RequestHeaders, &site.PathRoutes, &site.ExtraConfig, &currentStatus)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "站点不存在"})
		return
	}

	newStatus := "enabled"
	msg := "已启用"
	if currentStatus == "enabled" {
		newStatus = "disabled"
		msg = "已停用"
	}

	now := time.Now().UnixMilli()
	_, err = database.DB.Exec("UPDATE t_site SET status = ?, update_time = ? WHERE id = ?", newStatus, now, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	sitesDir := database.GetConfig("caddyfile_sites_dir")
	if newStatus == "enabled" {
		site.Status = "enabled"
		services.WriteSiteConfig(site, sitesDir)
	} else {
		services.DeleteSiteConfig(site.Domain, sitesDir)
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": msg})
}

func isUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return containsStr(msg, "UNIQUE") || containsStr(msg, "unique")
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstr(s, substr)
}

func searchSubstr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
