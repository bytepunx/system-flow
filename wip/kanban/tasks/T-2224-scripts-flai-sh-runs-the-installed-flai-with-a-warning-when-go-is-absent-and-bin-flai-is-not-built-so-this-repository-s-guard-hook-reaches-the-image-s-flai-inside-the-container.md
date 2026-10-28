---
id: T-2224
type: task
nature: remediation
title: scripts/flai.sh runs the installed flai, with a warning, when go is absent and bin/flai is not built, so this repository's guard hook reaches the image's flai inside the container
status: backlog
parent: S-0518
owner: alex
created: 2026-10-28T13:33:04Z
updated: 2026-10-28T13:33:51Z
transitions: []
stream: S-0518
tags: [cli]
touches: [scripts/flai.sh, scripts/README.md]
---
# T-2224 scripts/flai.sh runs the installed flai, with a warning, when go is absent and bin/flai is not built, so this repository's guard hook reaches the image's flai inside the container

## Work

In `scripts/flai.sh`, before the build, check `command -v go`. With `go` present nothing changes: build when `bin/flai` is missing or stale, then `exec "$ROOT/bin/flai"`. Without `go`: when `bin/flai` is executable, run it, warning on standard error when it is stale, since it cannot be rebuilt here; else when `flai` is on `PATH`, print one warning, `go not found; running the installed flai` with its path from `command -v flai`, to standard error and `exec flai "$@"`; else print `go not found and no flai on PATH` and exit 127, so the guard hook of the next task refuses the call rather than passing it. Keep the script POSIX `sh` under `set -eu`, quoted as `tooling.md` asks, and keep `scripts/flai-build.sh` as it is. Update the `flai.sh` row of `scripts/README.md` to say what it runs without `go`. Note for the inner harness: inside `flai-agent-claude-code:local` the worktree has no `bin/` and no `go`, so the hook's `scripts/flai.sh guard` runs `/usr/local/bin/flai`, built from the tree when the image was built; S-0510's layer brings `go` back, and the script then builds as on the host.

## Done when

From a copy of the scripts tree without `bin/`, with `PATH` holding a fake `flai` and no `go`, `scripts/flai.sh version` runs the fake and prints the warning on standard error; with neither on `PATH` it exits 127 with the message; with `go` present it builds and runs `bin/flai` as before; the `shellcheck` tier passes on it, and `scripts/README.md` describes the fallback.

## Notes

The fallback runs only without `go`, so a developer's host, where `go` is always present, never runs an older installed flai in place of the tree's.
