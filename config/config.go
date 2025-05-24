package config

import (
	"os"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
		DBName   string
		SSLMode  string
	}
	Server struct {
		Port int
	}
	WireGuard struct {
		ConfigPath    string
		ContainerName string
	}
}

func LoadConfig() *Config {
	config := &Config{}

	// Database config
	config.Database.Host = getEnv("DB_HOST", "localhost")
	config.Database.Port = getEnvAsInt("DB_PORT", 5432)
	config.Database.Username = getEnv("DB_USER", "admin")
	config.Database.Password = getEnv("DB_PASSWORD", "postgres")
	config.Database.DBName = getEnv("DB_NAME", "wg_proj")
	config.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Server config
	config.Server.Port = getEnvAsInt("SERVER_PORT", 8080)

	// WireGuard config
	config.WireGuard.ConfigPath = getEnv("WG_CONFIG_PATH", "/etc/wireguard/wg0.conf")
	config.WireGuard.ContainerName = getEnv("WG_CONTAINER_NAME", "wireguard")

	return config
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		if value, err := strconv.Atoi(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
