package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *Server
}

func NewConfig() (*Config, error) {
	cfg := &Config{
		Server: &Server{
			Host: "localhost",
			Port: "8080",
			//Timeout: time.Second * 7,
		},
	}
	return cfg, nil
}

func LoadConfig() (*Config, error) {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal("Failed to create config")
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Failed to get home dir, using default config")
		return cfg, err
	}
	fullPath := filepath.Join(homeDir, "recjmmendSystem/config/server.yaml")
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		log.Printf("Failed to found file, using default config")
		return cfg, err
	}

	file, err := os.ReadFile(fullPath)
	if err != nil {
		log.Printf("Failed read file, using default config")
		return cfg, err
	}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		log.Printf("Failed read file, using default config")
		return cfg, err
	}

	return cfg, err
}
