package handlers

import (
	"net/http"
	"time"

	"caddy_server/database"
	"caddy_server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetInitPage(c *gin.Context) {
	if database.HasAdmin() {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	c.HTML(http.StatusOK, "init.html", gin.H{"title": "初始化管理员"})
}

func PostInit(c *gin.Context) {
	if database.HasAdmin() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "管理员已存在"})
		return
	}

	var req struct {
		Username        string `json:"username"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	if len(req.Username) < 4 || len(req.Username) > 64 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名长度需为4-64个字符"})
		return
	}

	if len(req.Password) < 6 || len(req.Password) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "密码长度需为6-128个字符"})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "两次输入的密码不一致"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	now := time.Now().UnixMilli()

	result, err := database.DB.Exec(`
		INSERT INTO t_user (username, password_hash, created_time)
		VALUES (?, ?, ?)
	`, req.Username, string(hash), now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建管理员失败"})
		return
	}

	userID, _ := result.LastInsertId()

	database.DB.Exec(`
		INSERT INTO t_login_log (user_id, login_ip, login_time, user_agent, created_time)
		VALUES (?, ?, ?, ?, ?)
	`, userID, c.ClientIP(), now, c.GetHeader("User-Agent"), now)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "管理员创建成功"})
}

func GetLoginPage(c *gin.Context) {
	if !database.HasAdmin() {
		c.Redirect(http.StatusFound, "/init")
		return
	}
	redirect := c.Query("redirect")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":    "登录",
		"redirect": redirect,
	})
}

func PostLogin(c *gin.Context) {
	if !database.HasAdmin() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请先初始化管理员账号"})
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请求参数错误"})
		return
	}

	var user database.User
	err := database.DB.QueryRow(`
		SELECT id, username, password_hash FROM t_user
		WHERE username = ? AND delete_time = 0
	`, req.Username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	tokenStr := uuid.New().String()
	now := time.Now().UnixMilli()
	expiresAt := now + 24*60*60*1000

	_, err = database.DB.Exec(`
		INSERT INTO t_token (user_id, token, expires_at, created_time)
		VALUES (?, ?, ?, ?)
	`, user.ID, tokenStr, expiresAt, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成令牌失败"})
		return
	}

	database.DB.Exec(`
		INSERT INTO t_login_log (user_id, login_ip, login_time, user_agent, created_time)
		VALUES (?, ?, ?, ?, ?)
	`, user.ID, c.ClientIP(), now, c.GetHeader("User-Agent"), now)

	maxAge := 24 * 60 * 60
	c.SetCookie("auth_token", tokenStr, maxAge, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data": gin.H{
			"token": tokenStr,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	})
}

func GetLogout(c *gin.Context) {
	tokenStr := c.GetHeader("Token")
	if tokenStr == "" {
		tokenStr, _ = c.Cookie("auth_token")
	}
	if tokenStr != "" {
		database.DB.Exec("DELETE FROM t_token WHERE token = ?", tokenStr)
	}

	userID := middleware.GetUserID(c)
	now := time.Now().UnixMilli()
	database.DB.Exec(`
		UPDATE t_login_log SET logout_time = ?
		WHERE user_id = ? AND logout_time = 0
		ORDER BY id DESC LIMIT 1
	`, now, userID)

	c.SetCookie("auth_token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
