---
id: T-1225
type: task
nature: feature
title: The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical
status: backlog
parent: S-0232
owner: alex
created: 2026-10-07T22:15:41Z
updated: 2026-10-07T22:15:47Z
transitions: []
stream: S-0232
tags: []
touches: [flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts]
after: [T-1222]
---
# T-1225 The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical

## Work

Put the release public key, as PEM, in a constant in `flai/internal/buildinfo/releasekey.go` and in `flaiover/src/lib/server/release.ts`, each with a test that parses it as an ECDSA P-256 public key, and a flaiover test that reads the Go source from the checkout and finds the same PEM.

Waits for T-1222: the key is the one the operator generates with its runbook.

## Done when

- Both constants hold the operator's public key, each parsed by a test, and a test fails when the two differ.

## Notes
