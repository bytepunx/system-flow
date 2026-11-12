---
id: T-3164
type: task
nature: feature
title: flai/internal/agentimage lists the images a project's agents need and says whether each is present and current
status: backlog
parent: S-0735
owner: alex
created: 2026-11-12T15:58:34Z
updated: 2026-11-12T15:58:34Z
transitions: []
stream: S-0735
tags: [cli]
touches: [flai/internal/agentimage/agentimage.go, flai/internal/agentimage/status.go, flai/internal/agentimage/status_test.go]
---
# T-3164 flai/internal/agentimage lists the images a project's agents need and says whether each is present and current

## Work

Create the package `flai/internal/agentimage`.

`agentimage.go` holds the image decision moved from `containerImageHost` (`flai/internal/serve/container.go`, lines 76 to 102), unchanged in what it answers:

- the layer's `Manifest.LayerImage(inner)` (`flai/internal/manifest/manifest.go`, lines 585 to 591) when the manifest names one;
- else the host's `images.<inner>` (`config.HarnessHost.Image`);
- else `config.DefaultImage`.

Each image records where its name came from. `Needed(root, host)` lists every image the project's agents and roles need: the inner harnesses of its default agent, its stories' agents, and its roles' agents on `container`, and the layer over each.

`status.go` asks docker for each image with one `docker image inspect` call, through a runner the tests replace. It reports:

- present or missing;
- current or stale, by the label and rule of S-0733's ADR, with the reason;
- `own` for an image named by the host's `--image`, which is never stale.

Table tests cover each source, each state, and docker missing.

It is the first layer and waits for no task.

## Done when

- `go test ./internal/agentimage/` passes, with a case for each source and each state.
- The decision answers what `containerImageHost` answers today for every case in `TestContainerImageHostPutsTheLayerImageInACopyOfTheHostsImages`.

## Notes
