---
id: T-1688
type: task
nature: feature
title: An integration test and a smoke step run a story's agent on OpenCode only where opencode of the known version is installed, and flai serve's refusals for it are tested
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:38Z
updated: 2026-10-16T05:14:38Z
transitions: []
stream: S-0406
tags: [cli]
touches: [flai/tests/integration/opencode_test.go, flai/internal/serve/agents_test.go, scripts/opencode-smoke.sh, scripts/smoke.sh, scripts/README.md]
after: [T-1682, T-1687]
---
# T-1688 An integration test and a smoke step run a story's agent on OpenCode only where opencode of the known version is installed, and flai serve's refusals for it are tested

## Work

- Add `flai/tests/integration/opencode_test.go`: skip under `testing.Short()`, and when `opencode` is not on `PATH` or its version is not the known one T-1682 recorded, with a message naming what is missing. Otherwise start the adapter in a temporary project with a prompt that calls one flai MCP tool, and read its JSON output to the end.
- Add `scripts/opencode-smoke.sh`, called from `scripts/smoke.sh`, which prints `skipped` with the reason as `scripts/gateway-smoke.sh` does, and list it in `scripts/README.md`.
- In `flai/internal/serve/agents_test.go`, test that a story on `opencode` with `roles` is refused at start naming `agent.harnesses.opencode.guard`, and that one without `auto-approve` is refused unless the host sets `--deny-protected`.
- Waits for T-1682 (the known version) and T-1687 (a start that writes the configuration).

## Done when

The test and the smoke step pass with the binary and are skipped, saying so, without it; the refusal tests pass; `make integration` and `make smoke` pass on a host without `opencode`.

## Notes
