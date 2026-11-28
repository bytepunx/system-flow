---
id: S-0853
type: story
nature: improvement
title: The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator
status: ready
owner: alex
created: 2026-11-27T22:55:25Z
updated: 2026-11-28T00:02:10Z
transitions:
  - to: ready
    at: 2026-11-28T00:02:04Z
    by: claude-concierge
tags: [flai]
touches: [design/adrs, design/adrs/README.md, flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/cmd/accept_uncommitted_test.go, flai/cmd/accept_orchestrator_dirty_test.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/hostapi/writes_test.go, flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/digest/decay.go, flai/internal/digest/decay_test.go, flai/internal/digest/write.go, design/system/flai-cli.md, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, design/issues/I-0411-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md, design/issues/summary.md]
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
      seconds: 16
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 14
          output: 2770
          cache_read: 950191
          cache_write: 26901
          cost: 0.4351
    - kind: analyzer
      seconds: 86
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5533
          cache_read: 391177
          cache_write: 29514
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
cost_of_delay:
  inputs:
    time_lost_per_cycle: 76h
    by: flai
    at: 2026-11-27T22:55:25Z
  value: 11400
  by: planner-S-0853
  at: 2026-11-27T23:58:18Z
forecast:
  duration: 2h
  delivery: 2026-11-28T04:12:00Z
  basis: "Its own forecast of 2h; 2nd in the pull order with an in-progress limit of 3, behind S-0837, S-0844, S-0846 and S-0852."
  by: flai
  at: 2026-11-28T00:02:10Z
finalized:
  by: alex
  at: 2026-11-28T00:01:53Z
---
# S-0853 The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator

## Goal

This story remediates [I-0411](../../../design/issues/I-0411-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md), "The acceptance's lint fix refuses every acceptance while the curator's uncommitted digest entries in the main checkout fail markdownlint, so stories that passed review wait for the operator", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0411 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0411 is closed with `flai issue close I-0411 --reason` saying what fixed it

## Tasks
- T-3631 Record in an ADR that an acceptance no longer refuses for unfixable lint in a shared-set file the story does not change, and that digest absorb and compact lint what they rewrite
- T-3636 flai accept leaves out of its commits a shared-set file the story does not change whose lint the fix cannot fix, names it, and no longer refuses for it
- T-3638 The merge queue's preflight holds the queue on unfixable lint in a shared-set file only when an enqueued story changes it
- T-3639 flai digest absorb and compact lint the entry they rewrite and refuse a rewrite that brings a markdown finding, as new and edit do
- T-3640 Document that an acceptance leaves unfixable lint in a file the story does not change uncommitted, and that digest absorb and compact lint, then close I-0411

## Notes

The analyzer found I-0411: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0411. time_lost_per_cycle 76h carried over from I-0411's Impact section. Its Impact time_lost_per_cycle was taken rather than the 76h derived from its cost and count.

Strategic usage carried over from I-0411: analyzer 1m26s, 1.4522 USD.

### Planning

Where the refusal comes from: `flai accept` runs the lint fix over the main checkout's shared set (`wip/`, `design/`, and `claims.shared`) at `flai/cmd/accept.go` through `previewedToLint` and `refuseUnfixedUncommitted` in `flai/cmd/accept_markdown.go`, whether or not the story changes those files (ADR-0140 §1). Every uncommitted digest entry is in that set. The merge queue's landing, `flai accept --batch`, and the orchestrator's MCP `accept` run the same code. The merge queue's own preflight (`preflightLint` and `confirmLint` in `flai/internal/mergequeue/preflight.go`, ADR-0234) holds the queue on the same files. `flai digest absorb` and `compact` rewrite entries through `rewrite` in `flai/internal/digest/decay.go`, which validates but does not lint, unlike `new` and `edit` (ADR-0318).

The proposed fix, from I-0411's instances, for the story's agent to record in an ADR and build:

1. Every acceptance's preflight lint stops refusing for a shared-set file the story's change does not touch. A file the fix leaves failing is left out of the acceptance's commits and stays uncommitted, as ADR-0248 leaves "the rest". The acceptance names it in its output and its preview. A file the story's change touches still refuses, as ADR-0140 §2 has it.
2. The merge queue's preflight applies the same rule to a file no enqueued story changes, so the queue is not held on it.
3. `flai digest absorb` and `compact` lint what they rewrite, as ADR-0318 makes `new` and `edit` do. Then no flai write puts unlintable markdown in the digest.

Touches, and where each came from:

| Touch | Source |
|-------|--------|
| `flai/cmd/accept_markdown.go` | layout: the refusal "the main checkout's uncommitted files hold lint the fix … could not fix" is made here |
| `flai/cmd/accept.go`, `flai/cmd/accept_fix_test.go`, `flai/cmd/accept_uncommitted_test.go`, `flai/cmd/accept_orchestrator_dirty_test.go` | co-change with `accept_markdown.go` (6, 4, 2, 2 of 12 commits); `accept.go` calls the lint and holds the `--help` text; the tests hold the refusal text |
| `flai/cmd/accept_batch_dirty.go`, `flai/cmd/accept_batch_dirty_test.go` | co-change (3 of 12): the queue's landing on a dirty main, `leftUncommitted` |
| `flai/internal/preview/accept.go`, `flai/internal/preview/accept_test.go` | co-change and design: the preview lists `shared_uncommitted` and `uncommitted` (ADR-0248) and gains the files left out |
| `flai/internal/hostapi/writes_test.go` | layout: it asserts the refusal text |
| `docs/users/flai-reference.md` | co-change (3 of 12): generated from `flai accept --help`, which changes |
| `flai/internal/mergequeue/preflight.go`, `flai/internal/mergequeue/preflight_test.go` | co-change (4 of 12) and design: ADR-0234's preflight lints main's uncommitted files |
| `flai/internal/digest/decay.go`, `flai/internal/digest/decay_test.go`, `flai/internal/digest/write.go` | layout: `Absorb`, `Compact`, and `rewrite` in `decay.go`; the `lint` helper ADR-0318 added is in `write.go` |
| `design/adrs`, `design/adrs/README.md` | design: the decision refines ADR-0140 and ADR-0318 in a new ADR |
| `design/system/flai-cli.md`, `design/system/workflow.md`, `design/system/strategic-agents.md`, `docs/users/flai.md`, `template/CHANGELOG.md` | design and co-change: The acceptance gate and The merge queue pass, What waits on an acceptance, The curator, and the user docs of `flai accept` and `flai digest` |
| `design/issues/I-0411-…md`, `design/issues/summary.md` | criterion 2 |

The folder touch `design/adrs` is kept: the story adds one ADR whose file name `flai adr new` gives from its number and title, which no task can name yet. `flai task done` narrows it to the file once written.

No `after`: the story reads no output of an open story. S-0836's ADR-0318 is on main.

Tagged `flai`: the change is flai's; `template/CHANGELOG.md` only records it.

Figures:

- Cost of delay value 11,400 USD a week, as `flai cod` gives it from the input I-0411 carried: 76 h lost per 168 h cycle at 150 USD an hour. It stands: 76 h is the shortest of the five waits of 2026-11-18, and the hazard is open today while the curator commits nothing.
- Forecast 2h, raised from `flai forecast`'s 48m (109 s per unit of size over 187 large improvement stories, size 26). The plan has five tasks in three layers, each tested with `flai test`, and a flai close-out takes about 12m, as S-0847's and S-0808's did. Delivery moves by the same 72m, to 2026-11-28T08:30:00Z.
