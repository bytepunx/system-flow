---
id: T-3086
type: task
nature: improvement
title: Run the go-test tier, integration.sh, and test.sh without -race, one package at a time, on one core
status: cancelled
parent: S-0722
owner: alex
created: 2026-11-11T06:05:17Z
updated: 2026-11-11T06:27:07Z
transitions:
  - to: cancelled
    at: 2026-11-11T06:27:07Z
    by: agent-S-0722
stream: S-0722
tags: [tests, close-out]
touches: [system-flow.yaml, scripts/integration.sh, scripts/test.sh, flai/internal/verify/select_test.go]
---
# T-3086 Run the go-test tier, integration.sh, and test.sh without -race, one package at a time, on one core

## Work

First, before changing anything, measure the "before" figures criterion 3 asks for: run the close-out's Go tests as they stand for a change that touches `flai/internal/mcpserver` (the `go-test` tier's command for that package and its importers, then `scripts/integration.sh`), under `/usr/bin/time -v` or with the load sampled, and write the wall time, user CPU, and peak cores into the narrative's `## Decisions`.

Then change the three places the close-out's Go tests are set, so that none runs `-race`, each runs `-p 1`, and each test binary runs on one core:

- `system-flow.yaml`, tier `go-test`: drop `-race`, add `-p 1` and hold the binaries to one core. `go test -cpu 1` sets `GOMAXPROCS` to 1 for every test and needs no environment; a tier has no `env` key, so `GOMAXPROCS=1` would need `scripts/with-env.sh` or `env` in front of `go`, which hides `go test` from the manifest's `-timeout` check. Pick one and say why in `## Decisions`. Write the tier through `flai manifest set tests=...` if it takes the `tests` key, by hand otherwise, and run `flai check` after.
- `scripts/integration.sh`: the same flags on its `go test ... ./...`. Export `GOMAXPROCS=1` there too, so that the compiler and the test binaries it starts stay on one core.
- `scripts/test.sh`: the same, so that `make test` matches.

Keep each `-timeout` as it is: the third task sets it again from the measured figures.

Update `TestRepositoryManifestRunsChangedGoPackagesShortAndIntegrationWhole` in `flai/internal/verify/select_test.go`, which pins the `go-test` tier's flags, to the new ones, and have it fail if `-race` comes back.

Waits for nothing: it is the first layer.

## Done when

- `go-test`, `integration.sh`, and `test.sh` run `go test` with no `-race`, with `-p 1`, and on one core.
- The before figures are in the narrative.
- `flai test flai/internal/verify/select_test.go system-flow.yaml scripts/integration.sh scripts/test.sh` passes, and `flai check --strict` is clean.

## Notes
- 2026-11-11T06:27:07Z: moved to cancelled: duplicates agent-S-0722's set T-3090 to T-3092, written at the same time (TH-1069); kept the story's agent's set and folded in this one's points
