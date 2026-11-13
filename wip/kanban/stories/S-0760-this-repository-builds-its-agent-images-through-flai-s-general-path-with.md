---
id: S-0760
type: story
nature: improvement
title: This repository builds its agent images through flai's general path, with scripts/agent-image.sh as its own build command, and the documents and the template say how any project gets its images
status: ready
parent: E-0036
owner: alex
created: 2026-11-12T16:06:26Z
updated: 2026-11-13T15:22:36Z
transitions:
  - to: ready
    at: 2026-11-13T07:31:21Z
    by: alex
tags: [cli, template]
topics: [docker, template]
touches: [scripts/agent-image.sh, scripts/agent-image-test.sh, Makefile, system-flow.yaml, flai/images/system-flow/Dockerfile, docs/operators/runbooks/agent-container.md, docs/operators/runbooks/install.md, docs/operators/index.md, docs/users/flai.md, design/system/agent-adapters.md, design/tech/docker.md, template/root/docs/operators/index.md.tmpl, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md]
after: [S-0745, S-0759]
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
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 20
          cache_read: 112634
          cache_write: 18925
          cost: 0.0377
forecast:
  duration: 32m
  delivery: 2026-11-13T16:09:00Z
  basis: "Its own forecast of 32m; 1st in the pull order with an in-progress limit of 3, behind S-0663, S-0751 and S-0773."
  by: flai
  at: 2026-11-13T15:22:36Z
finalized:
  by: alex
  at: 2026-11-13T07:31:18Z
---
# S-0760 This repository builds its agent images through flai's general path, with scripts/agent-image.sh as its own build command, and the documents and the template say how any project gets its images

## Goal

This repository builds its agent images outside flai today: `make agent-images` runs `scripts/agent-image.sh`, which builds `flai` and `pilot` from the tree. It reads the layer's build arguments from `flai/go.mod` and `scripts/install-tools.sh`.

After this story the repository runs on the path every project uses:

- `system-flow.yaml` names its layer's build arguments the general way of S-0745;
- `scripts/agent-image.sh` is the build command the host names for this project (S-0759), building from the tree as the development tree S-0733's ADR describes;
- `make agent-images` stays as its entry.

The documents say what is true, and the template's documents tell a new project how it gets its images: from flai, with nothing of its own to build unless it names a layer.

## Acceptance criteria

- [ ] `system-flow.yaml`'s `agent.image` names the layer's build arguments in the shape S-0745 gives, written with the flai command that sets it. `flai image build --layer` builds `flai-agent-sf-claude-code:local` with the same arguments `scripts/agent-image.sh --layer` passes today
- [ ] `scripts/agent-image.sh` works as the host's build command for this project: it reads the variables S-0759 passes, builds the one image asked for from the tree, and keeps `--tag`, `--layer`, and `--check`. `scripts/agent-image-test.sh` covers the variables, and `make agent-images` still builds all three
- [ ] `docs/operators/runbooks/agent-container.md` and `docs/operators/runbooks/install.md` § Agents in containers give `flai image ensure`, and the automatic job, as the way to get images in any project. `make agent-images` and the build command are this repository's own
- [ ] `docs/operators/index.md` § The image, `design/system/agent-adapters.md` § A container harness of flai's own, and `design/tech/docker.md` § Agent images say what is built after E-0036: where the images come from, the label, the job, and the substitution. None of them says that only this repository can build the images
- [ ] The template's operators guide (`template/root/docs/operators/index.md.tmpl`) and the commented `agent:` block of `template/root/system-flow.yaml.tmpl` say that a new project on the `container` harness gets its images from flai, and how it names a layer and its arguments. `template/CHANGELOG.md` records it under `## Unreleased`

## Tasks
- T-3207 system-flow.yaml names this repository's layer build arguments the general way, and flai image build --layer builds the layer as scripts/agent-image.sh does
- T-3209 scripts/agent-image.sh serves as this project's build command, building from the tree the one image flai asks for, and make agent-images still builds all three
- T-3212 The template's operators guide and manifest tell a new project that flai supplies its agent images, and how it names a layer and its arguments
- T-3214 The operator runbooks, the operators guide, agent-adapters.md, and docker.md say how any project gets its agent images after E-0036, and what is this repository's own

## Notes

Proposed slice 6 of E-0036, its last. It is an improvement: it changes how this repository builds, and adds nothing a user of flai lacks after S-0759.

Naming `scripts/agent-image.sh` as the host's build command is the operator's to do, once, on the host: `flai guard` refuses it to agents (S-0759). T-3209 records the command for them.

### Planning

- Touches:
  - `system-flow.yaml`, `flai/images/system-flow/Dockerfile`: the first criterion. The manifest's key is written with flai's command, never by hand.
  - `scripts/agent-image.sh`, `scripts/agent-image-test.sh`, `Makefile`: the second criterion (`agent-images`, `Makefile` lines 69 to 70).
  - `docs/operators/runbooks/agent-container.md`, `docs/operators/runbooks/install.md`, `docs/operators/index.md`, `design/system/agent-adapters.md`, `design/tech/docker.md`: named in the third and fourth criteria.
  - `docs/users/flai.md` (co-change, 49%) still names `make agent-images` at line 501.
  - `template/root/docs/operators/index.md.tmpl`, `template/root/system-flow.yaml.tmpl`, `template/CHANGELOG.md`: the last criterion.
  - Left out from `flai touches suggest`: `design/system/flai-cli.md` (47%), which S-0733 to S-0759 keep current as they build, and the dashboard's documents.
- `after`:
  - S-0745: this story writes the manifest's layer arguments in the shape that story gives, and builds with its `flai image build --layer`.
  - S-0759: `scripts/agent-image.sh` becomes the build command that story's setting names, reading its variables.
- Tasks, in two layers: T-3207 (the manifest), T-3209 (the script), and T-3212 (the template) share no file and run together; T-3214 describes T-3207's and T-3209's result.
- Forecast 32m stands as `flai forecast` gives it (107 s per unit over 175 done large improvement stories, size 18). Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
