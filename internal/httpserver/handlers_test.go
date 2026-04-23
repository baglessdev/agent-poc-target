package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"strings"
	"testing"
	"time"
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

func TestNow(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/now", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	timeVal, ok := body["time"]
	if !ok {
		t.Fatalf("body: missing time field")
	}
	timeFloat, ok := timeVal.(float64)
	if !ok {
		t.Fatalf("body[time]: got type %T want float64", timeVal)
	}
	if timeFloat <= 0 {
		t.Fatalf("body[time]: got %v want positive value", timeFloat)
	}
}

func TestDatetime(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/datetime", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	datetimeVal, ok := body["datetime"]
	if !ok {
		t.Fatalf("body: missing datetime field")
	}
	datetimeStr, ok := datetimeVal.(string)
	if !ok {
		t.Fatalf("body[datetime]: got type %T want string", datetimeVal)
	}
	if _, err := time.Parse(time.RFC3339, datetimeStr); err != nil {
		t.Fatalf("body[datetime]: invalid RFC3339 format: %v", err)
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

func TestSum(t *testing.T) {
	mux := http.NewServeMux()
	RegisterRoutes(mux)

	tests := []struct {
		name       string
		method     string
		url        string
		wantStatus int
		wantBody   map[string]any
	}{
		{
			name:       "valid positive integers",
			method:     http.MethodGet,
			url:        "/sum?a=5&b=3",
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"sum": float64(8)},
		},
		{
			name:       "negative integers",
			method:     http.MethodGet,
			url:        "/sum?a=-10&b=5",
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"sum": float64(-5)},
		},
		{
			name:       "zero values",
			method:     http.MethodGet,
			url:        "/sum?a=0&b=0",
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"sum": float64(0)},
		},
		{
			name:       "large numbers",
			method:     http.MethodGet,
			url:        "/sum?a=1000000&b=2000000",
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"sum": float64(3000000)},
		},
		{
			name:       "missing parameter a",
			method:     http.MethodGet,
			url:        "/sum?b=5",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "parameter required"},
		},
		{
			name:       "missing parameter b",
			method:     http.MethodGet,
			url:        "/sum?a=5",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "parameter required"},
		},
		{
			name:       "missing both parameters",
			method:     http.MethodGet,
			url:        "/sum",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "parameter required"},
		},
		{
			name:       "invalid number a",
			method:     http.MethodGet,
			url:        "/sum?a=abc&b=5",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "invalid number"},
		},
		{
			name:       "invalid number b",
			method:     http.MethodGet,
			url:        "/sum?a=5&b=xyz",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "invalid number"},
		},
		{
			name:       "float input",
			method:     http.MethodGet,
			url:        "/sum?a=5.5&b=3.3",
			wantStatus: http.StatusBadRequest,
			wantBody:   map[string]any{"error": "invalid number"},
		},
		{
			name:       "post method allowed",
			method:     http.MethodPost,
			url:        "/sum?a=5&b=3",
			wantStatus: http.StatusOK,
			wantBody:   map[string]any{"sum": float64(8)},
		},
		{
			name:       "method not allowed",
			method:     http.MethodPut,
			url:        "/sum?a=5&b=3",
			wantStatus: http.StatusMethodNotAllowed,
			wantBody:   map[string]any{"error": "method not allowed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, nil)
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
