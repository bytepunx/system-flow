---
id: T-3418
type: task
nature: remediation
title: The review page and the merge queue panel show the close-out running at a story's new head, its job and step, beside a last report from before the head moved
status: backlog
parent: S-0799
owner: alex
created: 2026-11-16T01:46:36Z
updated: 2026-11-16T01:46:36Z
transitions: []
stream: S-0799
tags: [flaiover, merge-queue, verify]
touches: [flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/MergeQueue.svelte, flaiover/src/lib/components/MergeQueue.svelte.test.ts, design/components/flaiover/src/lib/components/MergeQueue.examples.json]
after: [T-3415, T-3416, T-3417]
---
# T-3418 The review page and the merge queue panel show the close-out running at a story's new head, its job and step, beside a last report from before the head moved

## Work

The review page (`flaiover/src/lib/components/Review.svelte`, Verification section, lines 1270-1290) shows `verify.commit` and "Stopped at <step>" from `verify.status`, and the merge queue panel (`MergeQueue.svelte`) shows the queue's reason text. Neither knows of a close-out running at the new head.

- In `flaiover/src/lib/review.ts`, give `VerifyReport` (lines 200-214) `moved` and `running` as T-3417 answers them, and have `verifyView` (lines 283-316) give the row a note: "This report is from before the head moved to <head>" when `moved`, and "Close-out running at <head>: job <id>, at <step>, since <time>" when `running`.
- In `Review.svelte`, show that note in the Verification row beside the last report's outcome, and the preview's `verify_running` from T-3416 beside its blockers. Leave Accept's rule as it is: the blockers decide it.
- In `MergeQueue.svelte`, give `QueueReason` (lines 3-11) the `job` field T-3415 adds, and show a `verify_running` reason as the running job and its step, not as a failure.
- Add a `verify_running` reason to `design/components/flaiover/src/lib/components/MergeQueue.examples.json`, so the component's examples keep step with its props.
- Test it: `review.test.ts` gives the note for a `moved` report and for a `running` job; `Review.svelte.test.ts` renders it; `MergeQueue.svelte.test.ts` renders a `verify_running` reason with the job and step. The first fails before the change: it reproduces I-0384 on the review page.

This task waits for T-3415, T-3416, and T-3417, whose JSON it reads.

## Done when

- The review page and the merge queue panel name the close-out job and its step while one runs at a new head, and say when the last report is from before the head moved.
- `flai test` on the flaiover files it changes passes, svelte-check included.

## Notes
