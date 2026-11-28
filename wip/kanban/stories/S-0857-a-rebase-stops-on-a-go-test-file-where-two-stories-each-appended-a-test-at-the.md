---
id: S-0857
type: story
nature: improvement
title: A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back
status: backlog
owner: alex
created: 2026-11-27T22:55:28Z
updated: 2026-11-28T00:26:08Z
transitions: []
tags: []
topics: [git, code]
touches: [design/adrs, design/adrs/README.md, flai/internal/storygit/gomerge.go, flai/internal/storygit/gomerge_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/storygit/trial_test.go, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0409-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md, design/issues/summary.md]
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
      seconds: 13
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 11
          cache_read: 309097
          cache_write: 6501
          cost: 0.1005
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-27T22:55:28Z
  value: 37.5
  by: planner-S-0857
  at: 2026-11-27T23:58:38Z
forecast:
  duration: 1h30m
  delivery: 2026-11-28T11:01:00Z
  basis: "Its own forecast of 1h30m; 16th in the pull order with an in-progress limit of 3, behind S-0835, S-0837, S-0844, S-0814, S-0852, S-0813, S-0822, S-0848, S-0853, S-0854, S-0824, S-0845, S-0241, S-0815, S-0849, S-0851, S-0855 and S-0856."
  by: flai
  at: 2026-11-28T00:26:08Z
finalized:
  by: alex
  at: 2026-11-28T00:01:57Z
---
# S-0857 A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back

## Goal

This story remediates [I-0409](../../../design/issues/I-0409-a-rebase-stops-on-a-go-test-file-where-two-stories-each-appended-a-test-at-the.md), "A rebase stops on a Go test file where two stories each appended a test at the same place, which flai cannot settle, so an acceptance sends the story back". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0409 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0409 is closed with `flai issue close I-0409 --reason` saying what fixed it

## Tasks
- T-3623 ADR: a rebase settles a stop on a Go file where each side only adds whole top-level declarations at one place
- T-3625 storygit merges a Go file's three stages when each side only adds whole top-level declarations at one place
- T-3628 Every rebase flai runs settles a stop on a Go file with the Go merge, and the trial merge leaves such a pair out
- T-3630 The rebase and resolve job keeps a story in review whose rebase stopped only on two appends to one Go test file
- T-3633 workflow.md and the user guide say a rebase settles two stories' whole-declaration appends to one Go file
- T-3634 Close I-0409 naming the Go merge that settles it

## Notes

Cost of delay inputs set by flai from I-0409. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-21T18:16:18Z, 6.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the one instance (S-0809 and S-0812 in `flai/internal/verify/select_test.go`, MS-0628, TH-1320): every rebase flai runs settles a stop on a `.go` file when each side only adds whole top-level declarations at the same place, keeping ours then theirs. It works from the three stages, as the issue-file merge does, not from git's markers, since git's hunk split the closing brace between the two sides. It parses base, ours, theirs, and the result with `go/parser` and refuses, leaving the stop to wait as today, when a side changes anything else, when the two add a declaration of the same name, or when the result does not parse or is not gofmt-clean. It covers every `.go` file, not only `_test.go`, since the rule is the same. A change to the import block is not settled.

Where the touches came from. `flai touches suggest` listed nothing, since the story declared no touches.

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs` | design | The new ADR that refines ADR-0175 takes the next free number when it is written. Other open stories take numbers meanwhile, so no task can name the file yet. Kept as a folder for that file alone. |
| `design/adrs/README.md` | design | The ADR index gains the new ADR. |
| `flai/internal/storygit/gomerge.go`, `gomerge_test.go` | layout | The new merge of a Go file's three stages, beside `mdmerge.go`, whose `merge3` it reuses. |
| `flai/internal/storygit/sync.go`, `sync_test.go` | layout | `Files`, `resolves`, `Unresolved`, and the merge at the stop gain the Go kind. |
| `flai/internal/storygit/trial_test.go` | layout | The trial merge leaves a settled Go pair out through `Files.Unresolved`. |
| `flai/internal/issues/generated.go`, `generated_test.go` | layout | `FilesIn` builds the set every rebase shares: stream sync, task done, accept, the merge queue's batch, and the trial merge. |
| `flai/internal/serve/resolvejob.go`, `resolvejob_test.go` | layout | `conflictText` names the kinds settled; the test reproduces I-0409 through the acceptance's rebase and resolve job. |
| `design/system/workflow.md` | design | The section on what a rebase settles (S-0401, ADR-0175, ADR-0276). |
| `docs/users/flai.md` | design | Shared documents and the rebase and resolve job's list of what it settles. |
| `design/issues/I-0409-…`, `design/issues/summary.md` | declared | Closing I-0409, criterion 2. |

Figures.

- Forecast: flai gave 2m from size 2 with no touches, then 31m once the touches were written (size 17, 109 s a unit over 187 large improvement stories). Raised to 1h30m, delivery 2026-11-28T08:20:00Z: the Go merge is new code on `go/ast` with real-git tests in three packages, and the close-out of a story touching `flai/` runs the integration and smoke tiers, about 12m alone on S-0847.
- Cost of delay: 37.50 USD a week stands, from flai's inputs (15m lost per 168h cycle at 150 USD an hour). One occurrence so far; D-0644 shows the same conflict worked around by hand by moving a test to its own file.

No `after`: the story reads nothing another open story writes.
