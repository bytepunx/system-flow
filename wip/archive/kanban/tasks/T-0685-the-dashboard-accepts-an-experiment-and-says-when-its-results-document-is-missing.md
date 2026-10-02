---
id: T-0685
type: task
nature: feature
title: The dashboard accepts an experiment and says when its results document is missing
status: done
parent: S-0194
owner: arobson
created: 2026-10-02T10:08:08Z
updated: 2026-10-02T10:14:58Z
transitions:
  - to: ready
    at: 2026-10-02T10:08:27Z
    by: agent-S-0194
  - to: in-progress
    at: 2026-10-02T10:14:02Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:14:58Z
    by: agent-S-0194
stream: S-0194
tags: []
usage:
  source: log
  seconds: 56
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 6712
      cache_read: 1942076
      cache_write: 22468
      cost: 0.7025
---

# T-0685 The dashboard accepts an experiment and says when its results document is missing

## Work

The dashboard's acceptance no longer disables an experiment story by nature; it shows the blocker flai reports when the results document is missing, and the release plan's no-release text for an experiment.

## Done when

- Component tests cover an experiment story's acceptance with and without its document.

## Notes
