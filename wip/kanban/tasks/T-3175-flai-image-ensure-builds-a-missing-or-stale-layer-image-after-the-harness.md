---
id: T-3175
type: task
nature: feature
title: flai image ensure builds a missing or stale layer image after the harness images, and flai image build --layer builds it alone
status: backlog
parent: S-0745
owner: alex
created: 2026-11-12T16:00:41Z
updated: 2026-11-12T16:00:41Z
transitions: []
stream: S-0745
tags: [cli]
touches: [flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/cmd/image.go, flai/cmd/image_test.go]
after: [T-3174]
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
