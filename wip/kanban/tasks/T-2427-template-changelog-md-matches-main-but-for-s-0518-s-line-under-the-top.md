---
id: T-2427
type: task
nature: feature
title: template/CHANGELOG.md matches main but for S-0518's line under the top Unreleased, after the rebase across release 1.1.18
status: done
parent: S-0518
owner: alex
created: 2026-11-02T04:12:50Z
updated: 2026-11-02T04:13:13Z
transitions:
  - to: ready
    at: 2026-11-02T04:13:02Z
    by: agent-S-0518
  - to: in-progress
    at: 2026-11-02T04:13:03Z
    by: agent-S-0518
  - to: done
    at: 2026-11-02T04:13:13Z
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
      input: 3
      output: 672
      cache_read: 372799
      cache_write: 6561
      cost: 0.1405
---

# T-2427 template/CHANGELOG.md matches main but for S-0518's line under the top Unreleased, after the rebase across release 1.1.18

## Work

The sync after S-0422's acceptance and release 1.1.18 rebased S-0518's changelog line into the `## 1.1.18` section, and the union merge dropped main's second `## 1.1.3 - 2026-10-28` heading, which the sync reported as `undoes_main` (I-0251). Take main's `template/CHANGELOG.md` and add S-0518's line under the top `## Unreleased`.

## Done when

`git diff main -- template/CHANGELOG.md` shows only S-0518's line under the top `## Unreleased`, the markdown tier passes, and the sync reports nothing in `undoes_main`.

## Notes

T-2425 fixed the same thing after release 1.1.17.
