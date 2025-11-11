package config

import (
	"github.com/wnmay/horo/shared/db"
	"github.com/wnmay/horo/shared/env"
)

type ChatMongoConfig struct {
	MongoCommonConfig     *db.MongoConfig
	MessageCollectionName string
	RoomCollectionName    string
}

type Config struct {
	HTTPPort          string
	MongoConfig       *ChatMongoConfig
	RabbitURI         string
	GRPCPort          string
	UserServiceAddr   string
	CourseServiceAddr string
}

const (
	MessageCollectionName = "messages"
	RoomCollectionName    = "rooms"
	dbName                = "chatdb"
)

func LoadConfig() *Config {
	mongoCommonConfig := db.NewMongoDefaultConfig(dbName)
	chatMongoConfig := &ChatMongoConfig{
		MongoCommonConfig:     mongoCommonConfig,
		MessageCollectionName: MessageCollectionName,
		RoomCollectionName:    RoomCollectionName,
	}
	return &Config{
		HTTPPort:          env.GetString("HTTP_PORT", "3005"),
		MongoConfig:       chatMongoConfig,
		RabbitURI:         env.GetString("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/"),
		GRPCPort:          env.GetString("GRPC_PORT", "50053"),
		UserServiceAddr:   env.GetString("USER_SERVICE_ADDR", "localhost:50051"),
		CourseServiceAddr: env.GetString("COURSE_SERVICE_ADDR", "localhost:50052"),
	}
}
