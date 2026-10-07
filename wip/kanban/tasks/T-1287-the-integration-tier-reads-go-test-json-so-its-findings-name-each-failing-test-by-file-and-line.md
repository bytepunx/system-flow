---
id: T-1287
type: task
nature: improvement
title: The integration tier reads go test -json, so its findings name each failing test by file and line
status: backlog
parent: S-0313
owner: alex
created: 2026-10-07T23:39:30Z
updated: 2026-10-07T23:39:30Z
transitions: []
stream: S-0313
tags: [scripts, verify]
touches: [scripts/integration.sh, system-flow.yaml, scripts/README.md]
---
# T-1287 The integration tier reads go test -json, so its findings name each failing test by file and line

## Work

The `integration` tier in `system-flow.yaml` runs `scripts/integration.sh`, which runs `go test -race -count=1 ./...` without `-json`, so the tier is `format: plain`. The `go-test-json` parser flai already has for the `go-test` tier never sees its output.

- Make `scripts/integration.sh` pass its arguments to `go test`, so `make integration` reads as before and the tier can ask for `-json`.
- Change the tier to run `["../scripts/integration.sh", "-json"]` with `dir: flai` and `format: go-test-json`. With `dir: flai`, the parser finds `flai/go.mod` and makes each finding's path relative to the root. It skips the script's `echo` line, which is not an event.
- Update the `integration.sh` row in `scripts/README.md`.
- The `tests` block is not a key `flai manifest set` writes. Edit it by hand, as S-0273 did, and run `flai check --strict` after.

It waits for no task: it shares no path with the `plain()` change.

## Done when

- With a temporary failing test under `flai/`, `flai test --all` reports the integration tier's finding by the test's name, file, and line. Remove the temporary test before the commit.
- `flai test --all` passes the integration tier without it.
- `flai check --strict` is clean.

## Notes

Drafted by the planner for S-0313.
