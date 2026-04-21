package main

import (
	"encoding/json"
	"net/http"
)

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthz)
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
