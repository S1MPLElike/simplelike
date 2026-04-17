package config

import "time"

var (
	MySQLConfig = mysqlConfig{
		Username: "root",
		Password: "187415157",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "GOYES",
		Charset:  "utf8mb4",
	}

	RedisConfig = redisConfig{
		Addr:     "127.0.0.1:6379",
		Password: "",
		Database: 3,
	}

	TokenConfig = tokenConfig{
		ExpireTime: 7 * 24 * time.Hour, // Token过期时间（7天）
		CookieName: "GOYES_token",      // Cookie名称（存储Token）
	}
)

type tokenConfig struct {
	ExpireTime time.Duration
	CookieName string
}

type mysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Charset  string
}
type redisConfig struct {
	Addr     string
	Password string
	Database int
}
