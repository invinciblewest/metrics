package main

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/server"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

func main() {
	cfg := server.GetConfig()
	st := storage.NewMemStorage()

	err := run(cfg.Address, st)
	if err != nil {
		panic(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)

	fmt.Println("Server is starting...")
	return http.ListenAndServe(addr, r)
}
