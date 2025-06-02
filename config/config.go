package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Profile       string `yaml:"profile"`
	MySQLUser     string `yaml:"mysql_user"`
	MySQLPassword string `yaml:"mysql_password"`
	MySQLHost     string `yaml:"mysql_host"`
	MySQLDB       string `yaml:"mysql_db"`
}

type yamlConfig struct {
	Dev        Config `yaml:"dev"`
	Production Config `yaml:"production"`
}

func LoadConfig() Config {
	profile := getEnv("PROFILE", "dev")
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		panic("failed to read config.yaml: " + err.Error())
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
	if v := os.Getenv("MYSQL_USER"); v != "" {
		cfg.MySQLUser = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		cfg.MySQLPassword = v
	}
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		cfg.MySQLHost = v
	}
	if v := os.Getenv("MYSQL_DB"); v != "" {
		cfg.MySQLDB = v
	}
	// Debug output
	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
