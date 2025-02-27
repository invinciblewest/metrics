package main

import (
	"github.com/invinciblewest/metrics/internal/server/config"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()
	st := storage.NewMemStorage()

	if err := run(cfg.Address, st); err != nil {
		log.Fatal(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)

	log.Println("Server is starting...")
	return http.ListenAndServe(addr, r)
}
