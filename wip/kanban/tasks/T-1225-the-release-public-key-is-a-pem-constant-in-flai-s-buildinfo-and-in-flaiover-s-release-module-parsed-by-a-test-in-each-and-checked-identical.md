---
id: T-1225
type: task
nature: feature
title: The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical
status: done
parent: S-0232
owner: alex
created: 2026-10-07T22:15:41Z
updated: 2026-10-12T15:06:27Z
transitions:
  - to: ready
    at: 2026-10-12T15:05:05Z
    by: agent-S-0232
  - to: in-progress
    at: 2026-10-12T15:05:06Z
    by: agent-S-0232
  - to: done
    at: 2026-10-12T15:06:27Z
    by: agent-S-0232
stream: S-0232
tags: []
touches: [flai/internal/buildinfo/releasekey.go, flai/internal/buildinfo/releasekey_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts, flai/internal/buildinfo/cosign.pub, flai/.goreleaser.yaml, docs/operators/runbooks/release-key.md]
after: [T-1222]
usage:
  source: log
  seconds: 80
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 137
      cache_read: 3659181
      cache_write: 8711
      cost: 1.6189
---
# T-1225 The release public key is a PEM constant in flai's buildinfo and in flaiover's release module, parsed by a test in each and checked identical

## Work

Put the release public key, as PEM, in a constant in `flai/internal/buildinfo/releasekey.go` and in `flaiover/src/lib/server/release.ts`, each with a test that parses it as an ECDSA P-256 public key, and a flaiover test that reads the Go source from the checkout and finds the same PEM.

Waits for T-1222: the key is the one the operator generates with its runbook.

## Done when

- Both constants hold the operator's public key, each parsed by a test, and a test fails when the two differ.

## Notes
