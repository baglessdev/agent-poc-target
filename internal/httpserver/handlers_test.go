package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	mux := http.NewServeMux()
	RegisterRoutes(mux)

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
