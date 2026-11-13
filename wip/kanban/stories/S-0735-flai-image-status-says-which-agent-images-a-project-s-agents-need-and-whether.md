---
id: S-0735
type: story
nature: feature
title: flai image status says which agent images a project's agents need, and whether each is on the host and current for the flai that runs it
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T15:58:14Z
updated: 2026-11-13T07:24:20Z
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
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 8
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 8
          cache_read: 61935
          cache_write: 11102
          cost: 0.0208
draft: true
forecast:
  duration: 33m
  delivery: 2026-11-13T13:56:00Z
  basis: "Its own forecast of 33m; 9th in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663, S-0730, S-0733 and S-0734."
  by: flai
  at: 2026-11-13T07:24:20Z
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

### Planning

- Touches:
  - `flai/internal/agentimage/agentimage.go`, `status.go`, `status_test.go`: new files of the package the first two criteria describe.
  - `flai/cmd/image.go`, `image_test.go`: the new command. `flai/cmd/root.go` registers commands (line 152).
  - `flai/internal/serve/container.go`, `container_test.go`: layout. `containerImageHost` (lines 76 to 102) and `containerImage` (lines 120 to 143) are where the image is decided and refused today. The close-out (`verifyjob_container.go`, line 115) and the checks (`checks.go`, line 250) call `containerImage` unchanged, so they are not touched.
  - `flai/internal/check/image.go`, `image_test.go`: layout (`AgentImage`, lines 43 to 89).
  - `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md`: named in the last criterion.
  - Left out from `flai touches suggest`: `docs/operators/index.md` (22%), `docs/operators/settings.md` (21%), and the dashboard's documents. This story adds no setting and no page.
- `after` S-0733: this story reads the ADR's label, its staleness rule, and the command's name.
- Tasks, in three layers: T-3164 writes the package; T-3165 (the command) and T-3166 (the start and the check) both call it and share no file, so they run together; T-3167 describes both.
- Forecast 33m stands as `flai forecast` gives it (109 s per unit over 118 done large feature stories, size 18). Most of the work is a move of code that exists, with its tests kept. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
