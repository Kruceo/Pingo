package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

type PingConfig struct {
	Name    string `json:"name"`
	Tool    string `json:"tool"`
	Target  string `json:"target"`
	Timeout int    `json:"timeout"` // Timeout in milliseconds
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

func SaveConfig(filename string, cfg *Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(filename)

	perm := os.FileMode(0o644)
	if st, err := os.Stat(filename); err == nil {
		perm = st.Mode()
	}

	tmpFile, err := os.CreateTemp(dir, ".pingo-config-*.json")
	if err != nil {
		return err
	}
	tmpName := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpName)
	}()

	if err := tmpFile.Chmod(perm); err != nil {
		_ = tmpFile.Close()
		return err
	}

	if _, err := bytes.NewBuffer(data).WriteTo(tmpFile); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return err
	}
	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpName, filename)
}

func ValidatePingConfig(item PingConfig) error {
	if strings.TrimSpace(item.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(item.Tool) == "" {
		return fmt.Errorf("tool is required")
	}
	if strings.TrimSpace(item.Target) == "" {
		return fmt.Errorf("target is required")
	}
	if item.Timeout <= 0 {
		return fmt.Errorf("timeout must be > 0")
	}
	return nil
}
