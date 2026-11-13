---
id: S-0733
type: story
nature: feature
title: An ADR decides where an installed flai's agent images come from, how each is tied to the flai that runs it, and how flai brings back a missing one or hands the build to the operator's own machinery
status: backlog
parent: E-0036
owner: alex
created: 2026-11-12T15:54:40Z
updated: 2026-11-13T07:28:33Z
transitions: []
tags: [cli]
topics: [release, docker, security]
touches: [design/adrs/README.md, design/system/flai-cli.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
forecast:
  duration: 25m
  delivery: 2026-11-13T13:17:00Z
  basis: "Its own forecast of 25m; 7th in the pull order with an in-progress limit of 3, behind S-0659, S-0241, S-0660, S-0661, S-0662, S-0663 and S-0730."
  by: flai
  at: 2026-11-13T07:23:41Z
finalized:
  by: alex
  at: 2026-11-13T07:28:33Z
---
# S-0733 An ADR decides where an installed flai's agent images come from, how each is tied to the flai that runs it, and how flai brings back a missing one or hands the build to the operator's own machinery

## Goal

Today only this repository can build the agent images. `scripts/agent-image.sh` builds `flai` and `pilot` from the tree, a released flai has no `pilot`, and nothing publishes the images ([ADR-0212](../../../design/adrs/0212-the-agent-images-are-layered-a-bedrock-layer-and-one-layer-per-inner-harness.md), decision 8). [ADR-0203](../../../design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md) intended a signed image, published and run by digest, which E-0029 left "to a later story".

An ADR refining ADR-0212 settles how the stories after it work, so that they build one design. It decides:

- where an installed flai's harness images come from;
- how an image is tied to the flai that runs it;
- what `flai serve` does when an image is missing or stale;
- how an operator customizes the images or replaces flai's build.

`design/system/flai-cli.md` then describes that design as decided.

## Acceptance criteria

- [ ] An accepted ADR, refining ADR-0212, decides where an installed flai's bedrock and claude-code images come from: pulled by a digest a signed list names, as `flai dashboard` pulls flaiover; built on the host from Dockerfiles flai carries, with `flai` and `pilot` for Linux from the release; or pulled first with the build as the fallback. It names the registry, the image names, and the signed list
- [ ] It says how an image is tied to the flai that runs it, by the `dev.system-flow.flai` label or another, and what makes an image stale: another flai version, another Dockerfile, or another base
- [ ] It says what a development tree such as this repository runs, where the images must hold the tree's flai rather than a release's
- [ ] It says how the project's layer (`agent.image.layer`) gets its build arguments in any project, in place of `scripts/agent-image.sh` reading `flai/go.mod` and `scripts/install-tools.sh`
- [ ] It says what `flai serve` does when a start finds an image missing or stale: a host job, one at a time, that holds the start, and shows its output in the dashboard. It says whether the first automatic build needs the operator's consent. It sets a limit on a build step with no progress, and a retry of a step that fails on the network
- [ ] It says whether a strategic agent, which never gets the project's layer, waits for the layer's build or starts once the harness image is there
- [ ] It says where an operator names a build command of their own, which flai runs in place of its own (host configuration or a protected path, never the manifest alone), and how an operator turns the automatic build off
- [ ] `design/system/flai-cli.md` describes the decided design under § Internal structure, links the ADR, and marks each part not yet built with the story that builds it

## Tasks
- T-3149 Write the ADR, proposed, refining ADR-0212 on where an installed flai's agent images come from and how flai keeps them
- T-3151 Describe in design/system/flai-cli.md where the agent images come from as decided, marking each part not yet built
- T-3152 Ask the operator to confirm the ADR's choices on a thread, and accept the ADR with their answer

## Notes

Proposed slice 1 of E-0036. The epic's § Open questions for planning are this ADR's to settle with the operator.

The sixth criterion's premise needs checking: the strategic roles go through `startOf` and `containerImageHost` too (`flai/internal/serve/agents_start.go`, lines 398 to 424), so today the project's layer applies to them as well. The ADR says whether it should.

### Planning

- Touches:
  - `design/adrs/README.md`: layout. `flai adr new` and `flai adr accept` write the index again. The new ADR's file takes the number `flai adr new` gives it when the story runs, and `flai adr new --commit` adds it to the story's touches, so there is no folder touch on `design/adrs`.
  - `design/system/flai-cli.md`: design. § The container harness (S-0508) is where the image is decided today (`containerImageHost`), and the last criterion has the decided design described there.
  - ADR-0212 is refined, not superseded, so its file is not touched.
  - Left out from `flai touches suggest`: `docs/users/flai.md` (53%), `docs/operators/index.md`, and the other hub documents. They describe what is built, and the stories after this one change them as they build it.
- Topics `docker`, `release`, and `security` select `design/tech/docker.md`, `design/system/release-signing.md`, and the signing ADRs, which the decision rests on.
- No `after`: the story reads no other story's output.
- Tasks, in three layers: T-3149 writes the ADR; T-3151 describes it in `flai-cli.md`, since it links the ADR by its number; T-3152 asks the operator and accepts, since it confirms both together.
- Forecast 25m: `flai forecast` gives 20m (116 s per unit over 5 done medium feature stories, size 10). Raised 5m: the ADR answers seven questions, read against ADR-0203, ADR-0212, and the signing design, against six for S-0659's comparable ADR. The wait on the operator in T-3152 is wall time, not agent time. Delivery is as flai gives it.
- Cost of delay: not set yet. E-0036 has no inputs; TH-1094 asks the operator for them.
