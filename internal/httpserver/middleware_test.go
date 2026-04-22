package httpserver

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func TestWithRequestID_GeneratesValidID(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	RegisterRoutes(mux)
	handler := WithRequestID(mux)

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
	RegisterRoutes(mux)
	handler := WithRequestID(mux)

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
	RegisterRoutes(mux)
	handler := WithRequestID(mux)

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

	handler := WithLogging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logOutput := buf.String()
	durationPattern := regexp.MustCompile(`(\d+)ms`)
	matches := durationPattern.FindStringSubmatch(logOutput)
	if len(matches) < 2 {
		t.Fatalf("duration not found in log output: %q", logOutput)
	}

	durationMs, err := strconv.Atoi(matches[1])
	if err != nil {
		t.Fatalf("failed to parse duration: %v", err)
	}
	if durationMs < 0 {
		t.Fatalf("duration is negative: %dms", durationMs)
	}
}

func TestWithLogging_CapturesStatusCode(t *testing.T) {
	tests := []struct {
		path   string
		status int
	}{
		{"/healthz", 200},
		{"/readyz", 200},
		{"/version", 200},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			var buf bytes.Buffer
			log.SetOutput(&buf)
			defer log.SetOutput(nil)

			mux := http.NewServeMux()
			RegisterRoutes(mux)
			handler := WithLogging(mux)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			logOutput := buf.String()
			expectedStatus := strconv.Itoa(tt.status)
			if !strings.Contains(logOutput, expectedStatus) {
				t.Fatalf("status code %d not found in log output: %q", tt.status, logOutput)
			}
		})
	}
}

func TestWithLogging_WrapsWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithLogging(WithRequestID(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logOutput := buf.String()
	requestIDPattern := regexp.MustCompile(`[a-f0-9]{16}`)
	if !requestIDPattern.MatchString(logOutput) {
		t.Fatalf("request ID not found in log output: %q", logOutput)
	}
}

func TestWithLogging_WrappedByWithRequestID(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithRequestID(WithLogging(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logOutput := buf.String()
	requestIDPattern := regexp.MustCompile(`[a-f0-9]{16}`)
	if !requestIDPattern.MatchString(logOutput) {
		t.Fatalf("request ID not found in log output: %q", logOutput)
	}
}

func TestWithLogging_ContainsAllComponents(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer log.SetOutput(nil)

	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithLogging(WithRequestID(mux))

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	logOutput := buf.String()

	// Check for method
	if !strings.Contains(logOutput, "GET") {
		t.Fatalf("method not found in log output: %q", logOutput)
	}

	// Check for path
	if !strings.Contains(logOutput, "/healthz") {
		t.Fatalf("path not found in log output: %q", logOutput)
	}

	// Check for status code
	if !strings.Contains(logOutput, "200") {
		t.Fatalf("status code not found in log output: %q", logOutput)
	}

	// Check for duration
	durationPattern := regexp.MustCompile(`\d+ms`)
	if !durationPattern.MatchString(logOutput) {
		t.Fatalf("duration not found in log output: %q", logOutput)
	}

	// Check for request ID
	requestIDPattern := regexp.MustCompile(`[a-f0-9]{16}`)
	if !requestIDPattern.MatchString(logOutput) {
		t.Fatalf("request ID not found in log output: %q", logOutput)
	}
}
