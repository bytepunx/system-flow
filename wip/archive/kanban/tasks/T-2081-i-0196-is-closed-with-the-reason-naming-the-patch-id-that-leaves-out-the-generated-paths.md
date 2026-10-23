---
id: T-2081
type: task
nature: improvement
title: I-0196 is closed with the reason naming the patch id that leaves out the generated paths
status: done
parent: S-0487
owner: alex
created: 2026-10-23T18:05:30Z
updated: 2026-10-23T18:27:31Z
transitions:
  - to: ready
    at: 2026-10-23T18:27:15Z
    by: agent-S-0487
  - to: in-progress
    at: 2026-10-23T18:27:16Z
    by: agent-S-0487
  - to: done
    at: 2026-10-23T18:27:31Z
    by: agent-S-0487
stream: S-0487
tags: [issues]
touches: [design/issues/I-0196-a-close-out-s-second-round-runs-go-test-again-when-the-rebase-regenerates-design-issues-summary-md-which-changes-the-branch-s-patch-id-so-a-close-out-main-outruns-still-takes-about-11-minutes-with-the-queue-on.md, design/issues/summary.md]
after: [T-2078, T-2079, T-2080]
usage:
  source: log
  seconds: 15
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 3119
      cache_read: 696125
      cache_write: 28912
      cost: 0.3886
---
# T-2081 I-0196 is closed with the reason naming the patch id that leaves out the generated paths

## Work

In the story's worktree, run `flai issue close I-0196 --reason` with a reason that names the fix: `flai verify`'s patch id now leaves out the paths flai regenerates at sync, so a close-out round after a rebase that regenerated `design/issues/summary.md` reuses the tiers main's new paths do not select, and names the ADR T-2079 wrote. The command rewrites the issue file and `design/issues/summary.md`.

Waits for T-2078, T-2079, and T-2080: the issue closes when the fix, its ADR, and its documentation are on the branch, and its reason names the ADR.

## Done when

I-0196 reads `status: closed` with the reason, `design/issues/summary.md` no longer lists it, and `flai check --strict` passes.

## Notes

The close of the issue is this story's second acceptance criterion; the story's agent ticks it with `flai criteria tick` once the close is committed.
