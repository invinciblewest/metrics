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
	mux := http.NewServeMux()

	mux.Handle("/update/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		handlers.UpdateMetricHandler(w, r, st)
	}))

	fmt.Println("Server is starting...")
	return http.ListenAndServe(addr, mux)
}
