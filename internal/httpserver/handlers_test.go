package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

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
	RegisterRoutes(mux)

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
	tests := []struct {
		name        string
		buildInfo   *debug.BuildInfo
		buildInfoOk bool
		wantVersion string
	}{
		{
			name: "build info present",
			buildInfo: &debug.BuildInfo{
				Main: debug.Module{Version: "v1.2.3"},
			},
			buildInfoOk: true,
			wantVersion: "v1.2.3",
		},
		{
			name:        "build info absent",
			buildInfo:   &debug.BuildInfo{},
			buildInfoOk: true,
			wantVersion: "dev",
		},
		{
			name:        "build info nil",
			buildInfo:   nil,
			buildInfoOk: false,
			wantVersion: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := readBuildInfo
			defer func() { readBuildInfo = orig }()
			readBuildInfo = func() (*debug.BuildInfo, bool) {
				return tt.buildInfo, tt.buildInfoOk
			}

			req := httptest.NewRequest(http.MethodGet, "/version", nil)
			rec := httptest.NewRecorder()
			version(rec, req)

			if rec.Code != http.StatusOK {
				t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}
			if body["version"] != tt.wantVersion {
				t.Fatalf("body[version]: got %v want %v", body["version"], tt.wantVersion)
			}
		})
	}
}

func TestEcho(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:       "valid message",
			method:     http.MethodPost,
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"echoed": "hello"},
		},
		{
			name:       "empty body",
			method:     http.MethodPost,
			body:       "",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "body required"},
		},
		{
			name:       "malformed json",
			method:     http.MethodPost,
			body:       `{bad json}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "invalid json"},
		},
		{
			name:       "missing message field",
			method:     http.MethodPost,
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "message required"},
		},
		{
			name:       "empty message field",
			method:     http.MethodPost,
			body:       `{"message":""}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "message required"},
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       "",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   map[string]any{"error": "method not allowed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/echo", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			for k, v := range tt.wantBody {
				if body[k] != v {
					t.Fatalf("body[%s]: got %v want %v", k, body[k], v)
				}
			}
		})
	}
}

func TestEchoWithAuth(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithAPIKey([]string{"test-key"})(mux)

	tests := []struct {
		name       string
		apiKey     string
		body       string
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:       "missing api key",
			apiKey:     "",
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusUnauthorized,
			wantBody:   map[string]any{"error": "unauthorized"},
		},
		{
			name:       "invalid api key",
			apiKey:     "wrong-key",
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusUnauthorized,
			wantBody:   map[string]any{"error": "unauthorized"},
		},
		{
			name:       "valid api key",
			apiKey:     "test-key",
			body:       `{"message":"hello"}`,
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"echoed": "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/echo", strings.NewReader(tt.body))
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			for k, v := range tt.wantBody {
				if body[k] != v {
					t.Fatalf("body[%s]: got %v want %v", k, body[k], v)
				}
			}
		})
	}
}

func TestVersionWithAuth(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithAPIKey([]string{"test-key"})(mux)

	tests := []struct {
		name       string
		apiKey     string
		wantStatus int
	}{
		{
			name:       "missing api key",
			apiKey:     "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid api key",
			apiKey:     "wrong-key",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "valid api key",
			apiKey:     "test-key",
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/version", nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status: got %d want %d", rec.Code, tt.wantStatus)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			if tt.wantStatus == http.StatusUnauthorized {
				if body["error"] != "unauthorized" {
					t.Fatalf("body: got %v want {error:unauthorized}", body)
				}
			} else {
				if _, ok := body["version"]; !ok {
					t.Fatalf("body: missing version field, got %v", body)
				}
			}
		})
	}
}

func TestHealthzWithAuth(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithAPIKey([]string{"test-key"})(mux)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d (healthz should be public)", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["ok"] != true {
		t.Fatalf("body: got %v want {ok:true}", body)
	}
}

func TestReadyzWithAuth(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	handler := WithAPIKey([]string{"test-key"})(mux)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d (readyz should be public)", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body["ready"] != true {
		t.Fatalf("body: got %v want {ready:true}", body)
	}
}
