---
id: S-0719
type: story
nature: improvement
title: A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance
status: in-progress
owner: alex
created: 2026-11-11T02:29:04Z
updated: 2026-11-11T14:21:49Z
transitions:
  - to: ready
    at: 2026-11-11T08:52:14Z
    by: alex
  - to: in-progress
    at: 2026-11-11T14:21:49Z
    by: agent-S-0719
tags: []
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go, flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go, design/adrs/0286-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences.md, design/adrs/README.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md, design/issues/summary.md]
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
      seconds: 279
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 137
          output: 20853
          cache_read: 5098755
          cache_write: 333161
          cost: 0.3799
        - model: claude-opus-5-5
          input: 314
          output: 36350
          cache_read: 13136536
          cache_write: 490422
          cost: 5.2219
    - kind: orchestrator
      seconds: 59
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 122
          cache_read: 529608
          cache_write: 78910
          cost: 0.1698
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-11T02:29:04Z
  value: 37.5
  by: planner-S-0719
  at: 2026-11-11T08:54:09Z
forecast:
  duration: 1h
  delivery: 2026-11-11T11:25:00Z
  basis: "Its own forecast of 1h; 2nd in the pull order with an in-progress limit of 3, behind S-0694, S-0702, S-0706 and S-0710."
  by: flai
  at: 2026-11-11T09:18:38Z
finalized:
  by: alex
  at: 2026-11-11T08:50:48Z
---
# S-0719 A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance

## Goal

This story remediates [I-0346](../../../design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md), "A trial merge that finds two open stories editing different sentences of one markdown paragraph is settled only by sending the second story back at acceptance". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0346 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0346 is closed with `flai issue close I-0346 --reason` saying what fixed it

## Tasks
- T-3094 ADR: the section merge merges by sentences the lines of a paragraph both sides change, refining ADR-0175 and ADR-0276
- T-3097 The section merge merges by sentences the lines of a paragraph both sides change, and still stops when both change one sentence
- T-3099 Say in the design and the user docs that the section merge settles a paragraph both sides change in different sentences
- T-3102 Reproduce I-0346 in real git: two stories change different sentences of one paragraph of a design document, and neither sync nor acceptance stops
- T-3103 Close I-0346, saying the section merge now settles a paragraph two stories change in different sentences

## Notes

Cost of delay inputs set by flai from I-0346. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-11T01:43:39Z, 0 days before this story; under one cycle counts as one).

### Planning

The proposed fix, from I-0346's one instance: S-0688 and S-0666 each changed a different sentence of one paragraph in `design/system/agent-adapters.md` and `docs/operators/index.md`. Both files are under the design or docs folder, so the rebase at S-0688's acceptance already ran the section merge (ADR-0175, ADR-0276). It stopped because this repository writes each paragraph on one line: `mergeLines` in `flai/internal/storygit/mdmerge.go` sees both sides change the same base line and reports "both sides change the same lines". The resolution by hand was mechanical: main's paragraph with S-0688's sentence applied. So where `mergeLines` finds the two sides' hunks touching the same base lines, the section merge takes the lines they span and merges them again by sentences, with the same three-way merge on sentence tokens in place of lines. It keeps both changes when they touch different sentences, and still leaves the stop when both change one sentence. Prose lines and list items only: fenced code, table rows, and headings keep the line merge. The one change settles the stop in every rebase flai runs, `flai stream sync`, `flai task done`'s sync, `flai accept`, and the merge queue's batch, since all of them call `MergeMarkdown`.

The issue's other remedy, rebasing the second story as soon as the first is accepted (the `overlapped` change), is already the rule in `work-management.md`; it did not help because S-0688 was in review, with no agent running. The trial merge (`git merge-tree`) still reports the pair's conflict, as ADR-0276 leaves it: making it run the section merge is deferred there to a story of its own.

Touches, and where each came from (`flai touches suggest` listed nothing: the story declared no touches; seeded with the merge code and the two documents that describe it, it listed `design/system/flai-cli.md`, `design/adrs/README.md`, and `design/issues/summary.md` among its co-changes):

| Touch | Source | Why |
|---|---|---|
| `flai/internal/storygit/mdmerge.go` | layout | `mergeLines`, where the same-lines stop is decided |
| `flai/internal/storygit/mdmerge_test.go` | layout | `TestMergeMarkdownSettles` and `TestMergeMarkdownLeavesWhatItCannotSettle` |
| `flai/cmd/accept_conflict_test.go` | layout | the acceptance reproduction, beside S-0637's of I-0315 |
| `flai/cmd/stream_sync_test.go` | layout | the sync reproduction, beside S-0637's |
| `design/adrs/0286-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences.md` | design | the ADR refining ADR-0175 and ADR-0276; `flai adr new` gives the number and slug, so the agent renames this touch when it differs |
| `design/adrs/README.md` | co-change | generated from the ADR files |
| `design/system/workflow.md` | design | the paragraphs on ADR-0175 and ADR-0276 say both changing the same lines stops |
| `design/system/agent-coordination.md` | design | its ADR-0175 paragraph |
| `design/system/flai-cli.md` | co-change | the `flai accept` and `flai stream sync` rows name the same-lines stop |
| `docs/users/flai.md` | design | the Shared documents bullet under sync |
| `design/issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md` | declared by the goal | closed by the last task |
| `design/issues/summary.md` | co-change | written again when the issue closes |

No folder touch and no `after`: the story reads no other open story's output.

- Forecast 1h, delivery 2026-11-11T22:30:00Z. `flai forecast` gave 2m from size 2 (2 criteria, 0 touches) at 58 s a unit, before any touch was declared; that size does not count an ADR, new merge code, two real-git tests, and four documents. S-0637, the same shape, took 30m against 1h30m; the sentence merge is new code, so 1h. The delivery is flai's 21:30Z moved by the 58m added.
- Cost of delay 37.50 USD a week, `flai cod`'s figure from the operator's input (15m lost a cycle at 150 USD an hour), kept: I-0346 still counts one occurrence.
