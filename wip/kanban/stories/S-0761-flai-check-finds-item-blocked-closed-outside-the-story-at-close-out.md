---
id: S-0761
type: story
nature: improvement
title: "flai check finds `item.blocked-closed` outside the story at close-out"
status: backlog
owner: alex
created: 2026-11-13T07:23:41Z
updated: 2026-11-13T07:33:44Z
transitions: []
tags: [flai, check]
topics: [cli, conventions, template]
touches: [wip/archive/kanban/stories/S-0678-agent-status-with-13-projects-at-once-is-not-under-200-ms.md, design/adrs, design/adrs/README.md, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md, design/issues/summary.md]
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
      seconds: 422
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 98
          output: 28758
          cache_read: 5519178
          cache_write: 157511
          cost: 2.9395
    - kind: orchestrator
      seconds: 137
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 39
          output: 284
          cache_read: 1851475
          cache_write: 204187
          cost: 0.5901
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1m
    by: orchestrator
    at: 2026-11-13T07:32:21Z
  value: 2.5
  by: planner-S-0761
  at: 2026-11-13T07:33:44Z
forecast:
  duration: 33m
  delivery: 2026-11-13T21:23:00Z
  basis: "Its own forecast of 33m; 33rd in the pull order with an in-progress limit of 3, behind S-0660, S-0733, S-0760, S-0735, S-0745, S-0759, S-0734, S-0751, S-0757, S-0740, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750, S-0752, S-0753, S-0754, S-0755, S-0756 and S-0758."
  by: flai
  at: 2026-11-13T07:32:52Z
finalized:
  by: alex
  at: 2026-11-13T07:27:32Z
---
# S-0761 flai check finds `item.blocked-closed` outside the story at close-out

## Goal

This story remediates [I-0354](../../../design/issues/I-0354-flai-check-finds-item-blocked-closed-outside-the-story-at-close-out.md), "flai check finds `item.blocked-closed` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0354 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0354 is closed with `flai issue close I-0354 --reason` saying what fixed it

## Tasks
- T-3268 Close S-0678's blocked interval at its close with the tree's flai unblock
- T-3269 Record in an ADR whether a close-out leaves out item.blocked-closed outside the story
- T-3270 Leave item.blocked-closed outside the story out of a check scoped to a story, with a test that reproduces I-0354
- T-3271 Say in the design, the user guide, and the work-management convention that a close-out leaves out item.blocked-closed
- T-3272 Close I-0354 saying what fixed it

## Notes

### Planning

Cause, from I-0354's one instance and I-0347's 36: every one names the archived S-0678. It was blocked at 2026-11-11T05:26:16Z ("pause for urgent fix") and accepted by the orchestrator at 06:15:07Z with the interval still open. S-0723 fixed how that happens: a move to `done` or `cancelled` now closes an open interval, and `flai unblock` on a closed item closes one left open at the item's close. But S-0723 could not repair S-0678 itself, which lives in `wip/archive` in the main checkout, and it closed I-0347 on its branch. Its own close-out then met S-0678's finding again and recorded it as I-0354. S-0678 still has an open `blocked` interval, so every close-out still notes it.

Proposed fix, in two parts:

1. Repair S-0678 once, with `scripts/flai.sh unblock S-0678` in the main checkout. The tree's flai has S-0723's fix and closes the interval at S-0678's close; the installed 1.86.17 would close it at now.
2. Leave `item.blocked-closed` outside the story out of a check scoped to a story, as S-0366 (ADR-0137) did for `threads.archived` and S-0555 (ADR-0256) for `narrative.missing`. No story branch can clear the finding on another item: only `flai unblock` on that item, in the main checkout, does. The ADR task weighs the alternative, repair alone, since S-0723 already stops new instances from flai's own moves; a hand edit of front matter is the one way left to make one.

Tasks and layers:

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-3268 Repair S-0678 | nothing |
| 1 | T-3269 The ADR | nothing |
| 2 | T-3270 The scoped check leaves the finding out, with tests | T-3269, for its number and decision |
| 2 | T-3271 The design, guide, convention, and changelog say so | T-3269, for its number and decision |
| 3 | T-3272 Close I-0354 | T-3268, T-3270, T-3271: the reason names the repair and the ADR |

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `wip/archive/kanban/stories/S-0678-…md` | Goal: the repair writes S-0678's blocked interval, in the main checkout |
| `flai/internal/check/scope.go`, `flai/internal/check/scope_test.go` | Co-change (89% for the test) and layout: `ScopeToStories` drops the other left-out rules; `scope_test.go` holds the `TestScopedCheckLeavesOut…` tests |
| `flai/cmd/check.go`, `flai/cmd/check_test.go` | Co-change (67%, 61%) and layout: `recordOutside`, the `--record-issues` help, and the `TestCheckRecordIssuesLeavesOut…` tests |
| `docs/users/flai-reference.md` | Co-change (56%): `make flai-reference` regenerates it from `flai check`'s help |
| `design/adrs`, `design/adrs/README.md` | Design: S-0366 and S-0555 each needed an ADR for the same kind of fix |
| `design/system/flai-cli.md`, `docs/users/flai.md`, `design/system/continuous-improvement.md` | Design: each lists what a scoped check leaves out, as S-0555 changed them |
| `design/conventions/work-management.md`, `template/root/design/conventions/work-management.md`, `template/CHANGELOG.md` | Design: the close-out rule lists the left-out findings; its baseline changes in `template/` in the same story |
| `design/issues/I-0354-…md`, `design/issues/summary.md` | Criterion 2: `flai issue close I-0354` |

The folder touch `design/adrs` stays: the ADR's file name comes from the title its task decides, so no task can name it yet. `flai/internal/check/check.go` and `check_test.go` (co-change 33%, 11%) are left out: the rule's text in `check_items.go` already names `flai unblock`, and the rule itself does not change. `design/system/workflow.md` lists no scoped finding, so it is left out too.

No `after`: the story reads only S-0723's output, and S-0723 is done.

Forecast 33m stands as `flai forecast` gives it, at size 18 from 2 criteria and 16 touches. S-0555, the same fix for `narrative.missing` with four tasks, took 40m. S-0366, for `threads.archived`, took 17m. This story adds a one-command repair task, so a figure between the two fits.

Cost of delay 2.50 USD a week stands as `flai cod` gives it. The input, `time_lost_per_cycle: 1m`, is the one recommended on TH-1137 and set by the orchestrator, as TH-1096 set it for S-0723 on the same finding. At 150 USD an hour and one 168h cycle a week, that is 2.50 USD.
