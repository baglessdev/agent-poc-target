package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/baglessdev/agent-poc-target/internal/httpserver"
)

func main() {
	apiKeysEnv := os.Getenv("API_KEYS")
	if apiKeysEnv == "" {
		log.Fatal("API_KEYS environment variable is required")
	}

	apiKeys := parseAPIKeys(apiKeysEnv)
	if len(apiKeys) == 0 {
		log.Fatal("API_KEYS environment variable must contain at least one key")
	}

	mux := http.NewServeMux()
	httpserver.RegisterRoutes(mux)
	handler := httpserver.WithLogging(httpserver.WithRequestID(httpserver.WithAPIKey(apiKeys)(mux)))

	addr := ":8080"
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func parseAPIKeys(s string) []string {
	parts := strings.Split(s, ",")
	keys := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			keys = append(keys, trimmed)
		}
	}
	return keys
}
