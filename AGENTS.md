# AGENTS.md

Binding contract every contributing agent (human or LLM) MUST read before
making changes in this repo.

---

## Process

- **One issue → one manifest → one PR.** No mixed-purpose PRs.
- Manifests live in issue comments (posted by the triager), not in this repo's
  history.
- The agent may only edit files listed in the manifest's `tier2.targets`.
- Before opening a PR, every step in the Verify section below must pass locally.
- Commit message references the manifest id: `<type>(scope): summary`.
- Agent-authored branches are prefixed `agent/<timestamp>-<task-id>`.

---

## Coding rules

- Go formatting: `gofmt -s`. No exceptions.
- `go vet` clean. No `//nolint` to bypass a methodology rule.
- Every new handler, exported function, or public method ships with a test
  in the same package.
- No `panic("TODO")`, no `// TODO` or `// FIXME` without a linked issue id.
- No global mutable state in domain code.
- `main.go` stays thin: route registration and lifecycle only.
- No `any` / `interface{}` in domain code — reach for concrete types.

---

## Forbidden paths

Files the agent must NEVER edit. Reviewer will flag any PR that touches these.

- `.github/workflows/*` — CI config is human-owned
- `Taskfile.yml` — build config is human-owned
- `AGENTS.md`, `DESIGN.md` — conventions are human-owned
- `.agent/*` — agent adoption artifacts (if added later)

Changes to any of the above require a human-authored PR, not an agent PR.

---

## Verify

Canonical commands that must pass before a PR opens. The triager picks from
this list when building a manifest.

Run from the repo root:

- `task lint`
- `task build`
- `task test`

Or equivalently, `task verify` (which chains the three above).

---

## Invariants

1. Handlers are pure functions of `(http.ResponseWriter, *http.Request)`. No
   globals, no package-level state.
2. Every route has a test covering at least one happy-path case.
3. `main` does not contain handler logic.
4. Responses always go through `writeJSON` (or an equivalent explicit helper).

---

## PR expectations

- Title: `<type>(scope): summary`
- Body: link to the manifest (embedded in the body for the reviewer to read),
  one-paragraph change summary, a line `Verify: <pass/fail>` with the outcome.
- Reviewer never approves. A human approves and merges.
