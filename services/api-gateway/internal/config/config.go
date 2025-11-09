package config

import "github.com/wnmay/horo/shared/env"

type Config struct {
	Port                     string
	UserManagementAddr       string
	OrderServiceURL          string
	UserManagementServiceURL string
	PaymentServiceURL        string
	ChatServiceURL           string
	RabbitMQURI              string
	ChatAddr       string
}

func LoadConfig() *Config {
	return &Config{
		Port:                     env.GetString("PORT", "8080"),
		UserManagementAddr:       env.GetString("USER_MANAGEMENT_ADDR", "localhost:50051"),
		UserManagementServiceURL: env.GetString("USER_MANAGEMENT_SERVICE_URL", "http://localhost:3003"),
		OrderServiceURL:          env.GetString("ORDER_SERVICE_URL", "http://localhost:3002"),
		PaymentServiceURL:        env.GetString("PAYMENT_SERVICE_URL", "http://localhost:3001"),
		ChatServiceURL:           env.GetString("CHAT_SERVICE_URL", "http://localhost:3004"),
		RabbitMQURI:              env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		ChatAddr: env.GetString("CHAT_ADDR", "localhost:50053"),
	}
}
