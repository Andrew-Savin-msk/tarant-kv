package config

import (
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	Srv Server   `toml:"server"`
	Db  Database `toml:"database"`
}

type Server struct {
	Port       string `toml:"port"`
	SessionKey string `toml:"session_key"`
	LogLevel   string `toml:"log_level"`
}

type Database struct {
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
		log.Fatal("enviromental variable doen't exists!")
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
