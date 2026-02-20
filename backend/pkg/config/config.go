package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// 服务器配置
	Port     string
	LogLevel string

	// CHSI登录配置
	ChsiUsername string
	ChsiPassword string

	// 邮件配置
	SMTPServer string
	SMTPPort   int
	SMTPUser   string
	SMTPPass   string

	// 数据库配置
	DatabaseDSN string

	// 查询配置
	QueryInterval      int
	ClearDBOnStart     bool
	InitialUserEntries string
}

func Load() (*Config, error) {
	// 尝试加载.env文件（可选）
	_ = godotenv.Load()

	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		ChsiUsername:       os.Getenv("CHSI_USERNAME"),
		ChsiPassword:       os.Getenv("CHSI_PASSWORD"),
		SMTPServer:         getEnv("SMTP_SERVER", "smtp.gmail.com"),
		SMTPPort:           getEnvInt("SMTP_PORT", 587),
		SMTPUser:           getEnv("SMTP_USER", ""),
		SMTPPass:           getEnv("SMTP_PASSWORD", ""),
		DatabaseDSN:        getEnv("DATABASE_DSN", "./data/chsi.db"),
		QueryInterval:      getEnvInt("QUERY_INTERVAL", 3600),
		ClearDBOnStart:     getEnvBool("CLEAR_DB_ON_START", false),
		InitialUserEntries: getEnv("INITIAL_USER_ENTRIES", ""),
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
