package main

import (
	"fmt"

	"github.com/recommender-system-for-MTUCI/backend/internal/config"
	"github.com/recommender-system-for-MTUCI/backend/internal/transport"
)

func main() {
	//here you app will start
	fmt.Println(config.LoadConfig())
	transport.ReccomendSystem()
}
