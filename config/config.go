package config

import (
    "log"
    "os"
    "strings"

    "github.com/joho/godotenv"
)

type Config struct {
    AzureEndpoint   string
    AzureAPIKey     string
    AzureModel      string
    AzureDeployment string
    ServerPort      string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    return &Config{
        AzureEndpoint:   getEnv("AZURE_OPENAI_ENDPOINT", ""),
        AzureAPIKey:     getEnv("AZURE_OPENAI_API_KEY", ""),
        AzureModel:      getEnv("AZURE_OPENAI_MODEL", "DeepSeek-V3.2"),
        AzureDeployment: getEnv("AZURE_OPENAI_DEPLOYMENT", "DeepSeek-V3.2"),
        ServerPort:      getEnv("SERVER_PORT", "8080"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}