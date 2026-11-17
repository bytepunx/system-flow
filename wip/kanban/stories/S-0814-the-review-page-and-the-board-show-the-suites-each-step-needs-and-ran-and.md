---
id: S-0814
type: story
nature: feature
title: The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream
status: backlog
parent: E-0034
owner: alex
created: 2026-11-16T23:42:10Z
updated: 2026-11-17T09:16:30Z
transitions: []
tags: [flaiover]
topics: [verify, workflow]
touches: [flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/verify_suites.go, flai/internal/hostapi/verify_suites_test.go, flai/internal/hostapi/resolve_status.go, flai/internal/hostapi/resolve_status_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flaiover/src/lib/components/VerifySuites.svelte, flaiover/src/lib/components/VerifySuites.svelte.test.ts, design/components/flaiover/src/lib/components/VerifySuites.examples.json, flaiover/src/lib/components/LogFeed.svelte, flaiover/src/lib/components/LogFeed.svelte.test.ts, design/components/flaiover/src/lib/components/LogFeed.examples.json, "flaiover/src/routes/api/items/[id]/resolve/+server.ts", "flaiover/src/routes/api/items/[id]/resolve/resolve.test.ts", flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/lib/server/board.ts, flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [S-0810, S-0807, S-0753]
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
      seconds: 81
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 47
          output: 17830
          cache_read: 4742358
          cache_write: 132280
          cost: 2.3635
    - kind: orchestrator
      seconds: 21
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 6
          output: 66
          cache_read: 318646
          cache_write: 52546
          cost: 0.1138
cost_of_delay:
  value: 71.32
  by: planner-S-0814
  at: 2026-11-16T23:52:43Z
forecast:
  duration: 1h1m
  delivery: 2026-11-17T15:28:00Z
  basis: "Its own forecast of 1h1m; 11th in the pull order with an in-progress limit of 3, behind S-0806, S-0831, S-0833, S-0754, S-0811, S-0809, S-0756, S-0812, S-0808, S-0807, S-0810, S-0241 and S-0813."
  by: flai
  at: 2026-11-17T09:16:30Z
finalized:
  by: orchestrator
  at: 2026-11-16T23:56:27Z
---
# S-0814 The review page and the board show the suites each step needs and ran, and follow a state's job in the activity stream

## Goal

[ADR-0312](../../../design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md) §4. The review page's verification section shows the last report's tiers, and after S-0804 a Rebase and resolve job's stream.

## Acceptance criteria

- [ ] The verification section groups the report's tiers by suite and says which suites the story's next transitions need and which have not passed at its head.
- [ ] A job a state's `runs:` started is followed in the collapsible activity stream S-0804 adds, with its current step, its skipped cluster steps and why, and how it ended.
- [ ] The host API gives the dashboard the suites of each transition, state, and point, and the board's column for a state with `runs:` says what it runs.
- [ ] The component examples and the dashboard design are in step.

## Tasks
- T-3476 board.get gives the dashboard the suites each transition's verify names, each state's runs, and verify_at's batch, publish, and ci
- T-3478 verify.status gives the report's tiers by suite, the suites the story's next transitions need, and which of them have not passed at its head
- T-3480 resolve.status gives the dashboard a job a state's runs started, with its state, its current step, its skipped cluster steps and why, and how it ended
- T-3481 The review page's verification section groups the report's tiers by suite in a VerifySuites component and names the suites the next transitions need that have not passed at the head
- T-3483 The board's column for a state with runs says what it runs, the suites and the cluster steps, from what board.get gives
- T-3485 The review page follows a job a state's runs started in the collapsible activity stream, with its current step, its skipped cluster steps and why, and how it ended
- T-3486 The dashboard design and the user guide describe the suites on the review page and the board, the job a state's runs started in the stream, and what board.get, verify.status, and resolve.status give

## Notes

- After S-0810: the page follows the job a state's runs key starts, which S-0810 adds.
- After S-0804: the stream is the component S-0804 adds to the review page.

### Planning

Written by planner-S-0814 on 2026-11-16.

`after`, each with what this story reads from it:

- S-0810: the job a state's `runs:` starts, its record with each step as it starts and ends, and the cluster steps it skips and why. S-0810 is itself after S-0804 and S-0806, so the Rebase and resolve job's record, `resolve.status`, and the review page's activity stream (S-0804's T-3430 and T-3431), and the parsed `verify:`, `runs:`, and `verify_at:` keys (S-0806), are on main before this story starts.
- S-0807: the verify report's record of each suite's tiers and the head they passed at, and the suites the transitions out of a story's state need. The verification section groups by the first and names the second.
- S-0753: the board's columns, with their kinds, drawn from what flai gives in place of the dashboard's own list of states. A column for a project's own state with `runs:`, such as S-0815's integration state, exists only once the board draws its columns from flai.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/hostapi/writes_host.go`, `writes_test.go` | layout: `verify.status` is answered there (`verifyStatus`); co-change with `Review.svelte` |
| `flai/internal/hostapi/verify_suites.go`, `verify_suites_test.go` | layout: a new file for the suites the next transitions need and which have not passed at the head, kept out of `writes_host.go` |
| `flai/internal/hostapi/hostapi.go`, `hostapi_test.go` | layout: `board.get` is answered there |
| `flai/internal/workitem/boardview.go`, `boardview_test.go` | layout: `NewBoardView` builds the answer `board.get` gives |
| `flai/internal/hostapi/resolve_status.go`, `resolve_status_test.go` | design: S-0804's T-3430 creates `resolve.status`, which gives the dashboard the job; this story gives it a `runs:` job's state and skipped steps |
| `flaiover/src/lib/review.ts`, `review.test.ts`, `components/Review.svelte`, `Review.svelte.test.ts` | layout: the verification section (`verifyView`, `verify-section`); co-change 33% and 13% |
| `flaiover/src/lib/components/VerifySuites.svelte`, its test, and `design/components/flaiover/src/lib/components/VerifySuites.examples.json` | design: the grouped view as a component of its own, so that it has examples, which `Review.svelte` has none of (criterion 4) |
| `flaiover/src/lib/components/LogFeed.svelte`, its test, and `LogFeed.examples.json` | layout: the collapsible live log the stream is built on; its step and outcome props gain the skipped steps |
| `flaiover/src/routes/api/items/[id]/resolve/+server.ts`, `resolve.test.ts` | design: S-0804's T-3431 route for the job, which this story follows for a `runs:` job too |
| `flaiover/src/routes/board/+page.svelte`, `board.svelte.test.ts`, `flaiover/src/lib/server/board.ts`, `lanes.ts`, `lanes.test.ts` | layout: the lane header and the board's load; co-change 13% and 9% |
| `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md` | criterion 4 and the documentation convention; co-change 25% and 24% |

`flai touches suggest` from the four layout files also listed `flaiover/src/routes/items/[id]/+page.svelte` (30%) and `BoardCard.svelte` (15%). Both are left out: the item page and the cards show neither the suites nor the job.

Figures:

- Forecast 1h1m, delivery 2026-11-17T11:14:00Z, as `flai forecast` gives it: 110 s per unit of size times size 33 (4 criteria, 29 touches), 16th in the pull order. It stands unadjusted. Two tasks may shrink if S-0804 or S-0810 already give part of what T-3480 and T-3485 add, and the three stories this one builds on may name their interfaces otherwise than predicted. Neither is a measured reason to move the figure either way.
- Cost of delay 71.32 USD a week, as `flai cod` gives it after the forecast: this story's share of E-0034's 1500.00 USD a week, from the operator's input on the epic, `time_lost_per_cycle` 10h, by its 1h1m of the 21h23m forecast over the epic's 20 open stories without inputs. The story has no inputs of its own, so no input was asked for.
