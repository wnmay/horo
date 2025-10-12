package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustOpen() *gorm.DB {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASS", "postgres")
	name := getenv("DB_NAME", "postgres")
	ssl  := getenv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		host, port, user, pass, name, ssl)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil { panic(err) }
	return db
}

func getenv(k, def string) string { if v := os.Getenv(k); v != "" { return v }; return def }
