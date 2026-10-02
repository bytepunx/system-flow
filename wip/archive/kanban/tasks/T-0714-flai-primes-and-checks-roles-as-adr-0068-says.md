---
id: T-0714
type: task
nature: improvement
title: flai primes and checks roles as ADR-0068 says
status: done
parent: S-0196
owner: arobson
created: 2026-10-02T23:29:19Z
updated: 2026-10-02T23:37:06Z
transitions:
  - to: ready
    at: 2026-10-02T23:29:45Z
    by: agent-S-0196
  - to: in-progress
    at: 2026-10-02T23:29:45Z
    by: agent-S-0196
  - to: done
    at: 2026-10-02T23:37:06Z
    by: agent-S-0196
stream: S-0196
tags: []
touches: [flai/internal/conventions, flai/internal/context, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/prime_test.go]
usage:
  source: log
  seconds: 441
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 40
      output: 9760
      cache_read: 1809322
      cache_write: 41276
      cost: 0.8525
---
# T-0714 flai primes and checks roles as ADR-0068 says

## Work

`flai/internal/conventions`: accept `story`, `explore`, `verify`, `orchestrator`, `planner`, and `analyzer` in a convention's `roles`, and warn about any other value. Add a way to ask whether a convention is read by a role: its roles are empty or list it. `flai/internal/context`: `ForStory` keeps only the conventions read by `story`, and `ForRole` those read by the role (empty roles included). `--role` still takes only `explore` or `verify`. `--cat` and the plain listing still print every convention. Update the help in `flai/cmd/prime.go` and the MCP `prime` description in `flai/internal/mcpserver` to match. Waits for nothing: it is the code the other tasks describe and configure.

## Done when

- Tests cover the story pack (no roles, `[story]`, `[verify]` alone left out), each role pack (no roles included, listed role included, other left out), `--cat` printing every convention, and the check (each of the six accepted, another warned)
- `go test` passes for `flai/internal/conventions`, `flai/internal/context`, `flai/cmd` (prime), and `flai/internal/mcpserver` (prime)

## Notes
