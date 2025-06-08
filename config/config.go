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
	ClientKey     string `yaml:"client_key"`
	ClientSecret  string `yaml:"client_secret"`
	AppID         string `yaml:"app_id"`
	AppSecret     string `yaml:"app_secret"`
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
	// config client credentials
	if v := os.Getenv("CLIENT_KEY"); v != "" {
		cfg.ClientKey = v
	}
	if v := os.Getenv("CLIENT_SECRET"); v != "" {
		cfg.ClientSecret = v
	}
	if v := os.Getenv("APP_ID"); v != "" {
		cfg.AppID = v
	}
	if v := os.Getenv("APP_SECRET"); v != "" {
		cfg.AppSecret = v
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
