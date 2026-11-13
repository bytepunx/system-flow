---
id: T-3175
type: task
nature: feature
title: flai image ensure builds a missing or stale layer image after the harness images, and flai image build --layer builds it alone
status: done
parent: S-0745
owner: alex
created: 2026-11-12T16:00:41Z
updated: 2026-11-13T14:12:27Z
transitions:
  - to: ready
    at: 2026-11-13T13:56:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T13:56:51Z
    by: system-flow
  - to: done
    at: 2026-11-13T14:12:27Z
    by: agent-S-0745
stream: S-0745
tags: [cli]
touches: [flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/cmd/image.go, flai/cmd/image_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-3174]
usage:
  source: log
  seconds: 936
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 47
      output: 17463
      cache_read: 3581551
      cache_write: 137189
      cost: 1.8842
---
# T-3175 flai image ensure builds a missing or stale layer image after the harness images, and flai image build --layer builds it alone

## Work

Extend `Ensure` in `flai/internal/agentimage/ensure.go`: once the harness images are present and current, build each missing or stale layer image with T-3174's build. Skip a layer whose harness image failed, naming why. Add `--layer` to `flai image build` in `flai/cmd/image.go`: it builds the layer alone over the harness image already on the host, and refuses with the harness image's name when that is missing.

Tests cover:

- the layer built after a pull;
- a layer skipped after a failed harness;
- `--layer` alone;
- `--layer` with no harness image.

It waits for T-3174, whose build it calls.

## Done when

- `go test ./internal/agentimage/ ./cmd/ -run 'Ensure|Image'` passes.
- After an `ensure` over the fake docker, `flai image status` reports the layer current.

## Notes
