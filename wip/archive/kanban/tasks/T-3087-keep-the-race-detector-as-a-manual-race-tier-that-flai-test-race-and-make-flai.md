---
id: T-3087
type: task
nature: improvement
title: Keep the race detector as a manual race tier that flai test race and make flai-test run
status: cancelled
parent: S-0722
owner: alex
created: 2026-11-11T06:05:28Z
updated: 2026-11-11T06:27:07Z
transitions:
  - to: cancelled
    at: 2026-11-11T06:27:07Z
    by: agent-S-0722
stream: S-0722
tags: [tests, close-out]
touches: [scripts/race.sh, system-flow.yaml, scripts/flai-test.sh, Makefile, flai/internal/verify/select_test.go]
after: [T-3086]
---
# T-3087 Keep the race detector as a manual race tier that flai test race and make flai-test run

## Work

Add `scripts/race.sh`: the full Go suite under `-race`, as `scripts/integration.sh` ran it before the first task, `go test -race -count=1 -timeout 20m -skip '^TestGateway' "$@" ./...` in `flai/`, with the same flaiover install in a worktree and a header comment saying it is the `race` tier's and why no close-out runs it. POSIX `sh` under `set -eu`, ShellCheck clean.

Add a tier `race` to `system-flow.yaml`'s `tests`, next to `gateway`: `command: ["../scripts/race.sh", "-json"]`, `dir: flai`, `format: go-test-json`, `all_only: true`, `manual: true`, with `paths: ["scripts/race.sh"]`. It has no `queue` and no `when`, so no close-out and no merge queue batch runs it, and `--all` leaves it out as it leaves `gateway` out.

Have `scripts/flai-test.sh` run `race.sh` in place of `integration.sh`, since it runs the same full suite and adds the detector, so `make flai-test` runs the suite once and still under `-race`; update its header comment. Add a `race` target to the `Makefile` that calls `scripts/race.sh`.

Add `TestRepositoryManifestRunsTheRaceDetectorOnlyByName` to `flai/internal/verify/select_test.go`, beside `TestRepositoryManifestRunsTheGatewayChecksOnlyByName`: `SelectTier(tiers, "race")` finds a manual tier that is not `queue` and runs `-race`, and neither a paths run, an `--all` run, nor a story's selection picks it. Not gated under `-short`.

Waits for the first task: both write the `tests` list in `system-flow.yaml` and `select_test.go`, and the race tier keeps the command the first task takes out of `integration.sh`.

## Done when

- `flai test race` runs the full suite under `-race` and passes.
- `make flai-test` runs `race.sh`; no close-out, `--all` run, or batch selects `race`, and the new test says so.
- `flai test scripts/race.sh scripts/flai-test.sh system-flow.yaml flai/internal/verify/select_test.go` passes, and `flai check --strict` is clean.

## Notes
- 2026-11-11T06:27:07Z: moved to cancelled: duplicates agent-S-0722's set T-3090 to T-3092, written at the same time (TH-1069); kept the story's agent's set and folded in this one's points
