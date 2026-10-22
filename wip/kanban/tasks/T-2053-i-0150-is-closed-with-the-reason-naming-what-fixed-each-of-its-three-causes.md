---
id: T-2053
type: task
nature: remediation
title: I-0150 is closed with the reason naming what fixed each of its three causes
status: backlog
parent: S-0467
owner: alex
created: 2026-10-22T21:01:26Z
updated: 2026-10-22T21:01:26Z
transitions: []
stream: S-0467
tags: [flai]
touches: [design/issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md, design/issues/summary.md]
after: [T-2049, T-2052]
---
# T-2053 I-0150 is closed with the reason naming what fixed each of its three causes

## Work

Waits for the prompt and docs tasks, so that the issue closes on the whole remedy.

From the story's worktree, run `flai issue close I-0150 --reason "<reason>"` with a reason that names the ADR and the three remedies: the orchestrator's acceptance lands on a dirty main and refuses only a change that would overwrite an uncommitted file, and finishes a stopped acceptance of its own; a rebase it cannot settle is aborted and the story sent back to its agent; and a thread the orchestrator alone wrote no longer gates the story. `flai issue close` writes the issue file and `design/issues/summary.md`; edit neither by hand. Then tick the story's second criterion with `flai criteria tick S-0467 2`, and the first once every code task's test passes, with `flai criteria tick S-0467 1`.

## Done when

- I-0150's status is closed with the reason above, `design/issues/summary.md` no longer lists it, and both are committed on the story branch.
- Both acceptance criteria of S-0467 are ticked through flai.

## Notes

Layer 5, the last. Issue commands write under the current checkout: run them in the story's worktree, not the main checkout.
