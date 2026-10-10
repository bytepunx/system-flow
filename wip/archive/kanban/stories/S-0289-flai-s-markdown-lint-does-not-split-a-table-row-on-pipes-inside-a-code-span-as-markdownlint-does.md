---
id: S-0289
type: story
nature: remediation
title: flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does
status: done
owner: alex
created: 2026-10-06T09:56:52Z
updated: 2026-10-10T18:58:58Z
transitions:
  - to: ready
    at: 2026-10-09T18:23:35Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-09T18:23:41Z
    by: agent-S-0289
  - to: review
    at: 2026-10-09T18:33:36Z
    by: agent-S-0289
  - to: done
    at: 2026-10-10T18:58:58Z
    by: alex
tags: [flai]
touches: [flai/internal/mdlint/mdlint.go, flai/internal/mdlint/rules.go, flai/internal/mdlint/doc.go, flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/testdata/cases/tables.md, flai/internal/mdlint/testdata/cases/expected.txt, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0077-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 608
  turns:
    - day: 2026-10-09
      ceremony: 1
      hand_edits: 1
      work: 18
  models:
    - model: claude-opus-5-5
      input: 100
      output: 33836
      cache_read: 4411281
      cache_write: 236761
      cost: 3.1733
  strategic:
    - kind: planner
      seconds: 193
      estimated: true
      models:
        - model: claude-haiku-4-5-20251001
          input: 13
          output: 390
          cache_read: 56152
          cache_write: 10835
          cost: 0.0211
        - model: claude-opus-5-5
          input: 76
          output: 20211
          cache_read: 3350012
          cache_write: 213120
          cost: 2.8039
    - kind: orchestrator
      seconds: 5424
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 912
          cache_read: 5835635
          cache_write: 55826
          cost: 1.4854
        - model: claude-sonnet-5-5
          input: 10
          output: 60
          cache_read: 164368
          cache_write: 16296
          cost: 0.146
cost_of_delay:
  inputs:
    time_lost_per_cycle: 6m
    by: flai
    at: 2026-10-06T09:56:52Z
  value: 15
  by: planner-S-0289
  at: 2026-10-09T18:23:07Z
forecast:
  duration: 22m
  delivery: 2026-10-09T21:32:00Z
  basis: "flai forecast: median 107 s per unit of size over 21 done remediation stories on claude-opus-5-5, times size 12 (2 criteria, 10 touches); the four earlier mdlint rule fixes S-0265, S-0262, S-0240, S-0324 took 12m, 16m, 21m, and 33m"
  by: planner-S-0289
  at: 2026-10-09T18:23:07Z
finalized:
  by: alex
  at: 2026-10-07T02:19:47Z
---
# S-0289 flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does

## Goal

This story remediates [I-0077](../../../design/issues/I-0077-flai-s-markdown-lint-does-not-split-a-table-row-on-pipes-inside-a-code-span-as-markdownlint-does.md), "flai's markdown lint does not split a table row on pipes inside a code span, as markdownlint does". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0077 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0077 is closed with `flai issue close I-0077 --reason` saying what fixed it

## Tasks
- T-1433 flai's markdown lint reports MD056, table column count, on a row whose unescaped pipes, those inside a code span among them, give it more or fewer cells than its header
- T-1434 I-0077 is closed with a reason naming MD056 as what fixed it

## Notes

Cost of delay inputs set by flai from I-0077. time_lost_per_cycle 6m: 6m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-05T03:58:53Z, 1.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0077's one instance: flai's `mdlint` gains MD056, table column count. A reproduction of the instance's row (commit adf27c49 escaped it) on today's `mdlint` and on markdownlint-cli2 0.20.0 shows the cause is narrower than the title:

| Linter | Findings on the row |
|--------|---------------------|
| markdownlint-cli2 0.20.0 | MD038 on the code span the split leaves; MD056, Expected: 2; Actual: 9 |
| flai `mdlint` | MD038 on the same code span |

`doc.go`'s `cells` already splits a row on every unescaped pipe, inside a code span too, and MD038 (S-0262) now agrees. What flai lacks is MD056, so a row the split gives too many cells still passes. MD055 and MD058, the other table rules the configuration turns on, are left out: I-0077 names neither.

Touches, all files, no folder touch:

| Touch | Source |
|-------|--------|
| `flai/internal/mdlint/mdlint.go` | layout: the rule list; co-change 60% |
| `flai/internal/mdlint/rules.go` | layout: where each rule is implemented |
| `flai/internal/mdlint/doc.go` | layout: where a table and its cells are recognised |
| `flai/internal/mdlint/mdlint_test.go` | co-change 80% |
| `flai/internal/mdlint/testdata/cases/tables.md` | layout: a new fixture, as S-0262 added `code-spans.md` |
| `flai/internal/mdlint/testdata/cases/expected.txt` | co-change 80%: regenerated by `make mdlint-fixtures` |
| `design/system/flai-cli.md` | design: the `flai check` row and the `mdlint/` layout line list the rules |
| `docs/users/flai.md` | design: the markdown lint paragraph names what flai checks |
| I-0077's file and `design/issues/summary.md` | criterion 2: `flai issue close` |

`flai touches suggest` from `rules.go` and `doc.go` also listed `flai/internal/mdlint/inline.go` (60%). It is left out: MD056 counts cells, which `doc.go` splits, and reads no inline content.

Forecast 22m stands. flai's first figure was 4m, from size 2 with no touches. With the 10 touches it gives 22m. The four earlier mdlint rule fixes took 12m, 16m, 21m, and 33m (S-0265, S-0262, S-0240, S-0324), median about 19m, so 22m fits. Delivery is flai's, 9th in the pull order.

Cost of delay value 15.00 USD a week stands, as `flai cod` works it out from the inputs: 6m lost per 168h cycle at 150 USD an hour. One instance in four days gives no reason to change it.
