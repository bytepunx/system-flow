---
id: T-0842
type: task
nature: improvement
title: scripts/flai-test.sh runs the Go tests once, the full run with the race detector, and still runs flaiover's unit tests
status: cancelled
parent: S-0267
owner: alex
created: 2026-10-05T00:18:11Z
updated: 2026-10-05T00:20:01Z
transitions:
  - to: cancelled
    at: 2026-10-05T00:20:01Z
    by: agent-S-0267
stream: S-0267
tags: [scripts, tests]
touches: [scripts/flai-test.sh, scripts/test.sh, scripts/flaiover-unit.sh]
---
# T-0842 scripts/flai-test.sh runs the Go tests once, the full run with the race detector, and still runs flaiover's unit tests

## Work

`scripts/flai-test.sh` calls `scripts/test.sh` and then `scripts/integration.sh`. The first runs `go test -race -short -count=1 ./...` and then flaiover's vitest. The second runs `go test -race -count=1 ./...` over the same packages, a superset of the short run. Make the sequence run the Go tests once:

- Move the flaiover block out of `scripts/test.sh` into a script of its own, such as `scripts/flaiover-unit.sh`. That block builds `bin/flai` when it is missing or stale and runs `pnpm test:unit` when `flaiover/node_modules` is present. `scripts/test.sh` calls the new script, so `make test` behaves as before.
- In `scripts/flai-test.sh`, keep gofmt, vet, and golangci-lint. Replace the call to `scripts/test.sh` with `scripts/integration.sh`, then the flaiover unit script, then `scripts/smoke.sh`, as now.
- Leave `scripts/integration.sh` and `scripts/smoke.sh` alone, so `make integration` and `make smoke` do not change.

If you find a simpler shape that keeps `make test` and `make integration` as they are, use it, and update this task's touches with `flai edit`. One example is an environment variable that `scripts/test.sh` reads to skip its Go run.

The dropped short run is not a gap in the tiers. The full run includes every test `-short` runs, so a failing behaviour test still fails `flai-test.sh`. Record this in the narrative's `## Decisions` against `code-quality.md`'s rule that each tier runs only once the prior one passes.

It waits for no task.

## Done when

- [ ] `scripts/flai-test.sh` runs `go test` once per run, `-race -count=1` without `-short`, as its output and a `grep` of the scripts it calls show
- [ ] with `flaiover/node_modules` present, `scripts/flai-test.sh` still runs vitest; without it, the script skips vitest as before
- [ ] `make test` still runs `go test -race -short -count=1 ./...` and then vitest, and `make integration` still runs the full `go test -race -count=1 ./...`
- [ ] `scripts/flai-test.sh` passes in the story's worktree, and its time before and after the change is recorded in the narrative

## Notes
- 2026-10-05T00:20:01Z: moved to cancelled: duplicate of T-0846, written by the story agent at the same time
