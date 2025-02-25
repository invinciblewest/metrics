package main

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

type ServerConfig struct {
	Address string `env:"ADDRESS"`
}

func main() {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	var config ServerConfig
	if err := env.Parse(&config); err != nil {
		panic(err)
	}

	if config.Address != "" {
		serverAddr = &config.Address
	}

	st := storage.NewMemStorage()

	err := run(*serverAddr, st)
	if err != nil {
		panic(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)

	fmt.Println("Server is starting...")
	return http.ListenAndServe(addr, r)
}
