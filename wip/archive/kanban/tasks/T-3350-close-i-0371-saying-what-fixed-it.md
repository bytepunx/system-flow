---
id: T-3350
type: task
nature: remediation
title: Close I-0371 saying what fixed it
status: done
parent: S-0781
owner: alex
created: 2026-11-13T22:54:13Z
updated: 2026-11-14T17:45:09Z
transitions:
  - to: ready
    at: 2026-11-14T17:44:48Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T17:44:48Z
    by: system-flow
  - to: done
    at: 2026-11-14T17:45:09Z
    by: agent-S-0781
stream: S-0781
tags: [issues]
touches: [design/issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [T-3343, T-3345, T-3348, T-3349]
usage:
  source: log
  seconds: 21
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 2790
      cache_read: 1176282
      cache_write: 28812
      cost: 0.4907
---
# T-3350 Close I-0371 saying what fixed it

## Work

Close I-0371 with `flai issue close I-0371 --reason "<reason>"`. The reason names the cause, S-0777's test run leaking `story/S-0001`, `story/S-0002`, `story/S-0003`, and `story/S-0005` into the main repository, and each fix: the git ceiling in the tests' `TestMain` (T-3343), the removal of the four branches (T-3345), and the ADR's scoped check (T-3348, T-3349), or that the ADR decided the root fix alone.

It waits for T-3343, T-3345, T-3348, and T-3349, since the reason names what each did. A task among them that was cancelled is left out of the reason. On I-0370, which the same four branches caused, say in this task's `## Notes` whether the removal cleared it. Its closing is the operator's call.

## Done when

- I-0371 is closed with a reason naming the cause and each fix.
- `design/issues/summary.md` no longer lists it.

## Notes

I-0371 closed in 67ca39abf. On I-0370: the removal cleared its cause. After T-3345, the main checkout's `flai check` reports no `story.unaccepted`, 0 errors and 4 advisory warnings that `--strict` passes over, so the repository check in `scripts/flai-test.sh` no longer fails on the four branches. `make flai-test` was not run in this story. I-0370 stays open: closing it is the operator's call.
