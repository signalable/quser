// quser/internal/config/config.go
package config

import (
    "os"
    "strconv"
    "time"
    "github.com/joho/godotenv"
)

type Config struct {
    Server      ServerConfig
    MongoDB     MongoDBConfig
    AuthService AuthServiceConfig
    LogLevel    string
}

type ServerConfig struct {
    Host string
    Port string
}

type MongoDBConfig struct {
    URI      string
    Database string
}

type AuthServiceConfig struct {
    URL     string
    Timeout time.Duration
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        return nil, err
    }

    timeoutSec, err := strconv.Atoi(getEnv("AUTH_SERVICE_TIMEOUT_SEC", "5"))
    if err != nil {
        timeoutSec = 5
    }

    return &Config{
        Server: ServerConfig{
            Host: getEnv("SERVER_HOST", "0.0.0.0"),
            Port: getEnv("SERVER_PORT", "8081"), // Auth는 8080, User는 8081 사용
        },
        MongoDB: MongoDBConfig{
            URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
            Database: getEnv("MONGODB_DATABASE", "user_db"),
        },
        AuthService: AuthServiceConfig{
            URL:     getEnv("AUTH_SERVICE_URL", "http://localhost:8080"),
            Timeout: time.Duration(timeoutSec) * time.Second,
        },
        LogLevel: getEnv("LOG_LEVEL", "debug"),
    }, nil
}

func getEnv(key, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}