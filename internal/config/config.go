package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Srv Server        `toml:"server"`
	VDb ValueDatabase `toml:"value_database"`
	UDb UserDatabase  `toml:"user_database"`
}

type Server struct {
	Port     string        `toml:"port"`
	LogLevel string        `toml:"log_level"`
	TokenTTL time.Duration `toml:"token_ttl" env-default:"1h"`
}

type ValueDatabase struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	DbName   string `toml:"db_name"`
	Port     string `toml:"port"`
}

type UserDatabase struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Host     string `toml:"host"`
	DbName   string `toml:"db_name"`
	Port     string `toml:"port"`
}

func Load() *Config {
	var envParamName = "DOCKER_CONFIG_PATH"
	switch runtime.GOOS {
	case "windows":
		envParamName = "CONFIG_PATH"
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalf("unable to load .env with error: %s", err)
		}
	}

	configPath := os.Getenv(envParamName)
	if configPath == "" {
		fmt.Println("default config path")
		configPath = "config/config.toml"
	}

	_, err := os.Stat(configPath)
	if err != nil {
		log.Fatalf("error with file stats: %s", err)
	}

	var cfg Config
	_, err = toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatalf("unable to load data from config file: %s", err)
	}
	return &cfg
}
