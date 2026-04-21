# DESIGN.md

## What this repo is

A minimal Go HTTP service, used as the target for an AI coding-agent POC.
Intentionally small so the loop (coder → PR → reviewer → human → merge) is the
interesting artifact, not the service.

## Shape

- `main.go` — lifecycle, listens on `:8080`.
- `handlers.go` — route registration and handler functions.
- `handlers_test.go` — unit tests using `httptest`.
- One `http.ServeMux` built in `main`, populated by `registerRoutes`.
- Responses are JSON via `writeJSON`.

## Invariants

1. Handlers are pure functions of `(http.ResponseWriter, *http.Request)`. No
   globals, no package-level state.
2. Every route has a test.
3. `main` does not contain handler logic.
4. Responses always go through `writeJSON` (or an equivalent explicit helper).

## How to add an endpoint

1. Add a handler to `handlers.go`.
2. Register it in `registerRoutes`.
3. Add a `TestX` in `handlers_test.go`.
4. `task verify`.

## Non-goals

- Database, auth, config, logging framework. This repo is a sandbox for the
  agent loop; do not grow it.
