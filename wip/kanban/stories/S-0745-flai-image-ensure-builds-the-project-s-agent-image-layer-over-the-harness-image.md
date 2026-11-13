---
id: S-0745
type: story
nature: feature
title: flai image ensure builds the project's agent.image.layer over the harness image itself, with build arguments the manifest gives in any project
status: ready
parent: E-0036
owner: alex
created: 2026-11-12T16:00:12Z
updated: 2026-11-13T07:32:52Z
transitions:
  - to: ready
    at: 2026-11-13T07:29:34Z
    by: alex
tags: [cli, template]
topics: [docker, manifest]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/agentimage/layer.go, flai/internal/agentimage/layer_test.go, flai/internal/agentimage/status.go, flai/internal/agentimage/status_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/cmd/image.go, flai/cmd/image_test.go, flai/internal/check/image.go, flai/internal/check/image_test.go, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, template/root/system-flow.yaml.tmpl]
after: [S-0740]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 12
          cache_read: 136891
          cache_write: 18181
          cost: 0.0443
forecast:
  duration: 40m
  delivery: 2026-11-13T08:52:00Z
  basis: "Its own forecast of 40m; 3rd in the pull order with an in-progress limit of 3, behind S-0660, S-0733, S-0760 and S-0735."
  by: flai
  at: 2026-11-13T07:32:52Z
finalized:
  by: alex
  at: 2026-11-13T07:27:59Z
---
# S-0745 flai image ensure builds the project's agent.image.layer over the harness image itself, with build arguments the manifest gives in any project

## Goal

A project's layer image, the Dockerfile `agent.image.layer` names, is built today only by `scripts/agent-image.sh`. That script reads its build arguments from this repository's own files: `GO_VERSION` from `flai/go.mod`, and the tool pins from `scripts/install-tools.sh`. No other project has those files, so no other project can build its layer through flai.

After this story `flai image ensure` builds the layer over the harness image it ensured, when the layer image is missing or stale. It takes the layer's build arguments in the general way S-0733's ADR decides: by this plan's assumption, a map under `agent.image` in the manifest, with the Dockerfile's `ARG` defaults for the rest. It is held to the same step limit, retry, and log as the harness images.

## Acceptance criteria

- [ ] The manifest's `agent.image` takes the layer's build arguments as S-0733's ADR names them, and `flai check` refuses a malformed one. A manifest that names a layer and no arguments builds with the Dockerfile's own `ARG` defaults
- [ ] `flai image ensure` builds the layer image (`Manifest.LayerImage`) over the harness image, once that image is present and current, with the manifest's arguments and `BASE` set to the harness image. It runs with the step limit, the network retry, and the log of S-0740's build
- [ ] The layer image carries the label S-0733's ADR names, recording the Dockerfile's sum and the arguments, so `flai image status` reports it stale when the Dockerfile, an argument, or its base changes
- [ ] `flai image build --layer` builds the layer alone over the harness image on the host, as `scripts/agent-image.sh --layer` does today
- [ ] `design/system/project-manifest.md`, the manifest rows of `docs/operators/settings.md`, `docs/users/flai.md`, `docs/operators/runbooks/agent-container.md`, and the commented `agent.image` block of `template/root/system-flow.yaml.tmpl` say how a project names its layer and its arguments, and that flai builds it

## Tasks
- T-3173 The manifest's agent.image takes the layer's build arguments, and flai check refuses a malformed one
- T-3174 flai builds the project's layer image over the harness image with the manifest's arguments, labelled with the Dockerfile's sum and the arguments
- T-3175 flai image ensure builds a missing or stale layer image after the harness images, and flai image build --layer builds it alone
- T-3176 The manifest design, the settings index, the guides, and the template's manifest say how a project names its layer and its arguments, and that flai builds it

## Notes

Proposed slice 3c of E-0036. This repository's own layer moves onto these arguments in S-0760, so `scripts/agent-image.sh` is left as it is here.

### Planning

- Touches:
  - `flai/internal/manifest/manifest.go`, `manifest_test.go`: layout (`Layer` and `LayerImage`, lines 577 to 591).
  - `flai/internal/agentimage/layer.go`, `layer_test.go`: new files of the layer's build.
  - `flai/internal/agentimage/status.go`, `status_test.go`: S-0735's status learns the layer's staleness (third criterion).
  - `flai/internal/agentimage/ensure.go`, `ensure_test.go`, `flai/cmd/image.go`, `image_test.go`: S-0740's `ensure` and `build` gain the layer (second and fourth criteria).
  - `flai/internal/check/image.go`, `image_test.go`: layout (`layerLine`, lines 91 to 119) for a malformed argument.
  - `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/users/flai.md`, `docs/operators/runbooks/agent-container.md`, `template/root/system-flow.yaml.tmpl`: named in the last criterion. `TestSettingsIndex` fails on a manifest key with no row.
  - `design/system/flai-cli.md` (co-change, 61%) describes the layer image in § The project's layer image, and `docs/users/flai-reference.md` (31%) gains `--layer`.
- `after`:
  - S-0733 is not named, though this story reads the ADR's shape for the layer's arguments and its label: it waits for S-0733 through S-0740 and S-0735, which names it.
  - S-0740: the layer is built through that story's runner (step limit, retry, log) and its `ensure` and `build`.
- Tasks, in four layers, one after another: T-3173 (the key) → T-3174 (the build, which reads `LayerArgs`) → T-3175 (`ensure` and `--layer`, which call it) → T-3176 (the documents).
- Forecast 40m: `flai forecast` gives 33m on the 13 touches the story declared when it ran (109 s per unit over 118 done large feature stories, size 18). Raised for the 19 its tasks name. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
