---
id: T-2138
type: task
nature: remediation
title: I-0184 is closed with the reason naming the kept commits, and the story's criteria are ticked
status: done
parent: S-0500
owner: alex
created: 2026-10-26T04:27:11Z
updated: 2026-10-26T04:43:11Z
transitions:
  - to: ready
    at: 2026-10-26T04:42:50Z
    by: agent-S-0500
  - to: in-progress
    at: 2026-10-26T04:42:50Z
    by: agent-S-0500
  - to: done
    at: 2026-10-26T04:43:11Z
    by: agent-S-0500
stream: S-0500
tags: [flai]
touches: [design/issues/I-0184-every-stats-get-and-flai-stats-runs-one-git-log-for-claims-drift-about-450-ms-of-a-900-ms-warm-answer.md, design/issues/summary.md]
after: [T-2137]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 3324
      cache_read: 557863
      cache_write: 20752
      cost: 0.3171
---
# T-2138 I-0184 is closed with the reason naming the kept commits, and the story's criteria are ticked

## Work

In the story's worktree, run `flai issue close I-0184 --reason "<what fixed it>"`, the reason saying that `storygit.Committed` keeps the stories' commits per repository root for the life of the process, keyed by `git show-ref --heads`, and reads the log again only when a branch tip moves, with the warm `stats.get` figure T-2137 measured. `flai issue close` regenerates `design/issues/summary.md`. Then tick the story's two criteria with `flai criteria tick S-0500 1` and `flai criteria tick S-0500 2`, each once verified: the first by T-2136's tests and T-2137's measurement, the second by the closed issue.

This task waits for T-2137, since the reason cites the measured figure and the criteria are verified by the two tasks before it.

## Done when

- I-0184's file under `design/issues/` has `status: closed` and a reason that names the kept commits and the measured warm figure
- `design/issues/summary.md` no longer lists I-0184 among the open issues
- Both criteria of S-0500 are ticked through flai
- `flai check --strict` passes scoped to the story

## Notes

Issue commands write under the current checkout, so they run in the story's worktree, not the main checkout.
