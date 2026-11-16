---
id: T-3419
type: task
nature: remediation
title: The design and the guides describe the running close-out's head and step, the verify_running reason, and a report from before the head moved, and I-0384 is closed
status: backlog
parent: S-0799
owner: alex
created: 2026-11-16T01:47:02Z
updated: 2026-11-16T01:47:02Z
transitions: []
stream: S-0799
tags: [flai, docs, verify]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0384-after-a-sync-moves-a-story-s-head-the-review-page-flai-queue-and-the-acceptance.md, design/issues/summary.md]
after: [T-3409, T-3415, T-3416, T-3417, T-3418]
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
