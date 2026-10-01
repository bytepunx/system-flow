---
id: S-0183
type: story
nature: remediation
title: "Agents and dev servers get only the environment they need: no host token, the project's own tools first"
status: done
owner: alex
created: 2026-10-01T08:00:32Z
updated: 2026-10-01T08:48:16Z
transitions:
  - to: ready
    at: 2026-10-01T08:32:09Z
    by: alex
  - to: in-progress
    at: 2026-10-01T08:32:36Z
    by: agent-S-0183
  - to: review
    at: 2026-10-01T08:48:00Z
    by: agent-S-0183
  - to: done
    at: 2026-10-01T08:48:16Z
    by: alex
tags: [flai, dashboard]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/dashboard_token.go, flai/internal/host/host.go, flai/internal/host/host_test.go, scripts/env.sh, scripts/flai-test.sh, scripts/install-tools.sh, scripts/flaiover-dev.sh, flaiover/compose.yaml, docs/operators/index.md, docs/operators/settings.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/devex.md, docs/contributors/index.md, design/issues/I-0001-golangci-lint-version-mismatch.md, design/issues/I-0039-scripts-flaiover-dev-sh-points-the-dev-server-at-a-dashboard-token-file-that-flai-no-longer-writes.md, design/issues/I-0044-agents-flai-serve-starts-inherit-the-host-s-address-and-token-so-a-flai-serve-an-agent-runs-by-hand-takes-over-the-operator-s-mcp-servers.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 961
  models:
    - model: claude-opus-5-5
      input: 266
      output: 66668
      cache_read: 13909773
      cache_write: 247654
      cost: 5.8938
---
# S-0183 Agents and dev servers get only the environment they need: no host token, the project's own tools first

## Goal

Three issues share a cause: the environment flai and the scripts hand on is wider or staler than what its receiver needs.

- **I-0044** (three occurrences): `flai serve` starts agents with `append(os.Environ(), ...)` (`flai/internal/serve/agents.go` ~606), so they inherit `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG` (set on purpose in `flai/cmd/serve.go` ~85). A `flai serve` an agent runs by hand then takes over the operator's MCP servers; `scripts/env.sh` keeps an inherited `FLAI_CONFIG`.
- **I-0001** (five occurrences): `scripts/env.sh` (~15) puts the worktree's own `bin/` on PATH but not the main checkout's, where `scripts/install-tools.sh` installs golangci-lint v2. A story worktree has no `bin/`, so its lint finds the host's v1 in `~/go/bin` and fails on the v2 config. It also still prepends `/usr/local/go/bin`, which no longer exists (I-0002, closed).
- **I-0039** (two occurrences): `scripts/flaiover-dev.sh` (~11) defaults `FLAIOVER_TOKEN_FILE` to `.flai-cache/dashboard.token`, and `flaiover/compose.yaml` (~18, ~20) mounts it and sets `PROJECT_DIR`, but since S-0080 flai writes the token to the host's serve directory (`flai/cmd/dashboard_token.go`).

## Acceptance criteria
- [x] Agents `flai serve` starts get no `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, or `FLAI_CONFIG` from its environment; a test checks the started agent's environment
- [x] A `flai serve` run with a `--config` other than the host's does not take over the host's MCP servers or dashboard connection, and says why
- [x] `scripts/env.sh` puts the main checkout's `bin/` on PATH after the worktree's own, so lint in a story worktree runs the pinned golangci-lint; the dead `/usr/local/go/bin` entry is gone; `scripts/flai-test.sh` says which golangci-lint it runs and fails clearly on a v1
- [x] `scripts/flaiover-dev.sh` and `flaiover/compose.yaml` take the dashboard token from where flai writes it (`flai dashboard token --json`, or the serve directory) and drop the stale mount, `PROJECT_DIR`, and comments
- [x] The design (`design/system/flai-cli.md`, `design/system/devex.md`) and the contributor guide describe what an agent's environment holds
- [x] I-0001, I-0039, and I-0044 are closed with what fixed them

## Tasks
- T-0640 flai serve starts agents without the host's address, token, or config
- T-0641 flai serve keeps no MCP servers for a host its config does not name, and says why
- T-0642 Scripts put the main checkout's tools on PATH and flai-test checks golangci-lint's version
- T-0643 The flaiover dev server and compose take the dashboard token from where flai writes it
- T-0644 Design and contributor guide say what an agent's environment holds; close I-0001, I-0039, I-0044

## Notes

- Verified 2026-10-01 by a fresh verifier in the worktree with `FLAI_HOST_URL`, `FLAI_HOST_TOKEN`, and `FLAI_CONFIG` unset: gofmt, go vet, golangci-lint (the main checkout's `bin/golangci-lint` 2.5.0, found from the worktree), the behavior and integration tiers, the template smoke, and `scripts/install-test.sh` pass.
- Two checks fail on files this story does not touch, identically on main: `flai check --strict` reports six warnings (done epics E-0003, E-0010, E-0011, E-0012 not archived; TH-0026 and TH-0032 on archived stories), and `scripts/lint-md.sh` reports MD012 in `wip/archive/agents/S-0175.md` and `wip/archive/agents/S-0180.md`, both from `918cbc4` on main.
- `flaiover/compose.yaml` and `scripts/flaiover-dev.sh` were checked by `docker compose config` and by running `flai dashboard token --json` against a scratch config; the dev server itself was not started.
- That another `--config` takes no dashboard connection holds by structure: flai serve's registry, and so every dashboard it dials, is in the `serve` folder beside its config. No test covers it.
