package main

import (
	"fmt"
	"github.com/invinciblewest/metrics/internal/server/handlers"
	"github.com/invinciblewest/metrics/internal/storage"
	"net/http"
)

func main() {
	st := storage.NewMemStorage()

	err := run(":8080", st)
	if err != nil {
		panic(err)
	}
}

func run(addr string, st storage.Storage) error {
	r := handlers.GetRouter(st)

	fmt.Println("Server is starting...")
	return http.ListenAndServe(addr, r)
}
