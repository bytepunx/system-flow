---
id: S-0222
type: story
nature: feature
title: "The orchestrator publishes by the release policy: judgement, theme, or cost of delay threshold"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-05T00:34:46Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, flai/internal/hostapi, flai/internal/release, design/system/strategic-agents.md, flai/internal/guard, flai/internal/manifest, design/system/project-manifest.md, docs/operators/settings.md]
after: [S-0218, S-0174]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:48:21Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T18:27:00Z
  basis: "Its own forecast of 1h15m; 19th in the pull order with an in-progress limit of 3, behind S-0249, S-0266, S-0268, S-0253, S-0257, S-0244, S-0262, S-0258, S-0260, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218, S-0219, S-0220 and S-0221."
  by: flai
  at: 2026-10-05T00:34:46Z
---
# S-0222 The orchestrator publishes by the release policy: judgement, theme, or cost of delay threshold

## Goal

With `publish` on, the orchestrator cuts releases when the project's release policy says so: a cost of delay or count threshold, a theme (an epic or tag whose stories are all accepted), or its own judgement when the policy is `judgement`.

## Acceptance criteria
- [ ] After each acceptance it runs `flai release --evaluate`; when the policy is met it publishes through `publish.run` (which requires the `push` host action too) and logs the release with the figures and the items bundled
- [ ] Under `judgement`, it publishes when it judges the unreleased work coherent and complete, and logs the reasoning; it never publishes a batch with a story whose epic is not in review or done when `orchestration.release.whole_epics` is set
- [ ] When publishing is refused (S-0174's tag check, a moved remote) it logs the refusal and opens a thread to the operator
- [ ] Tests cover each policy met and not met, and the refusal

## Tasks

## Notes
