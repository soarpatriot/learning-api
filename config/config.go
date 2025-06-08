package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profile       string `yaml:"profile"`
	MySQLUserName string `yaml:"mysql_username"`
	MySQLPassword string `yaml:"mysql_password"`
	MySQLAddress  string `yaml:"mysql_address"`
	MySQLDB       string `yaml:"mysql_db"`
	ApiKey        string `yaml:"api_key"`
	ApiSecret     string `yaml:"api_secret"`
	AppID         string `yaml:"app_id"`
}

type yamlConfig struct {
	Dev        Config `yaml:"dev"`
	Production Config `yaml:"production"`
}

func LoadConfig() Config {
	// Check for environment variable CLOUD_ENV first
	profile := os.Getenv("CLOUD_ENV")
	if profile == "" {
		profile = getEnv("PROFILE", "dev")
	}
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		// Try parent directory (for tests run from subfolders)
		data, err = os.ReadFile("../config.yaml")
		if err != nil {
			panic("failed to read config.yaml: " + err.Error())
		}
	}
	var yc yamlConfig
	err = yaml.Unmarshal(data, &yc)
	if err != nil {
		panic("failed to parse config.yaml: " + err.Error())
	}
	var cfg Config
	switch profile {
	case "production":
		cfg = yc.Production
	default:
		cfg = yc.Dev
	}
	// Allow env vars to override YAML
	cfg.Profile = profile
	if v := os.Getenv("MYSQL_USERNAME"); v != "" {
		cfg.MySQLUserName = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		cfg.MySQLPassword = v
	}
	if v := os.Getenv("MYSQL_ADDRESS"); v != "" {
		cfg.MySQLAddress = v
	}
	if v := os.Getenv("MYSQL_DB"); v != "" {
		cfg.MySQLDB = v
	}
	// config api key and secret
	if v := os.Getenv("API_KEY"); v != "" {
		cfg.ApiKey = v
	}
	if v := os.Getenv("API_SECRET"); v != "" {
		cfg.ApiSecret = v
	}
	if v := os.Getenv("APP_ID"); v != "" {
		cfg.AppID = v
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
