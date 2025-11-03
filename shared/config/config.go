package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv(serviceName string) error {
	// Check if running in Docker
	inDocker := os.Getenv("DOCKER_ENV") == "true" || fileExists("/.dockerenv")

	if inDocker {
		return nil
	}

	// Get the project root
	rootDir, err := findProjectRoot()
	if err != nil {
		return fmt.Errorf("failed to find project root: %w", err)
	}

	// Load shared.env first (lower priority)
	sharedEnvPath := filepath.Join(rootDir, "shared.env")
	if fileExists(sharedEnvPath) {
		if err := godotenv.Load(sharedEnvPath); err != nil {
			log.Printf("Warning: failed to load shared.env: %v", err)
		} else {
			log.Printf("Loaded shared.env from %s", sharedEnvPath)
		}
	}

	// Load service-specific .env
	serviceEnvPath := filepath.Join(rootDir, "services", serviceName, ".env")
	if fileExists(serviceEnvPath) {
		if err := godotenv.Overload(serviceEnvPath); err != nil {
			return fmt.Errorf("failed to load service .env: %w", err)
		}
		log.Printf("Loaded service .env from %s", serviceEnvPath)
	} else {
		if fileExists(".env") {
			if err := godotenv.Overload(".env"); err != nil {
				return fmt.Errorf("failed to load .env: %w", err)
			}
			log.Println("Loaded .env from current directory")
		}
	}

	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Walk up directories looking for go.mod
	for i := 0; i < 10; i++ {
		if fileExists(filepath.Join(dir, "go.mod")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	cwd, _ := os.Getwd()
	return filepath.Dir(cwd), nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}