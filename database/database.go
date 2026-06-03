package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbPath string) error {
	var err error
	dsn := fmt.Sprintf("file:%s?_journal_mode=WAL&_foreign_keys=off", dbPath)
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)
	DB.SetConnMaxLifetime(time.Hour)

	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	return nil
}

func migrate() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS t_user (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(64) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_by BIGINT DEFAULT 0,
			created_time BIGINT NOT NULL,
			update_by BIGINT DEFAULT 0,
			update_time BIGINT DEFAULT 0,
			delete_time BIGINT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS t_token (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id BIGINT NOT NULL,
			token VARCHAR(64) NOT NULL UNIQUE,
			expires_at BIGINT NOT NULL,
			created_by BIGINT DEFAULT 0,
			created_time BIGINT NOT NULL,
			update_by BIGINT DEFAULT 0,
			update_time BIGINT DEFAULT 0,
			delete_time BIGINT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS t_site (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(128) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			scheme VARCHAR(8) DEFAULT 'https',
			proxy_to VARCHAR(512) NOT NULL,
			compression INTEGER DEFAULT 0,
			request_headers TEXT DEFAULT '',
			path_routes TEXT DEFAULT '',
			extra_config TEXT DEFAULT '',
			status VARCHAR(16) DEFAULT 'enabled',
			created_by BIGINT DEFAULT 0,
			created_time BIGINT NOT NULL,
			update_by BIGINT DEFAULT 0,
			update_time BIGINT DEFAULT 0,
			delete_time BIGINT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS t_config (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			config_key VARCHAR(128) NOT NULL UNIQUE,
			config_value TEXT NOT NULL,
			description VARCHAR(256) DEFAULT '',
			created_by BIGINT DEFAULT 0,
			created_time BIGINT NOT NULL,
			update_by BIGINT DEFAULT 0,
			update_time BIGINT DEFAULT 0,
			delete_time BIGINT DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS t_caddy_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			log_level VARCHAR(16) DEFAULT 'info',
			log_message TEXT DEFAULT '',
			log_source VARCHAR(64) DEFAULT '',
			created_time BIGINT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS t_login_log (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id BIGINT NOT NULL,
			login_ip VARCHAR(64) DEFAULT '',
			login_time BIGINT NOT NULL,
			logout_time BIGINT DEFAULT 0,
			user_agent VARCHAR(512) DEFAULT '',
			created_time BIGINT NOT NULL
		)`,
	}

	for _, ddl := range tables {
		if _, err := DB.Exec(ddl); err != nil {
			return fmt.Errorf("migration failed: %w\nSQL: %s", err, ddl)
		}
	}

	return nil
}

func HasAdmin() bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM t_user WHERE delete_time = 0").Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func GetConfig(key string) string {
	var value string
	err := DB.QueryRow("SELECT config_value FROM t_config WHERE config_key = ? AND delete_time = 0", key).Scan(&value)
	if err != nil {
		return ""
	}
	return value
}

func SetConfig(key, value, description string) {
	now := time.Now().UnixMilli()
	DB.Exec(`INSERT INTO t_config (config_key, config_value, description, created_time)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(config_key) DO UPDATE SET config_value = ?, update_time = ?`,
		key, value, description, now, value, now)
}

func CleanExpiredTokens() {
	now := time.Now().UnixMilli()
	DB.Exec("DELETE FROM t_token WHERE expires_at < ?", now)
}

func CleanOldLogs() {
	DB.Exec("DELETE FROM t_caddy_log WHERE id NOT IN (SELECT id FROM t_caddy_log ORDER BY id DESC LIMIT 10000)")
	DB.Exec("DELETE FROM t_login_log WHERE id NOT IN (SELECT id FROM t_login_log ORDER BY id DESC LIMIT 1000)")
}
