---
id: S-0183
type: story
nature: remediation
title: "Agents and dev servers get only the environment they need: no host token, the project's own tools first"
status: backlog
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:00:32Z
transitions: []
tags: [flai, dashboard]
touches: [flai/internal/serve/agents.go, flai/cmd/serve.go, scripts/env.sh, scripts/flaiover-dev.sh, flaiover/compose.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0183 Agents and dev servers get only the environment they need: no host token, the project's own tools first

## Goal

Three issues share a cause: the environment flai and the scripts hand on is wider or staler than what its receiver needs.

- **I-0044** (three occurrences): `flai serve` starts agents with `append(os.Environ(), ...)` (`flai/internal/serve/agents.go` ~606), so they inherit `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG` (set on purpose in `flai/cmd/serve.go` ~85). A `flai serve` an agent runs by hand then takes over the operator's MCP servers; `scripts/env.sh` keeps an inherited `FLAI_CONFIG`.
- **I-0001** (five occurrences): `scripts/env.sh` (~15) puts the worktree's own `bin/` on PATH but not the main checkout's, where `scripts/install-tools.sh` installs golangci-lint v2. A story worktree has no `bin/`, so its lint finds the host's v1 in `~/go/bin` and fails on the v2 config. It also still prepends `/usr/local/go/bin`, which no longer exists (I-0002, closed).
- **I-0039** (two occurrences): `scripts/flaiover-dev.sh` (~11) defaults `FLAIOVER_TOKEN_FILE` to `.flai-cache/dashboard.token`, and `flaiover/compose.yaml` (~18, ~20) mounts it and sets `PROJECT_DIR`, but since S-0080 flai writes the token to the host's serve directory (`flai/cmd/dashboard_token.go`).

## Acceptance criteria
- [ ] Agents `flai serve` starts get no `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, or `FLAI_CONFIG` from its environment; a test checks the started agent's environment
- [ ] A `flai serve` run with a `--config` other than the host's does not take over the host's MCP servers or dashboard connection, and says why
- [ ] `scripts/env.sh` puts the main checkout's `bin/` on PATH after the worktree's own, so lint in a story worktree runs the pinned golangci-lint; the dead `/usr/local/go/bin` entry is gone; `scripts/flai-test.sh` says which golangci-lint it runs and fails clearly on a v1
- [ ] `scripts/flaiover-dev.sh` and `flaiover/compose.yaml` take the dashboard token from where flai writes it (`flai dashboard token --json`, or the serve directory) and drop the stale mount, `PROJECT_DIR`, and comments
- [ ] The design (`design/system/flai-cli.md`, `design/system/devex.md`) and the contributor guide describe what an agent's environment holds
- [ ] I-0001, I-0039, and I-0044 are closed with what fixed them

## Tasks

## Notes
