---
id: T-0993
type: task
nature: remediation
title: flai criteria lists, ticks, and unticks a story's acceptance criteria
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:56Z
updated: 2026-10-06T04:06:48Z
transitions:
  - to: ready
    at: 2026-10-06T03:54:09Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T03:54:10Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T04:06:48Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [flai/cmd/criteria.go, flai/cmd/criteria_test.go, flai/cmd/root.go, flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-0991]
usage:
  source: log
  seconds: 758
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 52
      output: 15203
      cache_read: 2435124
      cache_write: 71673
      cost: 1.2228
---
# T-0993 flai criteria lists, ticks, and unticks a story's acceptance criteria

## Work

Add `flai criteria list <story>` (a read: the numbered criteria, ticked or not, `--json` too) and `flai criteria tick <story> <n>...` and `flai criteria untick <story> <n>...` in `flai/cmd/criteria.go`, registered in `root.go`. Tick and untick call `itemedit.Tick`, then `itemedit.Apply` with the new body, with `--hash`, `--by`, `--autocommit`, `--trailer` as `flai edit` has them, and print the list after. In the guard, `criteria list` is a read a sub-agent may run (`cliReads`); `tick` and `untick` stay refused to sub-agents, the planner, and the orchestrator.

Waits for the itemedit task: it calls `Criteria` and `Tick`.

## Done when

- [ ] `flai criteria list|tick|untick` work on a story, with tests in `flai/cmd/criteria_test.go` for list, tick, untick, a bad number, and an archived story.
- [ ] The guard lets a sub-agent run `flai criteria list` and refuses `flai criteria tick`, with tests.
- [ ] `go test -race ./cmd/ ./internal/guard/` pass for what changed.

## Notes
