---
id: T-3085
type: task
nature: improvement
title: Measure the mcpserver package's CPU before and after the cache, with and without -race, and record it in the narrative
status: done
parent: S-0721
owner: alex
created: 2026-11-11T05:28:33Z
updated: 2026-11-11T05:40:31Z
transitions:
  - to: ready
    at: 2026-11-11T05:36:16Z
    by: agent-S-0721
  - to: in-progress
    at: 2026-11-11T05:36:16Z
    by: agent-S-0721
  - to: done
    at: 2026-11-11T05:40:31Z
    by: agent-S-0721
stream: S-0721
tags: [tests]
touches: [wip/agents/S-0721.md]
after: [T-3083]
usage:
  source: log
  seconds: 255
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2532
      cache_read: 1456868
      cache_write: 37048
      cost: 0.628
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 13
          cache_read: 94870
          cache_write: 22591
          cost: 0.0326
---
# T-3085 Measure the mcpserver package's CPU before and after the cache, with and without -race, and record it in the narrative

## Work

Measure `./internal/mcpserver` as S-0721's Notes say: `go test [-race] -short -c -o <bin> ./internal/mcpserver`, then from `flai/internal/mcpserver` run `time <bin> -test.short -test.count=1 -test.skip '^TestGateway'`. Take the before figure at the story's base commit and the after figure at the commit with T-3083's cache, each with and without `-race`, on the same host and as close together in time as you can, since the host's load moves the figures. Keep the binaries and profiles outside the worktree. Record a table of wall and user CPU for the four runs in the narrative's `## Decisions` with `flai stream log`, never by editing the file. Waits for T-3083, whose cache is what the after figure measures.

If the `-race` run's user CPU is not under half the before figure, profile it with `-test.cpuprofile` and `go tool pprof -top -cum`, say in the narrative what is left, and ask on a thread on S-0721 before going on.

## Done when

- The narrative holds the before and after figures, with and without `-race`, and the commits they were taken at.
- The `-race` run's user CPU is under half the before figure, or a thread on S-0721 says why not.
- Criterion 3 of S-0721 is ticked.

## Notes
