---
id: S-0745
type: story
nature: feature
title: flai image ensure builds the project's agent.image.layer over the harness image itself, with build arguments the manifest gives in any project
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T16:00:12Z
updated: 2026-11-12T16:00:56Z
transitions: []
tags: [cli, template]
topics: [docker, manifest]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/agentimage/layer.go, flai/internal/agentimage/layer_test.go, flai/internal/agentimage/ensure.go, flai/internal/agentimage/ensure_test.go, flai/internal/check/image.go, flai/internal/check/image_test.go, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, docs/operators/runbooks/agent-container.md, template/root/system-flow.yaml.tmpl]
after: [S-0733, S-0740]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
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

Proposed slice 3c of E-0036. This repository's own layer moves onto these arguments in the plan's last story, so `scripts/agent-image.sh` is left as it is here.
