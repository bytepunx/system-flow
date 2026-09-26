---
id: T-0450
type: task
nature: feature
title: flai's serve tests wait for the stub agents they start to end before their temp folders are removed
status: done
parent: S-0122
owner: alex
created: 2026-09-26T06:26:13Z
updated: 2026-09-26T06:32:36Z
transitions:
  - to: ready
    at: 2026-09-26T06:28:59Z
    by: agent-S-0122
  - to: in-progress
    at: 2026-09-26T06:28:59Z
    by: agent-S-0122
  - to: done
    at: 2026-09-26T06:32:36Z
    by: agent-S-0122
stream: S-0122
tags: []
touches: [flai/internal/serve]
---

# T-0450 flai's serve tests wait for the stub agents they start to end before their temp folders are removed

## Work

`newAgentLab` registers a cleanup, run before its temp folders are removed, that ends every stub still held, then waits for each stub it started to have exited and for the launcher to have recorded each agent it waits for as ended. The stub marks its own exit. Today the tests race: a stub or the launcher's goroutine writes into a temp folder while `testing` removes it, and `TempDir RemoveAll cleanup: directory not empty` fails the test, as it failed the flai/v1.18.3 release run.

## Done when

- `go test -race -count=20 ./internal/serve/` passes
- `make flai-test` passes

## Notes
