package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type PingConfig struct {
	Name    string        `json:"name"`
	Tool    string        `json:"tool"`
	Target  string        `json:"target"`
	Timeout time.Duration `json:"timeout"`
}

type Config struct {
	PingInterval int          `json:"ping_interval"`
	Items        []PingConfig `json:"items"`
}

func LoadConfig(filename string) (*Config, error) {
	// Load environment variables from .env file (if exists)
	godotenv.Load()

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
