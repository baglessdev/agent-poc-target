package main

import (
	"log"
	"net/http"

	"github.com/baglessdev/agent-poc-target/internal/httpserver"
)

func main() {
	mux := http.NewServeMux()
	httpserver.RegisterRoutes(mux)
	handler := httpserver.WithLogging(httpserver.WithRequestID(mux))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
