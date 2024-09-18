package config

import "time"

type Config struct {
	Server *Server
}

func New() (*Config, error) {
	cfg := &Config{
		Server: &Server{
			Host:    "localhost",
			Port:    "8080",
			Timeout: time.Second * 7,
		},
	}
	return cfg, nil
}
