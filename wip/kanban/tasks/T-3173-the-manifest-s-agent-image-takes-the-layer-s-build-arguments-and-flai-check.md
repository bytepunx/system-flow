---
id: T-3173
type: task
nature: feature
title: The manifest's agent.image takes the layer's build arguments, and flai check refuses a malformed one
status: in-progress
parent: S-0745
owner: alex
created: 2026-11-12T16:00:24Z
updated: 2026-11-13T13:38:01Z
transitions:
  - to: ready
    at: 2026-11-13T13:17:51Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T13:17:52Z
    by: system-flow
stream: S-0745
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/image.go, flai/internal/check/image_test.go, flai/internal/manifest/agent.go, flai/internal/manifest/agent_test.go, flai/internal/workitem/store_test.go, docs/operators/settings.md]
---
# T-3173 The manifest's agent.image takes the layer's build arguments, and flai check refuses a malformed one

## Work

Add the layer's build arguments to `agent.image` in `flai/internal/manifest/manifest.go`, beside `Layer` and `LayerImage` (lines 577 to 591), in the shape S-0733's ADR names. By this plan's assumption that is `args`, a map from an `ARG` name to a string. Give `Manifest.LayerArgs()`.

Make `layerLine` in `flai/internal/check/image.go` (lines 91 to 119) report an argument that is not a valid `ARG` name, or one the layer's Dockerfile does not declare, as an `agent.image` finding. Keep the scoping of ADR-0255 as it is: the finding stays in a story's check only when the story changes `system-flow.yaml` or the Dockerfile.

Test both, with a manifest that names no arguments read as before. It is the first layer.

## Done when

- `go test ./internal/manifest/ ./internal/check/` passes, with a good map, a bad name, an argument the Dockerfile lacks, and none.
- This repository's manifest still parses and checks clean.

## Notes
