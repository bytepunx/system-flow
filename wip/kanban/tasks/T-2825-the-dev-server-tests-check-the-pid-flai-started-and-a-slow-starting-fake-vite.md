---
id: T-2825
type: task
nature: remediation
title: The dev server tests check the pid flai started, and a slow-starting fake vite reproduces the missing pid
status: backlog
parent: S-0647
owner: alex
created: 2026-11-10T06:24:06Z
updated: 2026-11-10T06:24:06Z
transitions: []
stream: S-0647
tags: [flai, test]
touches: [flai/internal/components/devserver_test.go, flai/internal/components/devserver.go]
---
# T-2825 The dev server tests check the pid flai started, and a slow-starting fake vite reproduces the missing pid

## Work

First layer: it waits for no task.

- Reproduce first. Give the fake vite in `runFakeVite` a start delay, from an environment variable such as `FLAI_FAKE_VITE_DELAY`, slept before it writes its pid. With a delay longer than 300 ms, `TestStartDevServerStopsViteWhenItsContextEnds` fails as I-0311 says, "the fake vite wrote no pid", on every run and without load. Record that run in the task's Notes.
- Fix. Add a package seam in `devserver.go`, beside `pickPort`, that `launch` calls with `cmd.Process.Pid` once vite starts. The script `newFakeVite` writes `exec`s the test binary, so that pid is the fake's. Have the tests whose vite may be stopped before it runs, the context test and `TestStartDevServerStopsAViteThatDoesNotAnswerInTime`, check that pid with `waitGone` in place of `readPID`. Keep `readPID` where a test needs the fake to have run, such as its args or its child.
- Keep the delay in the context test, longer than its context, so the test covers a vite stopped before it ran on every run.
- If the seam proves wrong, write down why in the narrative's Decisions before taking another way, such as having the script write `$$` before its `exec`.

## Done when

- With the delay and before the fix, the context test fails with "the fake vite wrote no pid"; the run is in the task's Notes.
- After the fix, `go test ./internal/components -run 'TestStartDevServer' -count 20` passes in `flai/`, with the delay in place.
- `flai test flai/internal/components/devserver.go flai/internal/components/devserver_test.go` passes.

## Notes
