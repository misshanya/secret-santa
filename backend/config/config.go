package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort  int
	DatabaseURL string
	JWTSecret   string
}

var (
	config  *Config
	onceCfg sync.Once
)

// loadConfig loads environment variables from a .env file (or directly from env vars) and initializes the config.
func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, trying to use env variables")
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		log.Fatal("Error getting server port variable")
	}
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error getting database url variable")
	}
	JWTSecret := os.Getenv("JWT_SECRET")
	if JWTSecret == "" {
		log.Fatal("Error getting jwt secret variable")
	}

	serverPortInt, err := strconv.Atoi(serverPort)
	if err != nil {
		log.Fatal("Error parsing server port to int")
	}

	config = &Config{
		ServerPort:  serverPortInt,
		DatabaseURL: databaseURL,
		JWTSecret:   JWTSecret,
	}
}

// GetConfig returns the application configuration, loading it only once.
func GetConfig() *Config {
	onceCfg.Do(loadConfig)
	return config
}
