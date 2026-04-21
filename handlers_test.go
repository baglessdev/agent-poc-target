package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("body: got %v want {ok:true}", body)
	}
}

func TestReadyz(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["ready"] != true {
		t.Fatalf("body: got %v want {ready:true}", body)
	}
}

func TestVersion(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["version"] != "dev" {
		t.Fatalf("body: got %v want {version:dev}", body)
	}
}
