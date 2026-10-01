package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	DBDriver    string
	JWTSecret   string
	UploadDir   string
	BaseURL     string
}

func Load() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbURL := os.Getenv("DATABASE_URL")
	dbDriver := "sqlite"
	if dbURL != "" && (len(dbURL) >= 8 && dbURL[:8] == "postgres") {
		dbDriver = "pgx"
	} else if dbURL == "" {
		dbURL = "defect_management.db"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		DBDriver:    dbDriver,
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-for-development-change-in-prod"),
		UploadDir:   getEnv("UPLOAD_DIR", "uploads"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:"+port),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
