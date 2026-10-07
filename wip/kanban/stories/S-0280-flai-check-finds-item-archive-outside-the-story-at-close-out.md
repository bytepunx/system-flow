---
id: S-0280
type: story
nature: improvement
title: "flai check finds `item.archive` outside the story at close-out"
status: backlog
owner: alex
created: 2026-10-05T07:09:05Z
updated: 2026-10-07T21:42:43Z
transitions: []
tags: [flai, template]
topics: [cli, conventions, template]
touches: [design/adrs, flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, docs/users/flai-reference.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md, design/issues/summary.md]
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
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 10
          cache_read: 40634
          cache_write: 303
          cost: 0.0105
draft: true
forecast:
  duration: 25m
  delivery: 2026-10-08T05:54:00Z
  basis: "Its own forecast of 25m; 10th in the pull order with an in-progress limit of 3, behind S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239 and S-0241."
  by: flai
  at: 2026-10-07T21:42:43Z
---
# S-0280 flai check finds `item.archive` outside the story at close-out

## Goal

This story remediates [I-0078](../../../design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md), "flai check finds `item.archive` outside the story at close-out". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0078 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0078 is closed with `flai issue close I-0078 --reason` saying what fixed it

## Tasks
- T-1164 An ADR refining ADR-0085 and ADR-0115 records the remedy for I-0078, proposed from its instances
- T-1165 A close-out records no item.archive in an issue, and its scoped check leaves out one that does not name the story
- T-1166 I-0078 is closed with flai issue close, saying that a close-out records no item.archive
- T-1167 The convention, the design, and the users' guide say that a close-out records no item.archive

## Notes

### Planning

Planned by planner-S-0280 on 2026-10-07. The plan thread on S-0280 lists the tasks, their layers, and the assumptions.

The proposed remedy comes from I-0078's 38 instances. T-1164 records it in an ADR before any code is built.

- **One cause.** All 38 instances name S-0250. The operator cancelled it from backlog on 2026-10-05, and it stayed in `wip/kanban/stories` until commit 1f7b7296 archived it by hand on 2026-10-07. No instance came after that.
- **Never the story's.** `item.archive` warns on a done or cancelled epic or story that is not archived. The closing story is in progress, so the finding never names it. No story branch makes it, and no story's agent can fix it in its worktree.
- **Recorded every time.** `ScopeToStory` marks it outside, and `recordOutside` records every outside rule but `wip.overlap`.
- **Remedy.** As ADR-0115 did for `wip.overlap`: a close-out records no `item.archive` and leaves out one that does not name the story. `flai check` in the main checkout still warns.
- **Not taken.** Archiving on cancel. It would reverse ADR-0055's move back from cancelled and ADR-0028's narrative left in place.

Layers:

1. T-1164, the ADR.
2. T-1165, the scoped check and the recording, with tests that reproduce I-0078.
3. T-1166, closing I-0078; T-1167, the convention, the design, and the guides. They share no path and run together.

Touches:

- **Declared:** none. The story declared no touches.
- **Layout:**
  - `flai/internal/check/scope.go` and `scope_test.go`: `ScopeToStory`.
  - `flai/cmd/check.go` and `check_test.go`: `recordOutside` and the help.
- **Co-change:** `docs/users/flai-reference.md`, `design/system/flai-cli.md`, and `docs/users/flai.md` change with `flai/cmd/check.go` in 31 to 38% of its commits. The reference is generated from the changed help.
- **Design:**
  - The two copies of `work-management.md`: the close-out rule names `wip.overlap`.
  - `template/CHANGELOG.md`: the template convention changes.
  - `design/system/continuous-improvement.md`: it describes what a close-out records (ADR-0085).
- **Criterion 2:** the I-0078 file and `design/issues/summary.md`, which `flai issue close` writes.
- **Folder touch kept:** `design/adrs`. T-1164 adds an ADR whose number and slug `flai adr new` picks, so no task can name the file yet. It lies inside `claims.shared`, so it holds no story.
- **Not taken from `flai touches suggest`:** `flai/internal/check/check.go` and `check_test.go`, since the rule itself does not change. Also `design/system/workflow.md`, `docs/operators/settings.md`, and the other co-changes the goal names none of. T-1165 widens its touches if `check.go` must change.

Forecast 25m, delivery 2026-10-07T21:27Z.

- `flai forecast` gave 21m: 78 s per unit of size over 43 done large-band improvement stories on claude-opus-5-5, times size 16 (2 criteria, 14 touches). It delivered at 21:23Z, 10th in the pull order.
- Raised by 4m. The story has three serial layers, each with its own commit, sync, and test cycle. S-0279, the same remedy for I-0076 with one more layer, took 26m of agent time.
- Delivery is shifted by the same 4m.

Cost of delay: not set yet.

- `flai cod` cannot work it out: the story has no inputs and no epic.
- The inputs are the operator's, asked for on the cost thread on S-0280. The recommendation there is `time_lost_per_cycle: 45m`, as S-0279 took.
