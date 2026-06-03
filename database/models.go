package database

type User struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	CreatedBy    int64  `json:"created_by"`
	CreatedTime  int64  `json:"created_time"`
	UpdateBy     int64  `json:"update_by"`
	UpdateTime   int64  `json:"update_time"`
	DeleteTime   int64  `json:"delete_time"`
}

type Token struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	Token       string `json:"token"`
	ExpiresAt   int64  `json:"expires_at"`
	CreatedBy   int64  `json:"created_by"`
	CreatedTime int64  `json:"created_time"`
	UpdateBy    int64  `json:"update_by"`
	UpdateTime  int64  `json:"update_time"`
	DeleteTime  int64  `json:"delete_time"`
}

type Site struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Domain         string `json:"domain"`
	Scheme         string `json:"scheme"`
	ProxyTo        string `json:"proxy_to"`
	Compression    int    `json:"compression"`
	RequestHeaders string `json:"request_headers"`
	PathRoutes     string `json:"path_routes"`
	ExtraConfig    string `json:"extra_config"`
	Status         string `json:"status"`
	CreatedBy      int64  `json:"created_by"`
	CreatedTime    int64  `json:"created_time"`
	UpdateBy       int64  `json:"update_by"`
	UpdateTime     int64  `json:"update_time"`
	DeleteTime     int64  `json:"delete_time"`
}

type Config struct {
	ID          int64  `json:"id"`
	ConfigKey   string `json:"config_key"`
	ConfigValue string `json:"config_value"`
	Description string `json:"description"`
	CreatedBy   int64  `json:"created_by"`
	CreatedTime int64  `json:"created_time"`
	UpdateBy    int64  `json:"update_by"`
	UpdateTime  int64  `json:"update_time"`
	DeleteTime  int64  `json:"delete_time"`
}

type CaddyLog struct {
	ID          int64  `json:"id"`
	LogLevel    string `json:"log_level"`
	LogMessage  string `json:"log_message"`
	LogSource   string `json:"log_source"`
	CreatedTime int64  `json:"created_time"`
}

type LoginLog struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	LoginIP     string `json:"login_ip"`
	LoginTime   int64  `json:"login_time"`
	LogoutTime  int64  `json:"logout_time"`
	UserAgent   string `json:"user_agent"`
	CreatedTime int64  `json:"created_time"`
}
