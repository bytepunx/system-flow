---
id: T-2106
type: task
nature: remediation
title: I-0204 is closed with flai issue close naming the fix, and the story's criteria are ticked
status: done
parent: S-0493
owner: alex
created: 2026-10-24T13:19:19Z
updated: 2026-10-24T13:27:07Z
transitions:
  - to: ready
    at: 2026-10-24T13:26:56Z
    by: agent-S-0493
  - to: in-progress
    at: 2026-10-24T13:26:56Z
    by: agent-S-0493
  - to: done
    at: 2026-10-24T13:27:07Z
    by: agent-S-0493
stream: S-0493
tags: [flai]
touches: [design/issues/I-0204-flai-accept-fails-at-git-branch-d-when-the-story-branch-tracks-a-pushed-upstream-the-acceptance-rebase-left-behind-after-the-merge-into-main-is-done.md, design/issues/summary.md]
after: [T-2105]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 483
      cache_read: 158951
      cache_write: 5195
      cost: 0.083
---
# T-2106 I-0204 is closed with flai issue close naming the fix, and the story's criteria are ticked

## Work

Once T-2105 has made `flai accept` delete the story branch with `git branch -D` after the fast-forward and finish an acceptance run again from the state I-0204 left, close the issue from the story's worktree with `flai issue close I-0204 --reason "<what fixed it>"`, the reason naming `mergeStoryBranch`'s `-D` in `flai/cmd/branch.go` and the tests that pin it, which regenerates `design/issues/summary.md`. Then tick the story's two acceptance criteria with `flai criteria tick`, each once it is verified: the first by the tests of T-2105 passing, the second by the issue's status. Waits for T-2105 because the reason names the fix and the tests it added.

## Done when

- I-0204's status is `closed`, its reason names the fix, and `design/issues/summary.md` no longer lists it.
- Both acceptance criteria of S-0493 are ticked.
- `flai check --strict` is clean in the story's worktree.

## Notes
