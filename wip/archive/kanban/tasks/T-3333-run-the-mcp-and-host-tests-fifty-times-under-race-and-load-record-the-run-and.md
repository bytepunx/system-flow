---
id: T-3333
type: task
nature: remediation
title: Run the MCP and host tests fifty times under race and load, record the run, and close the issue
status: done
parent: S-0779
owner: alex
created: 2026-11-13T20:33:58Z
updated: 2026-11-13T22:05:39Z
transitions:
  - to: ready
    at: 2026-11-13T20:48:23Z
    by: agent-S-0779
  - to: in-progress
    at: 2026-11-13T20:48:23Z
    by: agent-S-0779
  - to: done
    at: 2026-11-13T22:05:39Z
    by: agent-S-0779
stream: S-0779
tags: [flai]
touches: [design/issues/summary.md]
after: [T-3329, T-3331, T-3332, T-3337, T-3338, T-3339]
usage:
  source: log
  seconds: 4636
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 112
      output: 33754
      cache_read: 14610691
      cache_write: 179356
      cost: 5.0325
---
# T-3333 Run the MCP and host tests fifty times under race and load, record the run, and close the issue

## Work

From `flai/`, run `go test -race -count=50 -run 'MCP|TheHost' ./cmd/` beside the whole package's tests, `go test -race ./cmd/` at the same time, or under like load, so that other parallel tests open listeners and connections while it runs, as they did in CI. When it fails, find which test lost which port, fix it on the task that owns that test, and run it again.

Record the command, the load beside it, how long it took, and its result in the narrative with `flai stream log`. Then close the issue T-3329 recorded with `flai issue close`, saying what fixed it: the tests decide for themselves which ports are free and held, and name this story.

Waits for T-3331 and T-3332, whose tests it runs, and for T-3329, whose issue it closes.

## Done when

- `go test -race -count=50 -run 'MCP|TheHost' ./cmd/` passed beside the whole package's tests or under like load, and the narrative records the run.
- The issue is closed with what fixed it, and `design/issues/summary.md` no longer lists it as open.

## Notes
