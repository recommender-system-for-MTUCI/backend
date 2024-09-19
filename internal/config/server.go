package config

import (
	"fmt"
)

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	//waiting ml for true time
	//Timeout time.Duration `yaml:"timeout"`
	//need timeout for connection server
}

func (s Server) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
