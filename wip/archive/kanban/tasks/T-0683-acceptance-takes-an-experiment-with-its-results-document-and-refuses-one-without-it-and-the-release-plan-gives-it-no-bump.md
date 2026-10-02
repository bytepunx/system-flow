---
id: T-0683
type: task
nature: feature
title: Acceptance takes an experiment with its results document and refuses one without it, and the release plan gives it no bump
status: done
parent: S-0194
owner: arobson
created: 2026-10-02T10:08:08Z
updated: 2026-10-02T10:11:10Z
transitions:
  - to: ready
    at: 2026-10-02T10:08:27Z
    by: agent-S-0194
  - to: in-progress
    at: 2026-10-02T10:08:27Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:11:10Z
    by: agent-S-0194
stream: S-0194
tags: []
usage:
  source: log
  seconds: 163
  estimated: true
  models:
    - model: claude-haiku-4-5-20251001
      input: 173
      output: 8396
      cache_read: 1145095
      cache_write: 71536
      cost: 0.2461
    - model: claude-opus-5-5
      input: 35
      output: 8330
      cache_read: 2410258
      cache_write: 27884
      cost: 0.8719
---

# T-0683 Acceptance takes an experiment with its results document and refuses one without it, and the release plan gives it no bump

## Work

`release.LevelFor` gives an `experiment` no release level, as for research, and the plan's skipped reason names the nature. `preview.Accept`, which `flai accept`, `flai move <story> done`, and the dashboard's acceptance all run before merging, blocks an experiment story that has no results document `design/experiments/<S-nnnn>-*.md` on its branch, and names the document it expects.

## Done when

- Tests cover an experiment accepted with its document, refused without it (naming the path), and planned with no bump and its components listed as unreleased.

## Notes
