package config

import (
	"log"
	"os"
)

type Config struct {
	ServerPort string
}

var AppConfig *Config

func LoadConfig() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	AppConfig = &Config{
		ServerPort: port,
	}

	log.Println("配置加载完成")
}