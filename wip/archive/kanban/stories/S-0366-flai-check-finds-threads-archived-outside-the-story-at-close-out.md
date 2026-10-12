---
id: S-0366
type: story
nature: improvement
title: "flai check finds `threads.archived` outside the story at close-out"
status: done
owner: alex
created: 2026-10-12T14:12:35Z
updated: 2026-10-12T15:48:40Z
transitions:
  - to: ready
    at: 2026-10-12T15:26:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T15:26:09Z
    by: agent-S-0366
  - to: review
    at: 2026-10-12T15:42:39Z
    by: agent-S-0366
  - to: done
    at: 2026-10-12T15:48:40Z
    by: alex
tags: [flai, check]
touches: [design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0135-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/summary.md]
after: [S-0365]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1008
  turns:
    - day: 2026-10-12
      ceremony: 1
      hand_edits: 1
      work: 31
  models:
    - model: claude-opus-5-5
      input: 140
      output: 39851
      cache_read: 6606428
      cache_write: 306898
      cost: 4.1943
  strategic:
    - kind: orchestrator
      seconds: 882
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 81
          output: 1164
          cache_read: 12543456
          cache_write: 42819
          cost: 3.173
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-12T15:05:05Z
  value: 12.5
  by: planner-S-0366
  at: 2026-10-12T15:05:52Z
forecast:
  duration: 23m
  delivery: 2026-10-13T03:23:00Z
  basis: "Its own forecast of 23m; 27th in the pull order with an in-progress limit of 5, behind S-0232, S-0337, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0304, S-0305, S-0306, S-0313, S-0349, S-0350, S-0351, S-0352, S-0353, S-0354, S-0355, S-0356, S-0357, S-0358, S-0359, S-0360, S-0363 and S-0364."
  by: flai
  at: 2026-10-12T15:17:00Z
finalized:
  by: orchestrator
  at: 2026-10-12T15:06:18Z
---
# S-0366 flai check finds `threads.archived` outside the story at close-out

## Goal

This story remediates [I-0135](../../../design/issues/I-0135-flai-check-finds-threads-archived-outside-the-story-at-close-out.md), "flai check finds `threads.archived` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0135 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0135 is closed with `flai issue close I-0135 --reason` saying what fixed it

## Tasks
- T-1451 Record in an ADR that a close-out records no threads.archived and a check scoped to a story leaves it out
- T-1452 Leave threads.archived outside the story out of a check scoped to a story, with a test that reproduces I-0135
- T-1453 Say in the design, the user guide, and the work-management convention that a close-out leaves out threads.archived
- T-1454 Close I-0135 saying what fixed it

## Notes

### Planning

Proposed fix, from I-0135's three instances: a check scoped to a story leaves out `threads.archived` outside the story. It already leaves out `item.archive` (ADR-0122) and `board.wip-limit` (ADR-0133) the same way, and S-0365 adds `issues.no-story`. The close-outs of S-0343, S-0289, and S-0342 each found it on TH-0389, which none of them touched. The orchestrator opened TH-0389 on S-0338 at 2026-10-08T20:55:51Z, after S-0338's acceptance had archived it, so ADR-0109's resolve-on-archive could not reach it. The thread awaited the operator, and threads live in `wip/` on the main branch, so no story branch can clear it. `flai check` in the main checkout still warns. A `threads.archived` on a thread anchored on the story or its tasks stays inside.

Alternative, for T-1451 to weigh: stop `threads.archived` firing on a thread opened after its item was archived, in every check. It fixes TH-0389's cause, but not a thread an archive that stopped part way left open.

Where the code is:

- The finding comes from `flai/internal/check/check.go`.
- The scoped check is `ScopeToStory` in `flai/internal/check/scope.go`.
- The issue is recorded by `recordOutside` in `flai/cmd/check.go`, whose help uses `threads.archived` as its example of an issue title.

Tasks, in layers:

1. T-1451, the ADR.
2. T-1452, the code and its tests, and T-1453, the design, user guide, and convention. Both wait for T-1451 for its number, and they share no path.
3. T-1454 closes I-0135, after T-1452 and T-1453.

The story waits for S-0365 (`after`). Both change the same rule list in `ScopeToStory`, the same docs, and each adds an ADR. In progress, S-0365 would hold S-0366 on overlap anyway; `after` says why.

Touches, and where each came from (none declared):

| Touch | Source |
|-------|--------|
| `flai/internal/check/scope.go`, `flai/cmd/check.go` | layout: where the finding is scoped and recorded |
| `flai/internal/check/scope_test.go`, `flai/cmd/check_test.go`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, `design/adrs/README.md` | co-change with the two above (`flai touches suggest`, 17-44%); S-0348 changed each |
| `design/system/continuous-improvement.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design: they list what a close-out leaves out; S-0348 changed each |
| `design/adrs` | design: the new ADR, as in each precedent |
| `design/issues/I-0135-flai-check-finds-threads-archived-outside-the-story-at-close-out.md`, `design/issues/summary.md` | the issue the story closes |

Left out:

- `flai/internal/check/check.go` (33% co-change): the fix scopes the finding and leaves its production alone. If T-1451 takes the alternative, the story's agent adds it and `check_test.go`.
- `flai/internal/check/check_test.go` and `flai/cmd/check_stats_test.go`: they test production and stats, which do not change.

Folder touch kept: `design/adrs`. T-1451 adds an ADR whose number is not known until it is written, and other open stories add ADRs too.

Forecast: 23m, delivery 2026-10-13T03:06:00Z, as `flai forecast` gives it: 78 s per unit of size over 57 large-band improvement stories, times size 17. Kept: S-0365, the same fix for `issues.no-story` with the same four tasks, was forecast at 23m, and S-0348 took 24m.

Cost of delay: 12.50 USD a week, as `flai cod` gives it from `time_lost_per_cycle: 5m` at 150 USD an hour. The orchestrator set that input on TH-0406, as recommended. Kept: it matches S-0365 and S-0348, close-out findings of the same kind.
