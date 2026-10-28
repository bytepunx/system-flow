---
id: S-0525
type: story
nature: remediation
title: A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore
status: backlog
owner: alex
created: 2026-10-28T12:53:57Z
updated: 2026-10-28T21:36:41Z
transitions: []
tags: [flai, check, cli]
topics: [cli, testing]
touches: [flai/internal/itemnew/itemnew.go, flai/internal/itemnew/itemnew_test.go, flai/internal/itemedit/itemedit.go, flai/internal/check/scope.go, design/adrs, design/adrs/README.md, design/system/flai-cli.md, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0191-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md, design/issues/summary.md]
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
      seconds: 157
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 34
          output: 582
          cache_read: 12537273
          cache_write: 13570
          cost: 2.7553
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-10-28T12:53:57Z
  value: 37.5
  by: planner-S-0525
  at: 2026-10-28T14:53:26Z
forecast:
  duration: 45m
  delivery: 2026-10-29T11:40:00Z
  basis: "Its own forecast of 45m; 34th in the pull order with an in-progress limit of 5, behind S-0416, S-0419, S-0423, S-0418, S-0421, S-0422, S-0420, S-0417, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510, S-0511, S-0512, S-0514, S-0515, S-0516, S-0517, S-0518, S-0519, S-0520, S-0523 and S-0524."
  by: flai
  at: 2026-10-28T21:36:41Z
finalized:
  by: alex
  at: 2026-10-28T12:58:59Z
---
# S-0525 A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore

## Goal

This story remediates [I-0191](../../../design/issues/I-0191-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md), "A work item write in a project with no .gitignore makes .flai-cache/items.lock and then refuses itself on layout.gitignore". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0191 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0191 is closed with `flai issue close I-0191 --reason` saying what fixed it

## Tasks
- T-2251 A checked creation refuses only a new finding inside the new item's story, with a test that reproduces I-0191
- T-2252 An ADR records that a checked creation refuses only a new finding inside the new item's story, refining ADR-0197
- T-2253 The design and the user documentation say a checked creation is scoped to the new item's story as a checked edit is
- T-2254 I-0191 is closed with the reason that names the fix

## Notes

Cost of delay inputs set by flai from I-0191. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-22T20:40:04Z, 5.7 days before this story; under one cycle counts as one).

### Planning

Reproduced on 2026-10-28 on a temp copy of `flai/internal/metrics/testdata/good` (no `.gitignore`, no `.flai-cache`): `flai edit` now warns and stands, since ADR-0197 scoped the checked edit; `flai epic new --body-stdin` as the first write refuses itself with `layout.gitignore`, and the same creation run again stands, since the second run's before-check already knows the finding. The live cause is `itemnew.Create`, which still diffs two unscoped runs (ADR-0197 point 6), reached by `story new`, `task new`, and `epic new` with `--body-stdin`, `--after`, or `--autocommit`, by `item_new` with a body or `after`, and by `issue story`.

Remedy recommended, from the three the issue lists: scope the checked creation as ADR-0197 scopes the checked edit, so a new finding outside the new item's story, such as `layout.gitignore` on the repository root, is a warning and the creation stands. The lock stays where ADR-0184 (accepted) puts it, and flai writes no tracked `.gitignore`. The plan's thread records the alternatives.

Touches, declared: none. Layout: `flai/internal/itemnew/itemnew.go` and its test (the refusing writer and the reproduction), `flai/internal/itemedit/itemedit.go` and `flai/internal/check/scope.go` (where `scope` lives today and where it may move so both writers share it), `design/adrs/README.md` and the folder `design/adrs` (a new ADR refining ADR-0197, whose number and slug no task can name until it is written: `flai adr new` numbers it), `design/system/flai-cli.md`, `design/system/strategic-agents.md` § Checked creation, and `docs/users/flai.md` (where checked creation is described), `design/issues/I-0191-…md` and `design/issues/summary.md` (`flai issue close` rewrites both). Co-change: `flai touches suggest` refused the story with no touches to start from; not rerun, since the layout names the files.

Forecast: `flai forecast` gave 22m (median 99 s per unit over 39 done remediation stories on claude-opus-5-5 in the large band, size 13: 2 criteria and 11 touches) and delivery 2026-10-29T05:16Z at 35th in the pull order. Raised to 45m, delivery 2026-10-29T05:39Z, because the size counts touches, not the four tasks, and the ADR and the three documentation files are each a write of their own beside the fix. Cost of delay: 37.50 USD a week from the inputs (15m lost per 168h cycle at 150 USD an hour), as `flai cod` worked it out; it stands. No `after`: the fix needs no other story's change on main.
