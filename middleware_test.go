package main

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
)

func TestWithRequestID_GeneratesValidID(t *testing.T) {
	handler := withRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("X-Request-Id header not set")
	}

	matched, err := regexp.MatchString("^[a-f0-9]{16}$", requestID)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("X-Request-Id format invalid: got %q want [a-f0-9]{16}", requestID)
	}
}

func TestWithRequestID_PassthroughExisting(t *testing.T) {
	handler := withRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	inboundID := "1234567890abcdef"
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-Id", inboundID)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID != inboundID {
		t.Fatalf("X-Request-Id not preserved: got %q want %q", requestID, inboundID)
	}
}

func TestWithRequestID_OnHealthz(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withRequestID(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("X-Request-Id header not set on /healthz")
	}

	matched, err := regexp.MatchString("^[a-f0-9]{16}$", requestID)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("X-Request-Id format invalid on /healthz: got %q", requestID)
	}
}

func TestWithRequestID_OnReadyz(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withRequestID(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("X-Request-Id header not set on /readyz")
	}

	matched, err := regexp.MatchString("^[a-f0-9]{16}$", requestID)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("X-Request-Id format invalid on /readyz: got %q", requestID)
	}
}

func TestWithRequestID_OnVersion(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withRequestID(mux)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	requestID := rec.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("X-Request-Id header not set on /version")
	}

	matched, err := regexp.MatchString("^[a-f0-9]{16}$", requestID)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("X-Request-Id format invalid on /version: got %q", requestID)
	}
}
