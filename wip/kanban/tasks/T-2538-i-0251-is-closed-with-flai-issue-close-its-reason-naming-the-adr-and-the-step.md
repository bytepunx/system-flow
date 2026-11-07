---
id: T-2538
type: task
nature: remediation
title: I-0251 is closed with flai issue close, its reason naming the ADR and the step that fixed it
status: ready
parent: S-0570
owner: alex
created: 2026-11-07T03:49:20Z
updated: 2026-11-07T04:16:40Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:40Z
    by: system-flow
stream: S-0570
tags: [issues]
touches: [design/issues/I-0251-a-story-branch-rebased-across-a-release-brings-back-the-changelog-s-old.md, design/issues/summary.md]
after: [T-2535, T-2536]
---
# T-2538 I-0251 is closed with flai issue close, its reason naming the ADR and the step that fixed it

## Work

Waits for T-2535 and T-2536: the issue closes only once the sync's, accept's, and the batch's rebases keep a story's changelog lines in place, and their tests pass.

In the story's worktree, run `flai issue close I-0251 --reason "<reason>"`. The reason names T-2524's ADR and the step after the rebase that writes each changelog again with the branch's lines under the top `## Unreleased`, and the tests that reproduce both instances. This is criterion 2. `flai issue close` writes `design/issues/summary.md` again.

## Done when

- I-0251's status is closed, with the reason.
- `design/issues/summary.md` no longer lists it, and `flai check --strict` reports no `issues.summary`.

## Notes

- Drafted by the planner. It runs beside T-2537, with no path in common.
