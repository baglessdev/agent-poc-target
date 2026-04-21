# AGENTS.md

Rules every agent (human or LLM) working in this repo must read before making
changes.

## Process

- **One issue → one manifest → one PR.** No mixed-purpose PRs.
- Manifests live in `.claude/manifests/<id>.toml` and are the task contract.
- The agent may only edit files listed in the manifest's `tier2.targets`.
- Before opening a PR, `task verify` must pass.
- Commit message references the manifest id: `<type>(scope): summary [<id>]`.

## Coding rules

- Go formatting: `gofmt -s`. No exceptions.
- `go vet` clean. No `//nolint` to bypass a rule.
- Every new handler has a test in the same package.
- No `panic("TODO")`, no `// TODO` without a linked manifest id.
- `main` stays thin: route registration and lifecycle only.

## What the agent must NOT touch

- `.github/workflows/*`
- `Taskfile.yml`
- `AGENTS.md`, `DESIGN.md`
- `.claude/manifests/*` (the contract is read-only to the coder)

Changes to these files require a human PR, not an agent PR.

## PR expectations

- Title: `<type>(scope): summary`
- Body: a link to the manifest, a one-paragraph change summary, a line
  `Verify: task verify` with pass/fail.
- Bot-authored PRs are branch-prefixed `agent/`.
