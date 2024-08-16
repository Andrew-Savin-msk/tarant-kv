package apiserver

import (
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	LogLevel    string `toml:"log_level"`
	DatabaseURL string `toml:"database_url"`
	SessionKey  string `toml:"session_key"`
}

func ConfigLoad() *Config {
	var configPathEnv = "CONFIG_PATH_DOCKER"
	if runtime.GOOS == "windows" {
		// Connect
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading local.env file whith error:", err)
		}
		configPathEnv = "CONFIG_PATH"
	}

	configPath := os.Getenv(configPathEnv)
	if configPath == "" {
		log.Fatal("config path is not set")
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file does not exists %s", configPath)
	}

	var cfg Config
	_, err = toml.DecodeFile(configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return &cfg
}
