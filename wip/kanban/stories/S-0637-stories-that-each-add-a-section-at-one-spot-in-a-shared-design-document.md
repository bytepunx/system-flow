---
id: S-0637
type: story
nature: improvement
title: Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back
status: ready
owner: alex
created: 2026-11-09T08:06:42Z
updated: 2026-11-10T07:44:15Z
transitions:
  - to: ready
    at: 2026-11-10T07:44:15Z
    by: alex
tags: [flai, sync, acceptance]
topics: [cli, go, code]
touches: [flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go, design/adrs/0273-a-rebase-merges-by-sections-every-markdown-document-under-the-design-and-docs.md, design/adrs/README.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0315-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md, design/issues/summary.md]
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
      seconds: 23
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 8
          cache_read: 83530
          cache_write: 3604
          cost: 0.0233
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h10m
    by: flai
    at: 2026-11-09T08:06:42Z
  value: 175
  by: planner-S-0637
  at: 2026-11-09T22:08:49Z
forecast:
  duration: 1h30m
  delivery: 2026-11-10T10:01:00Z
  basis: "Its own forecast of 1h30m; 8th in the pull order with an in-progress limit of 5, behind S-0628, S-0636, S-0638, S-0647, S-0649, S-0650, S-0655, S-0651, S-0653, S-0654, S-0657 and S-0241."
  by: flai
  at: 2026-11-10T07:12:18Z
finalized:
  by: alex
  at: 2026-11-10T07:44:10Z
---
# S-0637 Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back

## Goal

This story remediates [I-0315](../../../design/issues/I-0315-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md), "Stories that each add a section at one spot in a shared design document conflict at every sync and acceptance, and the acceptance sends them back". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0315 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0315 is closed with `flai issue close I-0315 --reason` saying what fixed it

## Tasks
- T-2811 Record in an ADR that a rebase merges by sections every markdown document under the design and docs folders
- T-2812 Merge by sections a stopped markdown file under the design or docs folder, not only one in claims.shared
- T-2813 Say in the design and user docs that a stopped design or docs markdown file merges by sections
- T-2814 Reproduce I-0315 in real git: two stories add a section at one spot of a design document outside claims.shared, and neither sync nor acceptance stops
- T-2815 Close I-0315 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0315. time_lost_per_cycle 1h10m: 35m per occurrence × 2 occurrences ÷ 1 cycle of 168h (first reported 2026-11-09T07:27:44Z, 0 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0315's three instances: every resolution was mechanical, "keep both sections", and ADR-0175's section merge already settles two sections added at one spot ("one both added appears once per side, ours first"). It runs only on markdown inside `claims.shared`, and `design/system/server-performance.md` is not shared. The gate is `sharedMarkdown` in `flai/internal/issues/generated.go`, which `issues.FilesIn` gives every rebase flai runs: `flai stream sync`, `flai task done`'s sync, `flai accept`'s rebase (`flai/cmd/branch.go`), and the merge queue's batch (`flai/internal/mergequeue/batch.go`). Widening that one matcher to every markdown document under the manifest's design and docs folders, the issue files and the generated indexes left to their own steps, settles the stop in all four. A section merge that cannot settle still leaves the stop as it is, so a real conflict is still the agent's. The plan's thread asks the operator to confirm this before the story is built.

Touches, and where each came from (`flai touches suggest` listed nothing: the story declared no touches):

| Touch | Source | Why |
|---|---|---|
| `flai/internal/issues/generated.go` | layout | `sharedMarkdown` and `FilesIn`, the gate |
| `flai/internal/issues/generated_test.go` | layout | `TestFilesInMergesTheSharedMarkdownFilesBySections` names a non-shared design file as not merged |
| `flai/cmd/accept_conflict_test.go` | layout | the send-back tests; the reproduction of the acceptance |
| `flai/cmd/stream_sync_test.go` | layout | the reproduction of the sync |
| `design/adrs/0273-a-rebase-merges-by-sections-every-markdown-document-under-the-design-and-docs.md` | design | the ADR refining ADR-0175; `flai adr new` gives the number and slug, so the agent renames this touch when it differs |
| `design/adrs/README.md` | design | the generated ADR index |
| `design/system/workflow.md` | design | § Branches and collisions names the section merge's scope |
| `design/system/agent-coordination.md` | design | its ADR-0175 paragraph names the scope |
| `design/system/flai-cli.md` | design | the `flai stream sync` row and the merge queue's build name the merge steps |
| `docs/users/flai.md` | design | § Shared documents tells users which files merge by sections |
| `design/issues/I-0315-stories-that-each-add-a-section-at-one-spot-in-a-shared-design-document.md` | design | closed by the last task |
| `design/issues/summary.md` | design | regenerated by the close |

No folder touch is kept. No `after`: the story reads nothing another open story writes.

Figures:

- Forecast 1h30m, delivery 2026-11-10T01:15:00Z. `flai forecast` gave 3m from size 2 (2 criteria, 0 touches) at 74 s a unit, before any touch was declared; that size does not count an ADR, a Go change, two real-git tests, and four documents. 1h30m matches the stories of this shape planned today. The delivery is flai's 23:45Z, 8th in the pull order, moved by the 1h27m added.
- Cost of delay 175.00 USD a week, `flai cod`'s figure from the operator's input (1h10m lost a cycle at 150 USD an hour), kept. I-0315 now counts 3 occurrences where the input counted 2; the input is the operator's, so the thread names it rather than changing it.
