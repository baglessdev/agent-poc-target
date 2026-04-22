package httpserver

import (
	"encoding/json"
	"io"
	"net/http"
	"runtime/debug"
)

var readBuildInfo = debug.ReadBuildInfo

func healthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func readyz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ready": true})
}

func version(w http.ResponseWriter, _ *http.Request) {
	v := "dev"
	if info, ok := readBuildInfo(); ok && info.Main.Version != "" {
		v = info.Main.Version
	}
	writeJSON(w, http.StatusOK, map[string]any{"version": v})
}

func echo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}

	if r.Body == nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "body required"})
		return
	}

	var req struct {
		Message string `json:"message"`
	}

	if err := readJSON(r, &req); err != nil {
		if err == io.EOF {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "body required"})
			return
		}
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json"})
		return
	}

	if req.Message == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "message required"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"echoed": req.Message})
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}

	writeJSON(w, http.StatusOK, map[string]any{"message": "hello, " + name})
}

func readJSON(r *http.Request, dest any) error {
	return json.NewDecoder(r.Body).Decode(dest)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
