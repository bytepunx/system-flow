---
id: S-0822
type: story
nature: remediation
title: A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it
status: backlog
owner: alex
created: 2026-11-17T07:13:18Z
updated: 2026-11-27T23:56:52Z
transitions: []
tags: [flai, flaiover]
touches: [flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0396-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: planner
      seconds: 17028
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 75
          output: 4474
          cache_read: 9673105
          cache_write: 27338
          cost: 2.2431
    - kind: orchestrator
      seconds: 233
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 63
          output: 436
          cache_read: 2555785
          cache_write: 374681
          cost: 0.9125
        - model: claude-sonnet-5-5
          input: 20
          output: 1495
          cache_read: 188416
          cache_write: 79678
          cost: 0.216
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: flai
    at: 2026-11-17T07:13:18Z
  value: 75
  by: planner-S-0822
  at: 2026-11-17T08:23:00Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T06:15:00Z
  basis: flai forecast's 41m (120 s per unit of size times 20) raised to 1h30m for five tasks in three layers, a cancel path through runJob with race tests in flai/internal/serve, and a close-out of about 12m like S-0847's.
  by: planner-S-0822
  at: 2026-11-27T23:56:52Z
---
# S-0822 A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it

## Goal

This story remediates [I-0396](../../../design/issues/I-0396-a-rebase-and-resolve-job-that-passes-records-stayed-in-review-without-reading.md), "A rebase and resolve job that passes records stayed in review without reading the story, and a move out of review neither stops the job nor is noticed by it". The issue recommends this solution:

When a rebase and resolve job passes, `settle` in `flai/internal/serve/resolvejob.go` reads the story again, as `sendBack` does. For a story no longer in review it records an outcome of its own, not `stayed in review`, and a reason that names the state the story is in, who moved it there and when, and the command that moves it back to review, since its close-out passed at its head. The review page shows that outcome as it shows the others. Decide in the story, with the operator, whether a move out of review while the job runs should also stop the job; record the decision in an ADR refining ADR-0313 if it changes what ADR-0313 says. A test in the resolve job's tests moves the story out of review while the job runs and reads the record.

## Acceptance criteria
- [ ] The cause I-0396 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0396 is closed with `flai issue close I-0396 --reason` saying what fixed it

## Tasks
- T-3520 settle reads the story again and records left review for a story no longer in review
- T-3521 The review page shows a rebase and resolve job that ended with its story out of review
- T-3522 The design and the user docs name the left review outcome, and I-0396 is closed
- T-3613 An ADR refining ADR-0313 §6 records that a move out of review stops a running rebase and resolve job
- T-3614 flai serve stops a running rebase and resolve job whose story left review and ends it stopped, outcome left review

## Notes

Cost of delay inputs set by flai from I-0396. time_lost_per_cycle 30m carried over from I-0396's Impact section. Its Impact time_lost_per_cycle was taken rather than the 2h30m derived from its cost and count.

### Planning

The decision the goal asks for is made: TH-1274 settled `stop`. The orchestrator recommended it and flai serve confirmed it after 48h. When the story leaves review while a rebase and resolve job runs, flai serve stops the job and ends it `stopped`, outcome `left review`, with a reason that names the state, who moved the story, and when. This changes ADR-0313 §6, "A story that left review while the job ran is left as it is", so the story writes an ADR that refines ADR-0313. TH-1275 asked for this replan.

Touches, and where each came from:

- `flai/internal/serve/resolvejob.go`, `flai/internal/serve/resolvejob_test.go`: declared. The goal names `settle` and the resolve job's tests.
- `flai/internal/closeout/closeout.go`, `flai/internal/closeout/closeout_test.go`: declared, from the code layout. The job's outcomes, `stayed in review` and `sent back`, are defined there.
- `flai/internal/serve/verifyjob.go`, `flai/internal/serve/verifyjob_test.go`: added for the `stop` answer, from the code layout. `look` and `runJob` there run every job and hold its context, so the stop is made there. `verifyjob_container.go` is left out, since the container close-out already runs every step under the job's context.
- `flaiover/src/lib/review.ts`, `flaiover/src/lib/review.test.ts`, `flaiover/src/lib/components/Review.svelte`, `flaiover/src/lib/components/Review.svelte.test.ts`: declared. The goal says the review page shows the outcome.
- `design/adrs`, `design/adrs/README.md`: added for the `stop` answer, from the design. The folder touch stays because `flai adr new` numbers and names the ADR file when the story's agent writes it, so no task can name it yet. The README is the ADRs' index.
- `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`: declared, from the design. They describe the job's record and the review page's job box.
- The I-0396 file and `design/issues/summary.md`: declared. The second criterion closes the issue.

`flai touches suggest` lists co-change paths across 1844 of 4304 commits, the operators' docs, the reference, and the workflow design among them. None is taken: the change is in the resolve job alone, and adds no command, flag, or setting.

Figures:

- Cost of delay value: 75 USD a week, as `flai cod` gives it from the operator's input, 30m lost per 168h cycle at 150 USD an hour. It stands.
- Forecast: 1h30m, raised from the 41m `flai forecast` gives (120 s per unit of size, size 20). The figure counts touches, not tasks. The plan has five tasks in three layers, and T-3614 threads a cancel through `runJob` with race tests in `flai/internal/serve`. The close-out alone took about 12m for S-0847 today. Delivery is moved from flai's 2026-11-28T05:26Z by the same 49m, to 2026-11-28T06:15Z.

The story has no `after`: it reads nothing another open story writes. S-0824 and S-0846 change `resolvejob.go` too; a trial merge judges that overlap.
