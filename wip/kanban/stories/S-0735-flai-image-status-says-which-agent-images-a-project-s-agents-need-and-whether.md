---
id: S-0735
type: story
nature: feature
title: flai image status says which agent images a project's agents need, and whether each is on the host and current for the flai that runs it
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T15:58:14Z
updated: 2026-11-12T15:58:56Z
transitions: []
tags: [cli]
topics: [docker]
touches: [flai/internal/agentimage/agentimage.go, flai/internal/agentimage/status.go, flai/internal/agentimage/status_test.go, flai/cmd/image.go, flai/cmd/image_test.go, flai/cmd/root.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/check/image.go, flai/internal/check/image_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [S-0733]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0735 flai image status says which agent images a project's agents need, and whether each is on the host and current for the flai that runs it

## Goal

The image a project's agent runs in is decided in `containerImageHost` (`flai/internal/serve/container.go`). Its absence is found in two places that say it differently: `containerImage`, which refuses the start, and `flai check`'s `agent.image` (`flai/internal/check/image.go`). Nothing says whether an image on the host is current for the flai that runs it.

After this story one package, `flai/internal/agentimage`, decides which images a project's agents need, and whether each is present and current by the label and the staleness rule of S-0733's ADR. `flai image status` reports it. `flai serve` and `flai check` read the same answer.

## Acceptance criteria

- [ ] `flai/internal/agentimage` gives, for a project and the host's configuration, each image its agents need. That is the harness image of each inner harness the project's agents and roles use, and the layer image when the manifest names `agent.image.layer`. For each it says where the name came from: the manifest's layer, the host's `--image`, or the default
- [ ] For each image it says present or missing, and current or stale, by the label and the rule S-0733's ADR sets, with the reason: another flai version, another Dockerfile, or another base. An image of the operator's own, named with `--image`, is reported as theirs and never stale
- [ ] `flai image status`, with `--json`, prints that list for the project and exits non-zero when an image is missing or stale
- [ ] `containerImageHost` and `containerImage` in `flai/internal/serve/container.go`, and `AgentImage` in `flai/internal/check/image.go`, take the image and its state from the package, so the start, the close-out container, the checks container, and `flai check` cannot disagree. Their tests still pass, and the refusal and the finding name `flai image status`
- [ ] `docs/users/flai.md` § In a container, the generated `docs/users/flai-reference.md`, and `design/system/flai-cli.md` § The container harness (S-0508) describe the command and the one place the image is decided

## Tasks
- T-3164 flai/internal/agentimage lists the images a project's agents need and says whether each is present and current
- T-3165 flai image status prints the project's agent images and their state, with --json, and exits non-zero on one missing or stale
- T-3166 The container start and flai check's agent.image take the image and its state from flai/internal/agentimage
- T-3167 The user guide, the command reference, and flai-cli.md describe flai image status and the one place an agent image is decided

## Notes

Proposed slice 3a of E-0036. The command's name is this plan's assumption, `flai image`, since `flai agent` already sets the project's default agent (`flai/cmd/agent.go`). S-0733's ADR may name it otherwise, and the story's agent follows the ADR.

The strategic roles go through `startOf` and `containerImageHost` too (`flai/internal/serve/agents_start.go`, lines 398 to 424), so today the project's layer applies to them as well. Whether it should is S-0733's ADR's point 6; this story keeps what the ADR decides.
