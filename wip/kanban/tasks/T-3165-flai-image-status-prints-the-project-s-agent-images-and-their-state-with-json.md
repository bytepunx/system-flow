---
id: T-3165
type: task
nature: feature
title: flai image status prints the project's agent images and their state, with --json, and exits non-zero on one missing or stale
status: backlog
parent: S-0735
owner: alex
created: 2026-11-12T15:58:41Z
updated: 2026-11-12T15:58:41Z
transitions: []
stream: S-0735
tags: [cli]
touches: [flai/cmd/image.go, flai/cmd/image_test.go, flai/cmd/root.go]
after: [T-3164]
---
# T-3165 flai image status prints the project's agent images and their state, with --json, and exits non-zero on one missing or stale

## Work

Add `flai image` with the subcommand `status` in `flai/cmd/image.go`, and register it in `flai/cmd/root.go` beside `newAgentCmd` (line 152). It opens the project and the host's configuration, as `flai serve agent show` does, calls T-3164's `Needed` and the status, and prints:

- a table: image, for which inner harness, source, present, current, and reason;
- with `--json`, the same as an array.

It exits 1 when an image is missing or stale and 0 otherwise, with docker missing named on standard error. Tests in `flai/cmd/image_test.go` run it over a fake docker. It waits for T-3164, whose package it calls.

## Done when

- `go test ./cmd/ -run Image` passes, covering present, missing, stale, own, and docker missing.
- `flai image status --help` describes the command and its exits.

## Notes
