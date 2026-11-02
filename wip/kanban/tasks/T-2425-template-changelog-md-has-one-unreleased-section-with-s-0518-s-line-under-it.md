---
id: T-2425
type: task
nature: feature
title: template/CHANGELOG.md has one Unreleased section, with S-0518's line under it, after the rebase across release 1.1.17
status: done
parent: S-0518
owner: alex
created: 2026-11-02T03:57:40Z
updated: 2026-11-02T03:58:00Z
transitions:
  - to: ready
    at: 2026-11-02T03:57:49Z
    by: agent-S-0518
  - to: in-progress
    at: 2026-11-02T03:57:50Z
    by: agent-S-0518
  - to: done
    at: 2026-11-02T03:58:00Z
    by: agent-S-0518
stream: S-0518
tags: []
touches: [template/CHANGELOG.md]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 548
      cache_read: 325854
      cache_write: 7454
      cost: 0.1358
---

# T-2425 template/CHANGELOG.md has one Unreleased section, with S-0518's line under it, after the rebase across release 1.1.17

## Work

Merge queue batch 20261102-035614 failed S-0518 at `markdown` (TH-0802): the rebase across release 1.1.17 union-merged a second `## Unreleased` into `template/CHANGELOG.md`, carrying S-0518's line and four lines the release had already moved under `## 1.1.17` (I-0251). Drop that block and put S-0518's line under the top `## Unreleased`.

## Done when

`template/CHANGELOG.md` differs from main only by S-0518's line under the top `## Unreleased`, the markdown tier passes on it, and the close-out passes.

## Notes

A fix found after review, so it is a task of its own rather than a reopened one.
