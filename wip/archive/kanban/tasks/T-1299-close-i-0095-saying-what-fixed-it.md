---
id: T-1299
type: task
nature: remediation
title: Close I-0095 saying what fixed it
status: done
parent: S-0317
owner: alex
created: 2026-10-08T00:03:24Z
updated: 2026-10-08T00:17:28Z
transitions:
  - to: ready
    at: 2026-10-08T00:17:21Z
    by: agent-S-0317
  - to: in-progress
    at: 2026-10-08T00:17:22Z
    by: agent-S-0317
  - to: done
    at: 2026-10-08T00:17:28Z
    by: agent-S-0317
stream: S-0317
tags: [flai]
touches: [design/issues/I-0095-an-agent-that-ended-to-wait-for-an-answer-is-recorded-as-failed-and-never-started-again-when-the-answer-comes-before-flai-serve-s-next-look.md, design/issues/summary.md]
after: [T-1297, T-1298]
usage:
  source: log
  seconds: 6
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 4
      output: 951
      cache_read: 160190
      cache_write: 7557
      cost: 0.1019
---
# T-1299 Close I-0095 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0095 --reason "<reason>"`, which marks the issue closed and regenerates `design/issues/summary.md`. The reason says what fixed it: since S-0317, flai serve judges a run whose agent asked on a thread of its story during the run, and was answered before the end was judged, as `asked`, and starts it again at once in its session on the answer, without spending an automatic restart; name the test that reproduces the case. Add to the issue's `## Remediation` that the first direction was taken and why the second was not.

It waits for T-1297 and T-1298 because the issue is closed only once the fix and its documentation are in.

## Done when

- The front matter of I-0095 says `status: closed`, and its file records the reason.
- `design/issues/summary.md` no longer lists I-0095 as open (`flai issue list` agrees).
- `flai check --strict` scoped to the story is clean.

## Notes
