---
id: T-2531
type: task
nature: remediation
title: The release package writes a changelog as the base branch has it with a branch's own added lines under its top Unreleased
status: done
parent: S-0570
owner: alex
created: 2026-11-07T03:48:02Z
updated: 2026-11-07T04:34:44Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:31Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:18:45Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:34:44Z
    by: agent-S-0570
stream: S-0570
tags: [cli, release]
touches: [flai/internal/release/changelog.go, flai/internal/release/changelog_test.go, flai/internal/release/release.go]
after: [T-2524]
usage:
  source: log
  seconds: 959
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 6463
      cache_read: 2850131
      cache_write: 43054
      cost: 1.025
---
# T-2531 The release package writes a changelog as the base branch has it with a branch's own added lines under its top Unreleased

## Work

Waits for T-2524: the ADR decides whether this merge is the fix, and what it keeps. If the ADR picks another fix, rewrite this task to match before you work it.

In `flai/internal/release/changelog.go`, new, write a pure function. It takes the base branch's changelog text and the lines a branch added. It returns the base's text with those lines at the end of the top `## Unreleased` section, its blank lines kept as `bumpVersionAndChangelog` writes them. It creates an empty `## Unreleased` above the first version heading when the base has none. It leaves out every added line the base already holds, in any section, so that a released line is never added again. Use the parsing `release.go` already has (`unreleasedHeading`, `cutUnreleased`, `firstSection`), moving a helper to `changelog.go` when both files need it.

In `changelog_test.go`, table tests:

- **Instance 1 of I-0251.** The base has `## Unreleased` (empty) above `## 1.1.17`, which holds the branch's old neighbours. The branch added one line. The result has one `## Unreleased`, holding only the branch's line, and `## 1.1.17` is as the base has it.
- **Instance 2.** The base has two `## 1.1.3` headings, and the branch's line must not land under `## 1.1.18`. The base's headings are all kept.
- **The plain cases.** No `## Unreleased` on the base, a line the base already holds, and no added lines, where the base's text comes back unchanged.

## Done when

- The function and its tests are in `flai/internal/release`.
- `flai test flai/internal/release` passes.

## Notes

- Drafted by the planner. A pure function keeps the git work in `storygit`, which the next task's wiring calls through a hook, so `storygit` does not import `release`.
