package config

import "github.com/wnmay/horo/shared/env"

type Config struct {
	Port               string
	UserManagementAddr string
}

func LoadConfig() *Config {
	return &Config{
		Port:               env.GetString("PORT", "3000"),
		UserManagementAddr: env.GetString("USER_MANAGEMENT_ADDR", "localhost:50051"),
	}
}
