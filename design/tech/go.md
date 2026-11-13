---
title: Go
updated: 2026-11-13
status: active
topics: [go]
---

# Go

| | |
|-|-|
| Version | 1.26 (toolchain pinned in `flai/go.mod`) |
| Used in | `flai` |
| Binaries | Two from the one module. `flai`, built from `flai/` and released by GoReleaser. `pilot`, the agent images' entrypoint, built from `flai/pilot`, for Linux only ([ADR-0212](../adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md), S-0507). GoReleaser builds it for linux/amd64 and linux/arm64 into `pilot_<version>_linux_<arch>.tar.gz`, named in the release's `checksums.txt`, and the published bedrock image holds the release's pilot; the bedrock Dockerfile's `tree` target still builds it from the tree with `go build ./pilot` in `flai/` (`flai/images/bedrock/Dockerfile`, [ADR-0295](../adrs/0295-an-installed-flai-pulls-its-agent-images-by-the-digests-a-signed-list-on-its.md)) |
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
| `go test` | built in | Tests, standard library `testing` only, no assertion framework. The close-out runs them at `GOMAXPROCS=1` with `-p 1` and without `-race`; `-race` runs in the manual `race` tier, `flai test race`, and in `make flai-test` (S-0722) |

## Considered

- Rust: better startup and smaller binary but slower iteration and a smaller prompt/CLI ecosystem for this use.
- Node/TypeScript: would share code with the dashboard but requires a runtime and is worse at single-binary distribution.
