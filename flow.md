# Agent POC — architecture flow

Three diagrams describing the issue-driven coding-agent POC that produced
the PRs in this repo. Rendered by GitHub natively.

---

## 1. Lifecycle — issue to merge

```mermaid
sequenceDiagram
    actor H as Human
    participant T as Triager<br/>sandbox
    participant C as Coder<br/>sandbox
    participant R as Reviewer<br/>sandbox
    participant GH as GitHub

    H->>GH: Open issue (Agent task template)
    H->>T: triage.sh issue-url
    T->>GH: read issue + clone main
    T->>GH: post manifest draft as comment
    H->>GH: comment "/approve"
    H->>C: coder.sh issue-url
    C->>GH: find /approve + manifest
    C->>C: edit files + task verify
    C->>GH: push branch, open PR (Closes #N)
    H->>R: reviewer.sh pr-url
    R->>GH: read diff + manifest from PR body
    R->>GH: post review (event=COMMENT)
    H->>GH: human approves, merges
```

---

## 2. Topology — three zones

```mermaid
flowchart LR
    subgraph HOST["HOST (your Mac)"]
        KC["macOS Keychain<br/>(Claude OAuth)"]
        POC["distrubuted-agents-poc<br/>wrappers + secrets"]
        AGT["coding-agent-poc<br/>agent roles + prompts"]
    end

    subgraph SB["SANDBOX (mkenv, per run)"]
        RUN["role/run.sh"]
        WS["workspace (target clone)"]
        CL["claude -p"]
    end

    subgraph GHUB["GitHub"]
        TGT["agent-poc-target<br/>(protected main)"]
        AGR["coding-agent-poc (mirror)"]
    end

    KC -->|refresh<br/>per run| POC
    POC -->|mkenv run<br/>+ mount session| SB
    AGT -->|git clone --local<br/>--depth 1| SB
    SB <-->|https :443| TGT
    AGT -.->|git push| AGR
```

---

## 3. Manifest as contract — where it lives across the lifecycle

```mermaid
flowchart LR
    I["Issue body<br/>(human prose)"]
    T["Manifest draft<br/>(issue comment,<br/>agent-manifest-draft marker)"]
    A["Approved<br/>(human posts /approve)"]
    S["manifest.toml<br/>(session dir)"]
    P["PR body<br/>(embedded TOML,<br/>agent-manifest marker)"]
    V["Reviewer verdict<br/>(against manifest)"]

    I -->|triager reads<br/>+ expands| T
    T -->|human decides| A
    A -->|coder extracts| S
    S -->|coder embeds| P
    P -->|reviewer extracts| V
```

---

## Key properties

- **Three sandboxes, three roles** — each run is its own mkenv container, no shared state between roles
- **GitHub-native plumbing** — every artifact is an issue, comment, PR, or review. No parallel tracking system
- **Manifest is the contract** — survives from issue comment → PR body → reviewer context, traceable end-to-end
- **Human gates merges** — reviewer can never approve (`event: COMMENT` is a hard API gate), CI never auto-merges
- **Agent repo cloned per session** — `agent.sha` recorded in session dir for reproducibility; agent code is versioned software
- **Credentials short-lived** — OAuth tokens refreshed from macOS Keychain before every run, never committed, mounted read-only in the sandbox
