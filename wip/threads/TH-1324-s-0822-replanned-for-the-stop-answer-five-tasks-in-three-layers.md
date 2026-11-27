---
id: TH-1324
title: "S-0822 replanned for the stop answer: five tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md
  item: S-0822
status: open
participants: [planner-S-0822]
created: 2026-11-27T23:57:06Z
updated: 2026-11-27T23:57:06Z
---

# TH-1324 S-0822 replanned for the stop answer: five tasks in three layers

On wip/kanban/stories/S-0822-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md.

## Entries

### 2026-11-27T23:57:06Z planner-S-0822
S-0822 is replanned for TH-1274's answer, `stop`, as TH-1275 asked. It is complete for `flai promote --drafts`. `flai check --strict` reports nothing on it or its tasks.

## Tasks and layers

| Layer | Task | What it does | after |
|-------|------|--------------|-------|
| 1 | T-3520 | `settle` reads the story again and records `left review` (kept as it was) | none |
| 1 | T-3521 | The review page shows `left review`, passed or stopped (kept as it was) | none |
| 1 | T-3613 | New: an ADR refining ADR-0313 §6, saying that a move out of review stops the job, and what the job records | none |
| 2 | T-3614 | New: flai serve cancels a running resolve job whose story left review. The job ends `stopped` at its step, outcome `left review`, with nothing sent back | T-3520, T-3613 |
| 3 | T-3522 | The design and user docs, now with the stop and the ADR; closes I-0396 (Work and `after` revised) | T-3521, T-3614 |

T-3614 waits for T-3520 because both change `resolvejob.go` and its tests, and it uses T-3520's outcome and reason. It waits for T-3613 because the ADR fixes the record's fields.

## Story

- Touches added: `flai/internal/serve/verifyjob.go`, `flai/internal/serve/verifyjob_test.go`, `design/adrs/README.md`, and the folder `design/adrs`. The folder stays because `flai adr new` names the ADR file. The reason is under `### Planning`.
- Forecast: 1h30m, delivery 2026-11-28T06:15Z. That is up from flai's 41m, for five tasks, the cancel path with race tests, and a close-out of about 12m.
- Cost of delay: 75 USD a week. It stands.
- `after`: none.

## Assumptions

- A left-review stop is told apart from a stop with flai serve, which still ends a job `interrupted` at the next start, and from the time limit.
- A stop during the rebase aborts a rebase the job started. A rebase the job found waiting is left as it was, as ADR-0313 §3 does.
- Only resolve jobs are stopped. A close-out job the story's agent asked for runs on, since that agent waits on its verdict. T-3613 records this in the ADR, and can recommend otherwise if the agent writing it finds a reason.
- I would not split, merge, or drop any task.
