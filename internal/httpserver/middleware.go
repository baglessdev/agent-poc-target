package httpserver

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = generateRequestID()
		}
		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r)
	})
}

func generateRequestID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)
		durationMs := int(duration.Milliseconds())
		requestID := lrw.Header().Get("X-Request-Id")
		authOutcome := lrw.Header().Get("X-Auth-Outcome")

		if authOutcome != "" {
			log.Printf("%s %s %d %dms %s %s", r.Method, r.URL.Path, lrw.statusCode, durationMs, requestID, authOutcome)
		} else {
			log.Printf("%s %s %d %dms %s", r.Method, r.URL.Path, lrw.statusCode, durationMs, requestID)
		}
	})
}

func WithAPIKey(keys []string) func(http.Handler) http.Handler {
	keySet := make(map[string]bool, len(keys))
	for _, k := range keys {
		keySet[k] = true
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isPublicPath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				w.Header().Set("X-Auth-Outcome", "auth=missing")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}` + "\n"))
				return
			}

			if !keySet[apiKey] {
				w.Header().Set("X-Auth-Outcome", "auth=invalid")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"unauthorized"}` + "\n"))
				return
			}

			w.Header().Set("X-Auth-Outcome", "auth=ok")
			next.ServeHTTP(w, r)
		})
	}
}

func isPublicPath(path string) bool {
	return path == "/healthz" || path == "/readyz"
}
