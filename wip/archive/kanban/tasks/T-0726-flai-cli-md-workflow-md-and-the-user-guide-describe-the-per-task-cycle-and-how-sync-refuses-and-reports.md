---
id: T-0726
type: task
nature: improvement
title: flai-cli.md, workflow.md, and the user guide describe the per-task cycle and how sync refuses and reports
status: done
parent: S-0197
owner: arobson
created: 2026-10-03T00:40:15Z
updated: 2026-10-03T00:59:21Z
transitions:
  - to: ready
    at: 2026-10-03T00:40:34Z
    by: agent-S-0197
  - to: in-progress
    at: 2026-10-03T00:52:36Z
    by: agent-S-0197
  - to: done
    at: 2026-10-03T00:59:21Z
    by: agent-S-0197
stream: S-0197
tags: []
touches: [design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0723]
usage:
  source: log
  seconds: 405
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 97
      output: 37539
      cache_read: 3287781
      cache_write: 130574
      cost: 2.2133
---
# T-0726 flai-cli.md, workflow.md, and the user guide describe the per-task cycle and how sync refuses and reports

## Work

Describe the per-task cycle (order settled on TH-0072) and the sync before review in `design/system/workflow.md` and `design/system/flai-cli.md`, and in the user guide (`docs/users/flai.md`, and `docs/users/flai-reference.md` for `stream sync`). Say that `flai stream sync` refuses uncommitted changes and how it reports a conflict. Waits for T-0723, whose behaviour it describes.

## Done when

- [ ] `flai-cli.md`, `workflow.md`, and the user guide describe the per-task cycle and the sync before review
- [ ] They describe sync's refusal and conflict report as T-0723 built them

## Notes
