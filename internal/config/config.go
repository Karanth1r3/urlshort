package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `json:"env"` // struct tags (yaml link),
	//																			(param name if it will be read from env variables)
	// 																			(default value - could be unsafe if config is lost)
	DB         DB         `json:"db"`
	HTTPServer HTTPServer `json:"httpserver"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Must before func name => can throw panics instead of errors (design choice, not language feature)
func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH") // Getting path from env vironment
	if configPath == " " {
		log.Fatal("CONFIG_PATH is not set")
	}
	// If file does not exist => fatal
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config
	// If writing data from file to cfg instance failed => fatal
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
