package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
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

const projectDirName = "gameserver"

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
		currentWorkDirectory, _ := os.Getwd()
		rootPath := projectName.Find([]byte(currentWorkDirectory))

		err := godotenv.Load(string(rootPath) + `/.env`) // load .env file
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
