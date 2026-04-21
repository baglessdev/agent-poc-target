package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"time"
)

func withRequestID(next http.Handler) http.Handler {
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

func withLogging(next http.Handler) http.Handler {
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

		log.Printf("%s %s %d %dms %s", r.Method, r.URL.Path, lrw.statusCode, durationMs, requestID)
	})
}
