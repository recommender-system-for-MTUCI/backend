package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server *Server
	JWT    *JWT
}

func NewConfig() (*Config, error) {
	//wd, err := os.Getwd()
	//if err != nil {
	//	return nil, err
	//}
	cfg := &Config{
		Server: &Server{
			Host: "0.0.0.0",
			Port: "8080",
			//Timeout: time.Second * 7,
		},
		JWT: &JWT{
			AccsesTokenLifetime:  5,
			RefreshTokenLifetime: 10000,
			PublicKeyPath:        "/home/relationskatie/recjmmendSystem/certs/public.pem",
			PrivateKeyPath:       "/home/relationskatie/recjmmendSystem/certs/key.pem",
		},
	}
	cfg, err := loadConfig(cfg)
	if err != nil {
		log.Fatal("failed to load config")
	}

	return cfg, err
}

func loadConfig(cfg *Config) (*Config, error) {
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
