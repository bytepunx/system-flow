---
id: T-3174
type: task
nature: feature
title: flai builds the project's layer image over the harness image with the manifest's arguments, labelled with the Dockerfile's sum and the arguments
status: backlog
parent: S-0745
owner: alex
created: 2026-11-12T16:00:29Z
updated: 2026-11-12T16:00:37Z
transitions: []
stream: S-0745
tags: [cli]
touches: [flai/internal/agentimage/layer.go, flai/internal/agentimage/layer_test.go, flai/internal/agentimage/status.go, flai/internal/agentimage/status_test.go]
after: [T-3173]
---
# T-3174 flai builds the project's layer image over the harness image with the manifest's arguments, labelled with the Dockerfile's sum and the arguments

## Work

`flai/internal/agentimage/layer.go` builds `Manifest.LayerImage(inner)` from the Dockerfile `agent.image.layer` names:

- with the project's root as the context, as `scripts/agent-image.sh`'s `build_layer` does (lines 95 to 110);
- with `BASE` set to the harness image and each of the manifest's arguments as `--build-arg`;
- with the label S-0733's ADR names, recording the Dockerfile's sha256 and the arguments.

It runs through S-0740's runner, so the step limit, the network retry, and the log file hold. Make S-0735's status (`status.go`) report the layer stale when the Dockerfile's sum, an argument, or the harness image's ID differs from the label's.

Tests over a fake docker cover:

- the arguments passed;
- a changed Dockerfile reported stale;
- a changed argument reported stale;
- a reset retried.

It waits for T-3173, whose `LayerArgs` it reads.

## Done when

- `go test ./internal/agentimage/ -run 'Layer|Status'` passes.
- The docker arguments a test records for this repository's manifest match `scripts/agent-image.sh --layer`'s, once the manifest names the arguments that script reads today.

## Notes
