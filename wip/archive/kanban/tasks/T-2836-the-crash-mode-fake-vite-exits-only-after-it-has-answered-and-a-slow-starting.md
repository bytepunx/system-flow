---
id: T-2836
type: task
nature: remediation
title: The crash-mode fake vite exits only after it has answered, and a slow-starting fake reproduces the early exit
status: done
parent: S-0650
owner: alex
created: 2026-11-10T06:58:43Z
updated: 2026-11-10T07:24:44Z
transitions:
  - to: ready
    at: 2026-11-10T07:21:00Z
    by: agent-S-0650
  - to: in-progress
    at: 2026-11-10T07:21:01Z
    by: agent-S-0650
  - to: done
    at: 2026-11-10T07:24:44Z
    by: agent-S-0650
stream: S-0650
tags: [flai, test]
touches: [flai/cmd/components_test.go]
usage:
  source: log
  seconds: 223
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 2225
      cache_read: 1431296
      cache_write: 51247
      cost: 0.7408
---
# T-2836 The crash-mode fake vite exits only after it has answered, and a slow-starting fake reproduces the early exit

## Work

First layer: it waits for no task.

- Reproduce first. Give `TestHelperComponentsVite` a start delay, from an environment variable such as `FLAI_FAKE_COMPONENTS_VITE_DELAY`, slept before it listens. Have `TestComponentsServeEndsWhenViteExits` set it longer than the 300 ms crash timer, such as 500 ms. Under the timer as it is, the test then fails as I-0325 says, "exited before it answered", on every run and without load. Record that run in the task's Notes.
- Fix. In `crash` mode, have the fake exit only after it has answered the harness's path once: the handler signals the first answer, such as by closing a channel once, and the crash goroutine waits for it, pauses briefly so the response reaches flai, then prints "the fake vite crashed" and exits 2. Drop the 300 ms timer. Update the helper's comment to say so.
- Keep the delay in the test, so it covers a slow-starting vite on every run.
- If this proves wrong, write down why in the narrative's Decisions before taking another way.

## Done when

- With the delay and before the fix, `TestComponentsServeEndsWhenViteExits` fails with "exited before it answered"; the run is in the task's Notes.
- After the fix, `go test ./cmd -run 'TestComponentsServe' -count 20` passes in `flai/`, with the delay in place.
- `flai test flai/cmd/components_test.go` passes.

## Notes
