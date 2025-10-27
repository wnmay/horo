package config

import "github.com/wnmay/horo/shared/env"

type Config struct {
	GRPCPort               string
	HTTPPort               string
	MongoURI               string
	MongoDBName            string
	UserCollectionName     string
	FirebaseAccountKeyFile string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort:               env.GetString("GRPC_PORT", "50051"),
		HTTPPort:               env.GetString("HTTP_PORT", "8080"),
		MongoURI:               env.GetString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:            env.GetString("MONGO_DB_NAME", "userdb"),
		UserCollectionName:     env.GetString("USER_COLLECTION_NAME", "users"),
		FirebaseAccountKeyFile: env.GetString("FIREBASE_KEY_PATH", "internal/firebase-key.json"),
	}
}
