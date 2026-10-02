---
id: T-0686
type: task
nature: feature
title: The design and the user and operator guides describe accepting an experiment
status: done
parent: S-0194
owner: arobson
created: 2026-10-02T10:08:08Z
updated: 2026-10-02T10:16:21Z
transitions:
  - to: ready
    at: 2026-10-02T10:08:27Z
    by: agent-S-0194
  - to: in-progress
    at: 2026-10-02T10:14:58Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:16:21Z
    by: agent-S-0194
stream: S-0194
tags: []
usage:
  source: log
  seconds: 83
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 51
      output: 12347
      cache_read: 3572601
      cache_write: 41332
      cost: 1.2923
---

# T-0686 The design and the user and operator guides describe accepting an experiment

## Work

`design/system/workflow.md`, `work-hierarchy.md`, `flai-cli.md`, `flaiover-dashboard.md`, `docs/users/`, and `docs/operators/` say that an experiment is accepted with its results document and released with nothing.

## Done when

- No document under `design/system` or `docs` says an experiment stays on its branch or is refused at acceptance.

## Notes
