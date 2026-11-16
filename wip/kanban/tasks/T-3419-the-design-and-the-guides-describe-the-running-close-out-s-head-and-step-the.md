---
id: T-3419
type: task
nature: remediation
title: The design and the guides describe the running close-out's head and step, the verify_running reason, and a report from before the head moved, and I-0384 is closed
status: done
parent: S-0799
owner: alex
created: 2026-11-16T01:47:02Z
updated: 2026-11-16T03:40:27Z
transitions:
  - to: ready
    at: 2026-11-16T02:43:42Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T03:36:19Z
    by: agent-S-0799
  - to: done
    at: 2026-11-16T03:40:27Z
    by: agent-S-0799
stream: S-0799
tags: [flai, docs, verify]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md, design/issues/summary.md, design/system/flaiover-dashboard.md]
after: [T-3409, T-3415, T-3416, T-3417, T-3418]
usage:
  source: log
  seconds: 248
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 94
      output: 32225
      cache_read: 4911769
      cache_write: 238107
      cost: 2.9866
---
# T-3419 The design and the guides describe the running close-out's head and step, the verify_running reason, and a report from before the head moved, and I-0384 is closed

## Work

Bring the documents in step with what the tasks before this one changed, in the same story as the code:

- `design/system/flai-cli.md` § The close-out job, the bullet **The state.**: the job's file names the head it runs at from its start, and `flai queue`, the preview, `flai verify --last`, and `verify.status` read a running job's head and step.
- `docs/users/flai.md` § Accept and release › The merge queue: add the `verify_running` row to the reasons table, and say what a `verify_failed` from before the head moved reads. In § Verify a story before review › The close-out job, and where `--last` is described, give the two lines `--last` now prints. In § Accept and release, give the preview's `verify_running`.
- `docs/users/flaiover.md` § Board › Reviewing a story: the Verification row and the queue panel name the running close-out and its step.
- Close I-0384 with `flai issue close I-0384 --reason`, naming the tasks that fixed it, and let flai regenerate `design/issues/summary.md`.

S-0796 changes `design/system/flai-cli.md` and `docs/users/flai.md` too. Sync before you start.

This task waits for every other task of the story, whose behaviour it describes.

## Done when

- Each section above says what the code now does, and the markdown lint passes on the three documents.
- I-0384 is closed, and the story's two criteria can be ticked.

## Notes
