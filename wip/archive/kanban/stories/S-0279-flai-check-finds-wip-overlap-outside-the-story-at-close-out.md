---
id: S-0279
type: story
nature: improvement
title: "flai check finds `wip.overlap` outside the story at close-out"
status: done
owner: alex
created: 2026-10-05T04:40:48Z
updated: 2026-10-07T14:35:49Z
transitions:
  - to: ready
    at: 2026-10-07T09:23:05Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T09:30:05Z
    by: agent-S-0279
  - to: review
    at: 2026-10-07T09:52:14Z
    by: agent-S-0279
  - to: done
    at: 2026-10-07T14:35:49Z
    by: alex
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/summary.md, design/system/work-hierarchy.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1555
  models:
    - model: claude-opus-5-5
      input: 290
      output: 81659
      cache_read: 12658618
      cache_write: 499048
      cost: 7.4418
  strategic:
    - kind: planner
      seconds: 574
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 438
          output: 678
          cache_read: 2741384
          cache_write: 131610
          cost: 0.5747
        - model: claude-opus-5-5
          input: 190
          output: 28578
          cache_read: 10046553
          cache_write: 284594
          cost: 5.2617
    - kind: orchestrator
      seconds: 1385
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 104
          output: 1553
          cache_read: 11686783
          cache_write: 88021
          cost: 3.0723
        - model: claude-sonnet-5-5
          input: 8
          output: 48
          cache_read: 72235
          cache_write: 38139
          cost: 0.074
cost_of_delay:
  inputs:
    time_lost_per_cycle: 45m
    by: planner-S-0279
    at: 2026-10-07T01:21:44Z
  value: 112.5
  by: planner-S-0279
  at: 2026-10-07T01:22:05Z
forecast:
  duration: 45m
  delivery: 2026-10-07T13:19:00Z
  basis: "Its own forecast of 45m; 11th in the pull order with an in-progress limit of 3, behind S-0215, S-0246, S-0275, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239 and S-0241."
  by: flai
  at: 2026-10-07T08:59:42Z
finalized:
  by: orchestrator
  at: 2026-10-07T08:20:27Z
---
# S-0279 flai check finds `wip.overlap` outside the story at close-out

## Goal

This story remediates [I-0076](../../../design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md), "flai check finds `wip.overlap` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0076 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0076 is closed with `flai issue close I-0076 --reason` saying what fixed it

## Tasks
- T-1141 An ADR refining ADR-0085 and ADR-0096 records the remedy for I-0076, proposed from its instances
- T-1143 wip.overlap compares the claims of the stories in progress, one finding per pair of stories, as the pull hold does
- T-1145 A close-out records no wip.overlap in an issue, and its scoped check lists as notes only the overlaps that name the story
- T-1149 I-0076 is closed with flai issue close, saying that a close-out records no wip.overlap and that wip.overlap compares claims
- T-1150 The convention, the design, and the users' guide say that wip.overlap compares claims and that a close-out records no overlap

## Notes

### Planning

Planned by planner-S-0279 on 2026-10-07. The plan thread on S-0279 lists the tasks, their layers, and the assumptions.

The proposed remedy comes from I-0076's 12 instances. T-1141 records it in an ADR before any code is built.

- **By design today.** ADR-0085 marks every `wip.overlap` as outside the story. So every close-out records every overlap on the board, even one that does not name the closing story, such as S-0251 recording S-0228 against S-0273.
- **Duplicates.** `overlap()` in `flai/internal/check/check.go` also pairs a story with another story's task, such as S-0262 against T-0877.
- **Raw touches.** It reads raw touches, not the claim the pull hold reads. So a folder touch is reported where its tasks narrow it (ADR-0096 §3), such as `flai/internal/serve` and `flai/internal/hostapi`.

Layers:

1. T-1141, the ADR.
2. T-1143, `wip.overlap` over claims, carrying the two story IDs on each finding.
3. T-1145, the close-out records no overlap and lists only those naming the story. It reads T-1143's IDs.
4. T-1149, closing I-0076; T-1150, the convention, the design, and the guides.

Touches:

- **Declared:** none. The story declared no touches.
- **Layout:**
  - `flai/internal/check/check.go` and `check_test.go`: `overlap()` and its tests.
  - `flai/internal/check/scope.go` and `scope_test.go`: `ScopeToStory`.
  - `flai/cmd/check.go` and `check_test.go`: `recordOutside` and the help.
- **Design:**
  - The two copies of `work-management.md`: the close-out rule names `wip.overlap`.
  - `template/CHANGELOG.md`: the template convention changes.
  - The `wip.overlap` mentions in `design/system/flai-cli.md`, `design/system/workflow.md`, `design/system/agent-coordination.md`, `design/system/continuous-improvement.md`, and `docs/users/flai.md`.
- **Co-change and design:** `docs/users/flai-reference.md`, generated from the changed help.
- **Criterion 2:** the I-0076 file and `design/issues/summary.md`, which `flai issue close` writes.
- **Folder touch kept:** `design/adrs`. T-1141 adds an ADR whose number and slug `flai adr new` picks, so no task can name the file yet. It lies inside `claims.shared`, so it holds no story.
- **Not taken from `flai touches suggest`:** the dashboard design, the operators' and flaiover guides, `template/template.yaml`, and the MCP, harness, and hostapi files. They are what any flai story co-changes with, and the goal names none of them. `flai/internal/hostapi/people.go` and `flai/internal/itemedit/itemedit.go` read `wip.overlap` findings. T-1143 runs their tests and widens its touches only if one must change.
- **Overlap with S-0277**, which remediates I-0073 the same way and is ready: it is likely to change `flai/internal/check/scope.go`, `flai/cmd/check.go`, and the close-out rule in `work-management.md`. Whichever story goes second rebases onto the other. The pull hold keeps them apart while both are open.

Forecast 45m, delivery 2026-10-07T10:26Z.

- `flai forecast` gave 26m: 78 s per unit of size over 32 done large-band improvement stories on claude-opus-5-5, times size 20 (2 criteria, 18 touches). It delivered at 10:07Z, 26th in the pull order.
- Raised by 19m. The story has four serial layers, each with its own commit, sync, and test cycle. It writes an ADR from the instances. Its claim-based `wip.overlap` is read by the designer's inbox and by itemedit too, so their tests must pass. Size counts none of this.
- Delivery is shifted by the same 19m.

Cost of delay 112.50 USD a week, as `flai cod` gives it.

- The input is `time_lost_per_cycle: 45m`, the operator's choice on the plan's cost thread: "accept alternative".
- I-0076 has 12 instances between 2026-10-05T03:24Z and 2026-10-07T01:04Z, about 44 a week. At about 1m of agent time and one issue-bump commit each, that is 45m per 168h cycle, at 150 USD an hour.
- The figure stands as computed.
