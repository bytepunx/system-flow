---
id: S-0296
type: story
nature: remediation
title: A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts
status: ready
owner: alex
created: 2026-10-06T11:44:51Z
updated: 2026-10-06T18:11:31Z
transitions:
  - to: ready
    at: 2026-10-06T11:49:31Z
    by: alex
tags: [flai]
topics: [orchestration]
touches: [flai/internal/serve/review_wait_test.go, flai/internal/serve/orchestrate.go, design/system/workflow.md, docs/operators/settings.md, docs/operators/index.md, design/issues/I-0088-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md, design/issues/summary.md]
after: [S-0221, S-0286]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 2h48m
    by: flai
    at: 2026-10-06T11:44:51Z
  value: 420
  by: planner-S-0296
  at: 2026-10-06T11:50:50Z
forecast:
  duration: 45m
  delivery: 2026-10-06T23:25:00Z
  basis: "Its own forecast of 45m; 8th in the pull order with an in-progress limit of 3, behind S-0226, S-0284, S-0223, S-0224, S-0227, S-0278, S-0229 and S-0295."
  by: flai
  at: 2026-10-06T18:11:31Z
finalized:
  by: alex
  at: 2026-10-06T11:45:52Z
---
# S-0296 A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts

## Goal

This story remediates [I-0088](../../../design/issues/I-0088-a-story-waits-in-review-for-the-operator-while-nothing-else-can-start-so-the-board-stands-idle-until-a-person-accepts.md), "A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts". The issue recommends this solution:

S-0221, already in ready, lets the orchestrator accept a story in review when the operator turns that permission on. S-0286 keeps the operator's acceptance for a story that changes a path Claude Code protects. Between them, the wait would remain only for those stories and for the ones the operator chooses to review. Until S-0221 is built, an agent in the operator's session can run the acceptance on the operator's word.

## Acceptance criteria
- [ ] The cause I-0088 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0088 is closed with `flai issue close I-0088 --reason` saying what fixed it

## Tasks
- T-1020 A flai serve test reproduces the idle board behind a story in review and shows the orchestrator's acceptance under accept_reviews clears it
- T-1021 The workflow design and the operators' guide say that accept_reviews ends the board's wait on the operator, and which stories still wait
- T-1022 I-0088 is closed with flai issue close, naming the orchestrator's acceptance under accept_reviews and the test that shows it

## Notes

Cost of delay inputs set by flai from I-0088. time_lost_per_cycle 2h48m: 1h24m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-10-06T10:32:27Z, 0.1 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0296 on 2026-10-06. The plan's thread, TH-0181, lists the tasks, their layers, and the assumptions. The operator took its recommendations on 2026-10-06:

- The story waits for S-0286 as well as S-0221.
- The operator turns `orchestrate` and `orchestration.permissions.accept_reviews` on for this project once both are accepted.
- The story is kept, not folded into S-0221.

The fix itself is S-0221's, accepted on 2026-10-06: the orchestrator accepts a story in review under `orchestration.permissions.accept_reviews`. S-0286 keeps the operator's acceptance for a story that changes a path Claude Code protects, so that the permission can be left on. What is left here is a test showing the board no longer idles behind a story in review, the design and guides saying so, and closing I-0088.

Touches:

- Declared: none; the story declared no touches, so `flai touches suggest` was started from the paths below.
- `flai/internal/serve/review_wait_test.go`: from the layout. The hold and orchestrator tests are `hold_test.go` and `orchestrate_test.go` in `flai/internal/serve`, and the launcher is what starts the held story once the review clears.
- `flai/internal/serve/orchestrate.go`: from the layout. The launcher looks at the orchestrator when items change and every minute; it changes only if the test shows the review noticed late.
- `design/system/workflow.md`: from the design, which describes acceptance and the review column.
- `docs/operators/settings.md`: from the design. It lists `accept_reviews`.
- `docs/operators/index.md` (co-change 27%): it describes the `orchestrate` host action.
- I-0088's file and `design/issues/summary.md` (co-change 9%): from the layout. `flai issue close` rewrites both.
- Not taken from `flai touches suggest`: `docs/users/flai.md` (68%), `design/system/flai-cli.md` (67%), `docs/users/flai-reference.md` (46%). This story adds no command or flag; S-0221 documents `flai accept --by orchestrator`.

Forecast: 45m, delivery 2026-10-07T04:00Z.

- `flai forecast` gave 18m: 116 s per unit of size over 12 done medium remediation stories, times size 9 (2 criteria, 7 touches).
- Raised to 45m. The test sets up a project with a story in review and a held ready story, and has a stand-in orchestrator accept under the permission. It must pass with the race detector, as S-0292's serve test did. The docs and the issue close are small.
- Delivery is not flai's 16:34Z, which leaves out the wait for S-0286. flai forecasts S-0286 done at 2026-10-07T02:56Z, 43rd in the pull order. Add 45m times a cycle factor of about 1.5, plus a margin, and delivery is 04:00Z. S-0286 is still a draft, so delivery holds only once the operator finalizes it.

Cost of delay: 420 USD/week, from `flai cod`, kept as computed.

- It comes from flai's input from I-0088: 2h48m lost per 168h cycle at 150 USD an hour.
- The input may undercount the loss. In S-0220's instance, ten held stories and three idle lanes waited 2h45m, against the 1h24m per occurrence the issue records. The input is the operator's to change, so the value is kept.
