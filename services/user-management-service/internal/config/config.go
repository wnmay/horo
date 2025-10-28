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

const (
	dbName             = "userdb"
	UserCollectionName = "users"
)

func LoadConfig() *Config {
	return &Config{
		GRPCPort:               env.GetString("GRPC_PORT", "50051"),
		HTTPPort:               env.GetString("HTTP_PORT", "8080"),
		MongoURI:               env.GetString("MONGO_URI", "mongodb://localhost:27017"),
		MongoDBName:            dbName,
		UserCollectionName:     UserCollectionName,
		FirebaseAccountKeyFile: env.GetString("FIREBASE_KEY_PATH", "internal/firebase-key.json"),
	}
}
