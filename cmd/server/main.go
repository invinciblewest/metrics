package main

import (
	"flag"
	"fmt"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

func main() {
	serverAddr := flag.String("a", "localhost:8080", "server address")

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
