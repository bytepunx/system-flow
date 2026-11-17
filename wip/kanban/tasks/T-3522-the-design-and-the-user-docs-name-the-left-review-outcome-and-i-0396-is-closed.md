---
id: T-3522
type: task
nature: remediation
title: The design and the user docs name the left review outcome, and I-0396 is closed
status: backlog
parent: S-0822
owner: alex
created: 2026-11-17T08:21:08Z
updated: 2026-11-17T08:21:08Z
transitions: []
stream: S-0822
tags: [docs]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0396-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md, design/issues/summary.md]
after: [T-3520, T-3521]
---
# T-3522 The design and the user docs name the left review outcome, and I-0396 is closed

## Work

Write down what the code tasks changed, in the same story, as `documentation.md` asks.

- `design/system/flai-cli.md`, the rebase and resolve job's `settle` and its state fields: add the outcome `left review` and what its reason names. Fix the drift found there too: a rebase that stops ends the job with verdict `stopped` and `stopped_at` `rebase`, not a verdict `rebase`.
- `design/system/flaiover-dashboard.md`: the job box's ends and the job's `outcome` field gain `left review`.
- `docs/users/flai.md`, `flai accept --resolve`: say what the job records when the story leaves review while it runs.
- `docs/users/flaiover.md`, the review page's Rebase and resolve box: the same, as the operator sees it.

Then close the issue with `flai issue close I-0396 --reason`, naming the fix and the test that reproduces it.

It waits for T-3520 and T-3521, because it describes the outcome and the wording those tasks write.

## Done when

- Each of the four documents names the `left review` outcome as the code records and shows it.
- I-0396 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it.
- The markdown lint passes with `flai test` on the changed paths.

## Notes
