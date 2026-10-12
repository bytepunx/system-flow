---
id: S-0365
type: story
nature: improvement
title: "flai check finds `issues.no-story` outside the story at close-out"
status: done
owner: alex
created: 2026-10-12T14:12:33Z
updated: 2026-10-12T15:17:00Z
transitions:
  - to: ready
    at: 2026-10-12T14:58:03Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T14:58:10Z
    by: agent-S-0365
  - to: review
    at: 2026-10-12T15:15:53Z
    by: agent-S-0365
  - to: done
    at: 2026-10-12T15:17:00Z
    by: orchestrator
tags: [flai, check]
touches: [design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1073
  turns:
    - day: 2026-10-12
      ceremony: 4
      hand_edits: 1
      work: 31
  models:
    - model: claude-opus-5-5
      input: 146
      output: 37473
      cache_read: 7636788
      cache_write: 340748
      cost: 4.5769
  strategic:
    - kind: orchestrator
      seconds: 617
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 68
          output: 1108
          cache_read: 10218151
          cache_write: 28173
          cost: 2.5831
        - model: claude-sonnet-5-5
          input: 16
          output: 111
          cache_read: 242075
          cache_write: 51213
          cost: 0.237
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-10-12T14:55:54Z
  value: 12.5
  by: planner-S-0365
  at: 2026-10-12T14:57:20Z
forecast:
  duration: 23m
  delivery: 2026-10-13T03:03:00Z
  basis: "flai forecast: 78 s per unit of size over 57 done large-band improvement stories on claude-opus-5-5, times size 17; 27th in the pull order with an in-progress limit of 5. Kept: S-0348, the same fix for board.wip-limit, took 24m."
  by: planner-S-0365
  at: 2026-10-12T14:57:20Z
finalized:
  by: orchestrator
  at: 2026-10-12T14:57:58Z
---
# S-0365 flai check finds `issues.no-story` outside the story at close-out

## Goal

This story remediates [I-0134](../../../design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md), "flai check finds `issues.no-story` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0134 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0134 is closed with `flai issue close I-0134 --reason` saying what fixed it

## Tasks
- T-1447 Record in an ADR that a close-out records no issues.no-story and a check scoped to a story leaves it out
- T-1448 Leave issues.no-story outside the story out of a check scoped to a story, with a test that reproduces I-0134
- T-1449 Say in the design, the user guide, and the work-management convention that a close-out leaves out issues.no-story
- T-1450 Close I-0134 saying what fixed it

## Notes

### Planning

Proposed fix, from I-0134's four instances: a check scoped to a story leaves `issues.no-story` outside the story out. It already leaves out `item.archive` (ADR-0122) and `board.wip-limit` (ADR-0133) the same way. The close-outs of S-0343, S-0289, S-0342, and S-0361 each found it on I-0052, which none of them touched. The finding names an issue, not a story. Only a story made for that issue clears it, and that is the operator's choice. The review page already lists every open issue no open story links. `flai check` in the main checkout still warns. An `issues.no-story` on an issue file the story changes stays inside.

Alternative, for T-1447 to weigh: a note without an issue, as ADR-0115 gives a `wip.overlap` naming the story. It keeps the nudge in `design/conventions/continuous-improvement.md` in front of each story's agent. But four close-outs saw it and none made the story.

Where the code is:

- The finding comes from `issuesNoStory` in `flai/internal/check/check.go`.
- The scoped check is `ScopeToStory` in `flai/internal/check/scope.go`.
- The issue is recorded by `recordOutside` in `flai/cmd/check.go`.

Tasks, in layers:

1. T-1447, the ADR.
2. T-1448, the code and its tests, and T-1449, the design, user guide, and convention. Both wait for T-1447 for its number, and they share no path.
3. T-1450 closes I-0134, after T-1448 and T-1449.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/check/scope.go`, `flai/cmd/check.go` | layout: where the finding is scoped and recorded |
| `flai/internal/check/scope_test.go`, `flai/cmd/check_test.go`, `docs/users/flai-reference.md`, `docs/users/flai.md`, `design/system/flai-cli.md`, `design/adrs/README.md` | co-change with the two above (`flai touches suggest`, 17-44%); S-0348 changed each |
| `design/system/continuous-improvement.md`, `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | design: they list what a close-out leaves out; S-0348 changed each |
| `design/adrs` | design: the new ADR, as in each precedent |
| `design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md`, `design/issues/summary.md` | the issue the story closes |

Left out:

- `flai/internal/check/check.go` (33% co-change): the fix scopes the finding and leaves its production alone.
- `flai/internal/check/check_test.go` and `flai/cmd/check_stats_test.go`: they test production and stats, which do not change.
- `design/conventions/continuous-improvement.md` and `flai/internal/harness/harness.go`: their rule to make a story when `flai check` warns `issues.no-story` stays true for the main checkout's check. If T-1447 decides otherwise, the story's agent adds them.

Folder touch kept: `design/adrs`. T-1447 adds an ADR whose number is not known until it is written, and other open stories, such as S-0363, add ADRs too.

Forecast: 23m, delivery 2026-10-13T03:03:00Z, as `flai forecast` gives it: 78 s per unit of size over 57 large-band improvement stories, times size 17. Kept: S-0348, the same fix for `board.wip-limit` with the same four tasks, was forecast at 21m and took 24m.

Cost of delay: 12.50 USD a week, as `flai cod` gives it from `time_lost_per_cycle: 5m`. The orchestrator set that input on TH-0404, as recommended. Kept: it matches S-0348 and S-0346, close-out findings of the same kind.

### Accepted by the orchestrator

- Verified: 2d9efc42214077bced555ea9a5e7b236aafcfa1f
- At: 2026-10-12T15:17:00Z

Verdict: pass at 2d9efc42214077bced555ea9a5e7b236aafcfa1f; flai verify passed every step at that commit, and the verifier matched both criteria to the diff.

- 1: flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/adrs/0136 (ADR), design/adrs/README.md, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md
- 2: design/issues/I-0134-flai-check-finds-issues-no-story-outside-the-story-at-close-out.md, design/issues/summary.md
