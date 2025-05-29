package main

import (
	"os"
)

type Config struct {
	Profile string
	MySQLUser string
	MySQLPassword string
	MySQLHost string
	MySQLDB string
}

func LoadConfig() Config {
	return Config{
		Profile:      getEnv("PROFILE", "dev"),
		MySQLUser:    getEnv("MYSQL_USER", "root"),
		MySQLPassword:getEnv("MYSQL_PASSWORD", ""),
		MySQLHost:    getEnv("MYSQL_HOST", "127.0.0.1:3306"),
		MySQLDB:      getEnv("MYSQL_DB", "learning"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
