---
id: T-1713
type: task
nature: improvement
title: flai serve's stream reads its entries through the run's adapter, and a Claude Code run's entries stay the same
status: backlog
parent: S-0411
owner: alex
created: 2026-10-16T05:19:40Z
updated: 2026-10-16T05:19:40Z
transitions: []
stream: S-0411
tags: [cli]
touches: [flai/internal/serve/stream.go, flai/internal/serve/stream_test.go]
---
# T-1713 flai serve's stream reads its entries through the run's adapter, and a Claude Code run's entries stay the same

## Work

- In `flai/internal/serve/stream.go`, move the parsing of Claude Code's stream-json behind a per-harness stream reader chosen by the run's adapter, from its start record, for `Stream`, `PlanStream`, and `RoleStream`.
- Keep the entry shape the dashboard reads as it is, and the Claude Code entries the same: the existing tests pass unchanged.
- Waits for nothing in this story: it changes no behaviour.

## Done when

The stream picks its reader by the run's harness, the Claude Code stream tests pass unchanged, and a run on an unknown harness gives an empty stream, not an error.

## Notes
