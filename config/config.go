package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
    // Server
    Env  string
    Port string
    
    // Database
    DatabaseURL        string
    DBMaxConns         int32
    DBMinConns         int32
    DBMaxConnLifetime  time.Duration
    DBMaxConnIdleTime  time.Duration
    
    // Gin
    GinMode  string
    LogLevel string
}

func Load() *Config {
    return &Config{
        // Server
        Env:  getEnv("ENV", "production"),
        Port: getEnv("PORT", "8080"),
        
        // Database
        DatabaseURL:        getEnv("DATABASE_URL", "postgres://userapi:userpass123@localhost:5432/userdb?sslmode=disable"),
        DBMaxConns:         getEnvAsInt32("DB_MAX_CONNS", 30),
        DBMinConns:         getEnvAsInt32("DB_MIN_CONNS", 5),
        DBMaxConnLifetime:  getEnvAsDuration("DB_MAX_CONN_LIFETIME", time.Hour),
        DBMaxConnIdleTime:  getEnvAsDuration("DB_MAX_CONN_IDLE_TIME", 30*time.Minute),
        
        // Gin
        GinMode:  getEnv("GIN_MODE", "debug"),
        LogLevel: getEnv("LOG_LEVEL", "info"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt32(key string, defaultValue int32) int32 {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return int32(intValue)
        }
    }
    return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return time.Duration(intValue) * time.Second
        }
    }
    return defaultValue
}