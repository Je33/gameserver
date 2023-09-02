package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	LogLevel  string `envconfig:"LOG_LEVEL"`
	MongoURL  string `envconfig:"MONGODB_URL"`
	MongoDB   string `envconfig:"MONGODB_DATABASE"`
	HTTPAddr  string `envconfig:"HTTP_ADDR"`
	JWTSecret string `envconfig:"JWT_SECRET"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		err := godotenv.Load() // load .env file
		if err != nil {
			log.Fatal(err)
		}
		err = envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
	return &config
}
