---
id: D-0529
title: "A task that adds or changes real-git tests, which skip under flai test's -short, runs them with scripts/integration.sh -run '<tests>' before it closes"
kind: decision
sources: [S-0609]
topics: [testing, go, flai]
paths: [scripts/integration.sh, system-flow.yaml, design/conventions/code-quality.md]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0529 A task that adds or changes real-git tests, which skip under flai test's -short, runs them with scripts/integration.sh -run '<tests>' before it closes

`flai test` and the close-out's `go-test` tier run the Go tests with `-short`, and the tests that need real git or the monorepo skip there; they are the integration tier's own (`design/conventions/code-quality.md`, S-0548). So a task that changes the real-git acceptance tests gets no run of them from `flai test`.

S-0609 ran them through the project's integration entry point with a filter: "The real-git acceptance tests skip under `flai test`'s `-short`, so each Go task's were run through `scripts/integration.sh -run '<tests>'`, the project's integration entry point with a filter: all passed" (wip/archive/agents/S-0609.md, Decisions).

`scripts/integration.sh` passes its arguments to `go test -race -count=1 -timeout 20m -skip '^TestGateway' ... ./...`, so a `-run` filter runs only the named tests, under race and the 20-minute timeout, and never the paid gateway checks (D-0451). Use it rather than a hand `go test`.

Source: S-0609.
