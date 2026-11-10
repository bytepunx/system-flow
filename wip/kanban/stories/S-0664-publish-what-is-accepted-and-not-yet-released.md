---
id: S-0664
type: story
nature: improvement
title: Publish what is accepted and not yet released
status: review
owner: alex
created: 2026-11-10T08:33:20Z
updated: 2026-11-10T08:35:50Z
transitions:
  - to: ready
    at: 2026-11-10T08:33:22Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:33:28Z
    by: agent-S-0664
  - to: review
    at: 2026-11-10T08:35:50Z
    by: agent-S-0664
tags: []
touches: [wip/kanban/stories/S-0664-publish-what-is-accepted-and-not-yet-released.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 152
  priced_by: harness
  turns:
    - day: 2026-11-10
      ceremony: 2
      hand_edits: 1
      work: 28
  models:
    - model: claude-opus-5-5
      input: 68
      output: 9199
      cache_read: 4775623
      cache_write: 159918
      cost: 2.4187
  strategic:
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 9
          output: 38
          cache_read: 298266
          cache_write: 55378
          cost: 0.0957
formula: publish
---
# S-0664 Publish what is accepted and not yet released

## Goal

Publish what is accepted and not yet released, by the project's release policy.

## Acceptance criteria

- [x] Fetch the remote's branch and tags
- [x] Evaluate the release policy, stopping unless it is met
- [x] Publish what is accepted and not yet released

## Tasks
- T-2890 Fetch the remote's branch and tags
- T-2891 Evaluate the release policy, stopping unless it is met
- T-2892 Publish what is accepted and not yet released

## Notes

Made by `flai formula run publish` from `design/workflows/formulas/publish.md`, with the inputs `judgement=`.

Published: flai 1.86.11 to 1.86.12 (patch), tag `flai/v1.86.12`, bundling S-0652 and S-0637; branch and tag pushed. No other component had anything pending. The steps' `flai` commands ran as `scripts/flai.sh` in the main checkout, the tree's flai, since the installed flai (1.86.6) lacks `--judged`.
