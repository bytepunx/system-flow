---
id: T-3166
type: task
nature: improvement
title: The container start and flai check's agent.image take the image and its state from flai/internal/agentimage
status: backlog
parent: S-0735
owner: alex
created: 2026-11-12T15:58:46Z
updated: 2026-11-12T16:01:26Z
transitions: []
stream: S-0735
tags: [cli]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/check/image.go, flai/internal/check/image_test.go]
after: [T-3164]
---
# T-3166 The container start and flai check's agent.image take the image and its state from flai/internal/agentimage

## Work

Make `containerImageHost` and `containerImage` in `flai/internal/serve/container.go` (lines 76 to 143) call T-3164's package for the image and its presence. The close-out container (`verifyjob_container.go`, line 115) and the checks container (`checks.go`, line 250) call `containerImage`, so they follow unchanged. The refusal keeps its words and adds `flai image status`.

Make `AgentImage` in `flai/internal/check/image.go` (lines 43 to 89) take presence from the package. Keep its findings and their scoping in `scope.go` (ADR-0255) as they are, and name `flai image status` in the advice. Neither yet reports a stale image as a refusal or a finding: S-0751, in which `flai serve` brings images back, decides that.

It waits for T-3164. It touches no file of T-3165, so the two run together.

## Done when

- `go test ./internal/serve/ ./internal/check/` passes, with `TestAContainerStartIsRefusedWhatItLacks`, `TestTheProjectsLayerImageRunsItsContainerAgents`, `TestContainerImageHostPutsTheLayerImageInACopyOfTheHostsImages`, and `TestAgentImage` unchanged in what they assert but the added command name.
- No code outside `flai/internal/agentimage` decides an agent image's name.

## Notes
