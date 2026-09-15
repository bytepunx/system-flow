---
title: Go
updated: 2026-09-15
status: active
---

# Go

| | |
|-|-|
| Version | 1.26 (toolchain pinned in `flai/go.mod`) |
| Used in | `flai` |
| Decision | [ADR 0006](../adrs/0006-go-for-the-cli.md) |

## Why

- Single static binary, trivial to install with `go install`, Homebrew, or a GitHub release download. No runtime for the user to manage.
- Excellent standard library for the job: `text/template` for rendering, `os/exec` for git and docker, `encoding/json` for config, `embed` for built-in templates.
- Strong CLI ecosystem (Cobra) and prompt libraries (Charm).
- Fast compile and test cycle for agent-driven development.

## Tooling

| Tool | Version | Use |
|------|---------|-----|
| `golangci-lint` | v2.5 | Lint, config in `flai/.golangci.yaml` in v2 format. A v1 binary cannot read it; install v2 with `go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.5.0` |
| GoReleaser | v2 | Cross-platform release builds on tag `flai/v*` |
| `go test` with `-race` | built in | Tests, standard library `testing` only, no assertion framework |

## Considered

- Rust: better startup and smaller binary but slower iteration and a smaller prompt/CLI ecosystem for this use.
- Node/TypeScript: would share code with the dashboard but requires a runtime and is worse at single-binary distribution.
