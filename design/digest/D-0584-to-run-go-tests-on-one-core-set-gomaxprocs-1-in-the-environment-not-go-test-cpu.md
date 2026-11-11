---
id: D-0584
title: "To run Go tests on one core, set GOMAXPROCS=1 in the environment, not go test -cpu 1, which bounds only the test binaries after init and leaves -parallel, the go command, the compiler, and the processes tests start unbounded"
kind: decision
sources: [S-0722]
topics: [go, testing, flai, code]
paths: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, flai/internal/manifest/manifest_tests.go, flai/internal/mcpserver]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0584 To run Go tests on one core, set GOMAXPROCS=1 in the environment, not go test -cpu 1, which bounds only the test binaries after init and leaves -parallel, the go command, the compiler, and the processes tests start unbounded

S-0722 put the close-out's Go tests on one core. Its agent measured the two ways of doing it before choosing. From S-0722's Decisions (`wip/archive/agents/S-0722.md`):

```text
One core is GOMAXPROCS=1 in the environment, not go test -cpu 1. Measured on ./internal/mcpserver -short with -p 1: -cpu 1 took 16.2 s wall at a mean of 1.36 cores and peaks of 3.9 (1 s) and 1.7 (5 s), since -cpu sets GOMAXPROCS only inside the test binaries, after init, and leaves -parallel at 24, the go command, the compiler, and the go processes the tests start unbounded; GOMAXPROCS=1 took 45.4 s at a mean of 0.47 and a peak of 1.00. -cpu 1 with -parallel 24 on one core also starved TestPermissionPromptHoldsARetryAgainUntilTheBound in a whole-module run (a 200 ms hold took 4.0 s end to end, over its 3.2 s limit); GOMAXPROCS=1 makes -parallel 1, and the same run passed. GOMAXPROCS=1 with -parallel 4 took 20.2 s at a mean of 1.02 and peaks of 2.9 and 1.3; not taken, since the story asks for one core and the tests that start git would push past it.
```

A tier has no `env` key, so the `go-test` tier in `system-flow.yaml` runs `env GOMAXPROCS=1 go test -p 1 ...`, and `scripts/test.sh` and `scripts/integration.sh` export it. `flai check`'s `manifest.tests-timeout` sees `go test` through a leading `env` since S-0722 (`design/system/project-manifest.md`).

So a later story that bounds a Go test run's cores sets `GOMAXPROCS`, not `-cpu`; `-cpu 1` with the default `-parallel` can also make a timing test fail that passes at `GOMAXPROCS=1`. `design/system/workflow.md` § "What a close-out's Go tests run (S-0722)" says what the close-out runs; this entry is why it is `GOMAXPROCS` and not `-cpu`.

Source: S-0722.
