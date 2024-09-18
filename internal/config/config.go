package config

type Config struct {
	Server *Server
}

func New() *Config {
	cfg := &Config{
		Server: &Server{
			Host: "local",
			Port: "8080",
		},
	}
	return cfg
}
