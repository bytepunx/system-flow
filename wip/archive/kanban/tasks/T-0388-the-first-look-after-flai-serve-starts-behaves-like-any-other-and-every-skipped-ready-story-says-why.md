---
id: T-0388
type: task
nature: remediation
title: The first look after flai serve starts behaves like any other, and every skipped ready story says why
status: done
parent: S-0112
owner: alex
created: 2026-09-24T06:19:57Z
updated: 2026-09-24T06:22:38Z
transitions:
  - to: ready
    at: 2026-09-24T06:20:01Z
    by: agent-S-0112
  - to: in-progress
    at: 2026-09-24T06:20:01Z
    by: agent-S-0112
  - to: done
    at: 2026-09-24T06:22:38Z
    by: agent-S-0112
stream: S-0112
tags: []
touches: [flai/internal/serve]
---
# T-0388 The first look after flai serve starts behaves like any other, and every skipped ready story says why

## Work

- Remove the launcher's `first` and `stale` from `flai/internal/serve/agents.go`: the first look settles orphans and starts what any later look would. `serve/agents.json` (`startedBefore`) is what keeps a restart from starting a story twice.
- Give every ready story the launcher skips a reason in `waiting`: tried since it entered ready (with when and how it ended), names no harness with no command set, someone attending, the in-progress limit. Log `agent not started` for a story when its reason changes, not at every look.
- Update the comments in `agents.go` and `serve.go` that say a restart starts nothing.
- Replace the tests that assert a restart starts nothing with ones that assert a restart starts nothing twice, and add a test that each skip is stated.

## Done when

- `go test -race ./internal/serve/...` passes with the new tests, which fail against the old code.

## Notes
