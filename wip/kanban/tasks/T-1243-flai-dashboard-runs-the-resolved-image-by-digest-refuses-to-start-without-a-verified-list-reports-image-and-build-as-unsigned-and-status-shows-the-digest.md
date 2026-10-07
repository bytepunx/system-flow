---
id: T-1243
type: task
nature: feature
title: flai dashboard runs the resolved image by digest, refuses to start without a verified list, reports --image and --build as unsigned, and status shows the digest
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:54:07Z
updated: 2026-10-07T22:54:07Z
transitions: []
stream: S-0236
tags: [cli]
touches: [flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, flai/cmd/dashboard_watch.go, flai/cmd/dashboard_watch_test.go]
after: [T-1241]
---
# T-1243 flai dashboard runs the resolved image by digest, refuses to start without a verified list, reports --image and --build as unsigned, and status shows the digest

## Work

In `flai/cmd/dashboard.go`, replace `dashboardSettings.ref()`'s `image:tag` for the configured image with T-1241's resolver: `flai dashboard` pulls and runs `ghcr.io/bytepunx/flaiover@sha256:…`, and when the resolver refuses, it starts no container and prints the reason. `--image` naming another image and `--build` run what they name, as today, and the command says the container is unsigned.

`flai dashboard status`, with `--json`, shows the running container's image digest, the release version whose verified list names it, and whether it is on a signed list, or `unsigned` for `--image` and `--build`.

The record `dashboard_watch.go` writes keeps the by-digest reference, so that `flai dashboard restart` and `flai host`'s watch restart the same release by digest without a fresh resolution (ADR-0118). Change the record only if it needs the version beside the reference; a record from an older flai, with an `image:tag` reference, still restarts.

Cover with `fakeRunner` and a GitHub stand-in: a start by digest, a refusal for a missing list and a bad one, `--image` and `--build` reported unsigned, status's new fields, and a restart that keeps the digest.

Waits for T-1241, whose resolver it calls.

## Done when

- `flai dashboard` pulls and runs the image by digest from a verified list, and refuses to start, saying why, when the list is missing or does not verify.
- `--image` naming another image and `--build` still run and are reported unsigned.
- `flai dashboard status` shows the digest, its release, and whether it is on a signed list.
- A restart keeps the by-digest reference.
- The cases are tested and `flai test flai/cmd` passes.

## Notes
