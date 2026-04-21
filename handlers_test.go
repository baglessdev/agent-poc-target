package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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

func TestQuicksort(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		want   []int
		status int
	}{
		{
			name:   "unsorted array",
			input:  []int{5, 2, 8, 1, 9},
			want:   []int{1, 2, 5, 8, 9},
			status: http.StatusOK,
		},
		{
			name:   "empty array",
			input:  []int{},
			want:   []int{},
			status: http.StatusOK,
		},
		{
			name:   "single element",
			input:  []int{42},
			want:   []int{42},
			status: http.StatusOK,
		},
		{
			name:   "already sorted",
			input:  []int{1, 2, 3, 4, 5},
			want:   []int{1, 2, 3, 4, 5},
			status: http.StatusOK,
		},
		{
			name:   "reverse sorted",
			input:  []int{9, 7, 5, 3, 1},
			want:   []int{1, 3, 5, 7, 9},
			status: http.StatusOK,
		},
		{
			name:   "duplicates",
			input:  []int{3, 1, 4, 1, 5, 9, 2, 6, 5},
			want:   []int{1, 1, 2, 3, 4, 5, 5, 6, 9},
			status: http.StatusOK,
		},
		{
			name:   "negative numbers",
			input:  []int{-5, 3, -1, 0, 2},
			want:   []int{-5, -1, 0, 2, 3},
			status: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			registerRoutes(mux)

			reqBody := map[string]any{"array": tt.input}
			bodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			req := httptest.NewRequest(http.MethodPost, "/quicksort", bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tt.status {
				t.Fatalf("status: got %d want %d", rec.Code, tt.status)
			}

			var body map[string]any
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("decode: %v", err)
			}

			sortedRaw, ok := body["sorted"]
			if !ok {
				t.Fatalf("missing 'sorted' field in response")
			}

			sortedSlice, ok := sortedRaw.([]any)
			if !ok {
				t.Fatalf("sorted field is not an array")
			}

			sorted := make([]int, len(sortedSlice))
			for i, v := range sortedSlice {
				fv, ok := v.(float64)
				if !ok {
					t.Fatalf("element %d is not a number", i)
				}
				sorted[i] = int(fv)
			}

			if !reflect.DeepEqual(sorted, tt.want) {
				t.Fatalf("sorted: got %v want %v", sorted, tt.want)
			}
		})
	}
}

func TestQuicksortInvalidRequest(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)

	req := httptest.NewRequest(http.MethodPost, "/quicksort", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusBadRequest)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["error"] != "invalid request body" {
		t.Fatalf("error message: got %v want 'invalid request body'", body["error"])
	}
}

func TestQuicksortMethodNotAllowed(t *testing.T) {
	mux := http.NewServeMux()
	registerRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/quicksort", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusMethodNotAllowed)
	}

	var body map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body["error"] != "method not allowed" {
		t.Fatalf("error message: got %v want 'method not allowed'", body["error"])
	}
}
