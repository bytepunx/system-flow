---
id: S-0245
type: story
nature: improvement
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
status: backlog
owner: alex
created: 2026-10-03T17:49:40Z
updated: 2026-10-06T22:56:13Z
transitions: []
tags: []
topics: [cli]
touches: [flai/internal/adr, flai/cmd/adr.go, flai/cmd/adr_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flai-cli.md, design/system/documentation-standard.md, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: planner-S-0245
    at: 2026-10-05T05:44:59Z
  value: 50
  by: planner-S-0245
  at: 2026-10-05T05:45:03Z
forecast:
  duration: 25m
  delivery: 2026-10-07T06:58:00Z
  basis: "Its own forecast of 25m; 21st in the pull order with an in-progress limit of 3, behind S-0299, S-0301, S-0300, S-0228, S-0261, S-0269, S-0270, S-0271, S-0212, S-0213, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239 and S-0241."
  by: flai
  at: 2026-10-06T22:56:13Z
---
# S-0245 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

## Goal

This story remediates [I-0063](../../../design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md), "flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0063 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0063 is closed with `flai issue close I-0063 --reason` saying what fixed it

## Tasks
- T-0915 flai adr new numbers past every ADR on main, every story worktree, and every story branch
- T-0916 The design and the user guides say how ADRs are numbered, and I-0063 is closed

## Notes

### Planning

Proposed fix, from I-0063's instances. `adr.NextNumber` (`flai/internal/adr/adr.go`) reads only the checkout's `design/adrs`. The fix makes it number one past the highest ADR in two places: the checkout's own files, and the names `storygit.FolderNames` returns for `design/adrs`. S-0252 added that helper for issue numbers (I-0065), and it already reads the main checkout, every linked worktree (uncommitted files included), the main branch, and every local `story/*` branch. `adr.New` already takes a runner, so no caller's signature changes. I-0063's issue-number instance was fixed by S-0252.

Touches, none declared before:

- `flai/internal/adr`: layout and design. `NextNumber` and `New` live there.
- `flai/cmd/adr.go`, `flai/cmd/adr_test.go`: layout and co-change. `flai touches suggest` found the test in all 3 commits that changed the package. These files hold the help text and the reproduction test.
- `docs/users/flai-reference.md`: layout. It is regenerated from the help text.
- `design/system/flai-cli.md`, `design/system/documentation-standard.md`, `docs/users/flai.md`, `docs/users/flaiover.md`: design. Each says the number is "one more than the highest file present".
- `design/issues/I-0063-…md`, `design/issues/summary.md`: the second criterion closes I-0063, and closing regenerates the summary.
- Left out: `flai/internal/docedit/docedit.go`, which `touches suggest` listed (2 of 3 commits). It guards edits to accepted ADRs, not numbering. `flai/internal/storygit` is also left out, because its helper is reused as it is.

Forecast: flai gave 18m (86 s per unit of size times size 12). I raised it to 25m. S-0252, the same fix for issues, took 31m. This story reuses that story's helper but still needs a git fixture of branches and worktrees, four documents, and a reference regeneration. The delivery, 2026-10-05T23:54Z, is flai's for 28th in the pull order, moved by the added 7m.

Cost of delay: 50.00 USD a week. `flai cod` works it out from the operator's input of `time_lost_per_cycle: 20m` (TH-0137): 20m per 168h cycle at 150 USD an hour. It stands unadjusted. The operator chose 20m over the recommended 10m, which allows for the links fixed by hand in S-0200's instances.

Topics: `cli` is added, since both tasks change flai or the CLI's design.
