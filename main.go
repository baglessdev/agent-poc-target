package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withRequestID(mux)

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}
