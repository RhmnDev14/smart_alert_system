package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DatabaseURL string

	// Waha Server
	WahaServerURL string
	WahaAPIKey    string

	// AI Configuration
	AIProvider string
	AIApiKey   string
	AIModel    string
	AIBaseURL  string // For Ollama or other OpenAI-compatible APIs

	// Application
	AppEnv   string
	AppPort  string
	Timezone string
	AppName  string

	// Scheduler
	MorningAlertTime   string
	EveningSummaryTime string
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	cfg := &Config{
		// Database
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "postgres"),
		DBName:      getEnv("DB_NAME", "smart_alert_db"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// Waha Server
		WahaServerURL: getEnv("WAHA_SERVER_URL", "http://localhost:3000"),
		WahaAPIKey:    getEnv("WAHA_API_KEY", ""),

		// AI Configuration
		AIProvider: getEnv("AI_PROVIDER", "openai"),
		AIApiKey:   getEnv("AI_API_KEY", ""),
		AIModel:    getEnv("AI_MODEL", "gpt-3.5-turbo"),
		AIBaseURL:  getEnv("AI_BASE_URL", ""), // For Ollama: http://localhost:11434/v1

		// Application
		AppEnv:   getEnv("APP_ENV", "development"),
		AppPort:  getEnv("APP_PORT", "8080"),
		Timezone: getEnv("TIMEZONE", "Asia/Jakarta"),
		AppName:  "Smart Alert System",

		// Scheduler
		MorningAlertTime:   getEnv("MORNING_ALERT_TIME", "05:00"),
		EveningSummaryTime: getEnv("EVENING_SUMMARY_TIME", "22:00"),
	}

	// Build DatabaseURL if not provided
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = buildDatabaseURL(cfg)
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func buildDatabaseURL(cfg *Config) string {
	return "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=" + cfg.DBSSLMode
}

func (c *Config) GetLocation() (*time.Location, error) {
	return time.LoadLocation(c.Timezone)
}
