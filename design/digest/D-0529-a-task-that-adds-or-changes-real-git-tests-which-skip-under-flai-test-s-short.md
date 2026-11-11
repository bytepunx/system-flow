---
id: D-0529
title: "A task that adds or changes real-git tests, which skip under flai test's -short, runs them with scripts/integration.sh -run '<tests>' before it closes"
kind: decision
sources: [S-0609, S-0702]
topics: [testing, go, flai]
paths: [scripts/integration.sh, system-flow.yaml, design/conventions/code-quality.md, flai/internal/mergequeue/batch_test.go]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0529 A task that adds or changes real-git tests, which skip under flai test's -short, runs them with scripts/integration.sh -run '<tests>' before it closes

`flai test` and the close-out's `go-test` tier run the Go tests with `-short`, and the tests that need real git or the monorepo skip there; they are the integration tier's own (`design/conventions/code-quality.md`, S-0548). So a task that changes the real-git acceptance tests gets no run of them from `flai test`.

S-0609 ran them through the project's integration entry point with a filter: "The real-git acceptance tests skip under `flai test`'s `-short`, so each Go task's were run through `scripts/integration.sh -run '<tests>'`, the project's integration entry point with a filter: all passed" (wip/archive/agents/S-0609.md, Decisions).

S-0702 did the same for the merge queue's batch tests in `flai/internal/mergequeue`: "T-3009's sub-agent ran the mergequeue batch tests through scripts/integration.sh with a -run filter, since the test tool's -short run skips real-git tests; the close-out's integration tier runs them all" (wip/archive/agents/S-0702.md, Decisions).

`scripts/integration.sh` passes its arguments to `go test`, with `-skip '^TestGateway'` and `./...`, so a `-run` filter runs only the named tests and never the paid gateway checks (D-0451). Since S-0722 it runs on one core without `-race` under `-timeout 10m`, and with `--race` first under `-race` on every core with `-timeout 20m` (D-0585). Use it rather than a hand `go test`.

Sources: S-0609, S-0702.
