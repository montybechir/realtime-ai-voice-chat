package config

import "os"

type Config struct {
	Host        string
	Port        string
	DatabaseURL string
	LogLevel    string
	APIKey      string
	Endpoint    string
}

func LoadConfig() (Config, error) {
	return Config{
		Host:        getEnv("HOST", "localhost"),
		Port:        getEnv("PORT", "5555"),
		DatabaseURL: getEnv("DATABASE_URL", "pg"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
		APIKey:      getEnv("AZURE_OPENAI_API_KEY", ""),
		Endpoint:    getEnv("AZURE_OPENAI_ENDPOINT", ""),
	}, nil
}

func (c Config) ServerAddress() string {
	return c.Host + ":" + c.Port
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
