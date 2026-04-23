package httpserver

import "net/http"

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/readyz", readyz)
	mux.HandleFunc("/now", now)
	mux.HandleFunc("/version", version)
	mux.HandleFunc("/echo", echo)
	mux.HandleFunc("/sum", sum)
}
