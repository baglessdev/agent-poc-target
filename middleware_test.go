package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
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

func TestWithLogging_DurationNonNegative(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	handler := withLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	matched, err := regexp.MatchString(`\d+ms`, logLine)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("log line missing duration: %q", logLine)
	}

	re := regexp.MustCompile(`(\d+)ms`)
	matches := re.FindStringSubmatch(logLine)
	if len(matches) < 2 {
		t.Fatalf("duration not found in log: %q", logLine)
	}
}

func TestWithLogging_StatusCodeHealthz(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withLogging(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "200") {
		t.Fatalf("log line missing status 200: %q", logLine)
	}
}

func TestWithLogging_StatusCodeReadyz(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withLogging(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "200") {
		t.Fatalf("log line missing status 200: %q", logLine)
	}
}

func TestWithLogging_StatusCodeVersion(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withLogging(mux)

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	if !strings.Contains(logLine, "200") {
		t.Fatalf("log line missing status 200: %q", logLine)
	}
}

func TestWithLogging_WrapsWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withLogging(withRequestID(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	matched, err := regexp.MatchString(`[a-f0-9]{16}`, logLine)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("log line missing request ID: %q", logLine)
	}
}

func TestWithLogging_WrappedByWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withRequestID(withLogging(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()
	matched, err := regexp.MatchString(`[a-f0-9]{16}`, logLine)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Fatalf("log line missing request ID: %q", logLine)
	}
}

func TestWithLogging_AllComponents(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	registerRoutes(mux)
	handler := withLogging(withRequestID(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logLine := buf.String()

	if !strings.Contains(logLine, "GET") {
		t.Errorf("log line missing method: %q", logLine)
	}
	if !strings.Contains(logLine, "/healthz") {
		t.Errorf("log line missing path: %q", logLine)
	}
	if !strings.Contains(logLine, "200") {
		t.Errorf("log line missing status: %q", logLine)
	}
	matched, err := regexp.MatchString(`\d+ms`, logLine)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Errorf("log line missing duration: %q", logLine)
	}
	matched, err = regexp.MatchString(`[a-f0-9]{16}`, logLine)
	if err != nil {
		t.Fatalf("regex error: %v", err)
	}
	if !matched {
		t.Errorf("log line missing request ID: %q", logLine)
	}
}
