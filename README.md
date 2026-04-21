# agent-poc-target

Minimal Go HTTP service used as the target of a two-sandbox coding + review loop.

## Endpoints

- `GET /healthz` — liveness. Returns 200 `{"ok":true}`.

## Run

    go run .
    # then: curl localhost:8080/healthz

## Test

    task test

## Agent rules

See `AGENTS.md`. One issue → one manifest → one PR.
