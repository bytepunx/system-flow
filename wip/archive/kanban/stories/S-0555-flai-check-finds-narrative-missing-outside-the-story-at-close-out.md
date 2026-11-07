---
id: S-0555
type: story
nature: improvement
title: "flai check finds `narrative.missing` outside the story at close-out"
status: done
owner: alex
created: 2026-11-01T05:45:09Z
updated: 2026-11-07T04:39:00Z
transitions:
  - to: ready
    at: 2026-11-07T03:56:08Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T03:56:17Z
    by: agent-S-0555
  - to: review
    at: 2026-11-07T04:36:21Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:39:00Z
    by: orchestrator
tags: [flai, check]
topics: [testing, cli]
touches: [design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, docs/users/flai.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md, design/issues/summary.md, design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2417
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      hand_edits: 1
      work: 77
  models:
    - model: claude-opus-5-5
      input: 238
      output: 39341
      cache_read: 30075828
      cache_write: 619483
      cost: 10.7978
  strategic:
    - kind: orchestrator
      seconds: 146
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 28
          output: 190
          cache_read: 1157905
          cache_write: 217439
          cost: 0.3337
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: orchestrator
    at: 2026-11-07T03:45:23Z
  value: 12.5
  by: planner-S-0555
  at: 2026-11-07T03:46:50Z
forecast:
  duration: 26m
  delivery: 2026-11-07T08:53:00Z
  basis: "flai forecast: 92 s per unit of size over 120 done large-band improvement stories on claude-opus-5-5, times size 17; 10th in the pull order with an in-progress limit of 5. Kept: S-0365 and S-0366, the same fix for two other rules, took 18m and 17m, and close-outs now run the integration and smoke tiers too."
  by: planner-S-0555
  at: 2026-11-07T03:48:40Z
finalized:
  by: orchestrator
  at: 2026-11-07T03:54:02Z
---
# S-0555 flai check finds `narrative.missing` outside the story at close-out

## Goal

This story remediates [I-0235](../../../design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md), "flai check finds `narrative.missing` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0235 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0235 is closed with `flai issue close I-0235 --reason` saying what fixed it

## Tasks
- T-2518 Record in an ADR that a close-out records no narrative.missing and a check scoped to a story leaves it out
- T-2521 Leave narrative.missing outside the story out of a check scoped to a story, with a test that reproduces I-0235
- T-2527 Say in the design, the user guide, and the work-management convention that a close-out leaves out narrative.missing
- T-2530 Close I-0235 saying what fixed it

## Notes

### Planning

The proposed fix follows S-0366 (ADR-0137, `threads.archived`) and ADR-0125 (`narrative.state` on another open story's narrative). Every instance of I-0235 is the same finding: S-0549 was in progress with no narrative, and five other stories' close-outs (S-0546, S-0547, S-0548, S-0550, S-0551) noted and recorded it. The operator moved S-0549 to in progress at 2026-10-31T09:47:53Z. `flai move` opens no narrative: only `flai story start` and `flai stream open` do (`flai/cmd/move.go`, `OpenStream` in `flai/internal/workitem/narrative.go`). No other story's branch can open S-0549's narrative. Only S-0549's agent or the operator can. So a check scoped to a story leaves out `narrative.missing` outside it, and the close-out records none. The ADR task weighs the alternative, a move to in progress that opens the narrative, and says why it is or is not chosen.

Tasks and layers: T-2518 (the ADR) first. Then T-2521 (code and tests) and T-2527 (design, user guide, convention) together, since they share no path. Then T-2530 (close I-0235).

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/check/scope.go`, `flai/internal/check/scope_test.go` | Layout: `ScopeToStories` drops the other rules outside the story, and `scope_test.go` holds the `TestScopedCheckLeavesOut…` tests. Co-change 14%. |
| `flai/cmd/check.go`, `flai/cmd/check_test.go` | Layout: `recordOutside` and the `--record-issues` help, and the `TestCheckRecordIssuesLeavesOut…` tests. Co-change 14%. |
| `docs/users/flai-reference.md` | Co-change 22%. `make flai-reference` generates it from the help of `flai check`. |
| `design/adrs`, `design/adrs/README.md` | Design: the fix needs an ADR, as S-0366 needed ADR-0137. Co-change 10%. |
| `design/system/flai-cli.md`, `docs/users/flai.md`, `design/system/continuous-improvement.md` | Design and co-change (50%, 51%, 11%): each lists what a scoped check leaves out. |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | Design and co-change (59%, 45%): the close-out rule lists the left-out findings. Its baseline changes in `template/` in the same story. |
| `design/issues/I-0235-…md`, `design/issues/summary.md` | Criterion 2: `flai issue close I-0235`. |

The folder touch `design/adrs` stays. The ADR's file name comes from its title, which the ADR task decides, so no task can name it yet. `design/system/workflow.md` (co-change 29%) is left out because it lists no scoped finding.

No `after`: the story reads no output of another open story.

Forecast: 26m, as `flai forecast` gives it at size 17. Kept: S-0365 and S-0366, the same fix for `issues.no-story` and `threads.archived`, took 18m and 17m of agent time. Close-outs now also run the integration and smoke tiers for a story that touches `flai/` or `template/`.

Cost of delay: 12.50 USD a week, as `flai cod` works it out from the operator's `time_lost_per_cycle: 5m` (TH-0837), at 150 USD an hour and one cycle a week. Kept as it is, since S-0556, a defect of the same kind, has the same figure.

### Accepted by the orchestrator

- Verified: 7bef3eda279499d5698a0237cb1e4ccbe9c17f04
- At: 2026-11-07T04:39:00Z

Verdict: pass — a check scoped to a story now leaves out narrative.missing outside the story and the close-out records no issue for it, tests reproduce the S-0549 case, and I-0235 is closed with a reason.
- 1: flai/internal/check/scope.go, flai/cmd/check.go, flai/internal/check/scope_test.go, flai/cmd/check_test.go
- 2: design/issues/I-0235-flai-check-finds-narrative-missing-outside-the-story-at-close-out.md, design/issues/summary.md
