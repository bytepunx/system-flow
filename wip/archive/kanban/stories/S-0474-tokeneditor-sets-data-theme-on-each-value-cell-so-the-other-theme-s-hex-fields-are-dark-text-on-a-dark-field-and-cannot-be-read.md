---
id: S-0474
type: story
nature: remediation
title: TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
status: done
owner: alex
created: 2026-10-21T11:59:01Z
updated: 2026-10-21T12:24:03Z
transitions:
  - to: ready
    at: 2026-10-21T12:10:31Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-21T12:10:38Z
    by: agent-S-0474
  - to: review
    at: 2026-10-21T12:17:16Z
    by: agent-S-0474
  - to: done
    at: 2026-10-21T12:24:03Z
    by: orchestrator
tags: [dashboard]
touches: [flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/design.svelte.test.ts, design/system/component-library.md, design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md, design/issues/summary.md, design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 411
  priced_by: harness
  turns:
    - day: 2026-10-21
      ceremony: 1
      hand_edits: 1
      work: 30
  models:
    - model: claude-opus-5-5
      input: 66
      output: 10397
      cache_read: 4028848
      cache_write: 143122
      cost: 2.1589
  strategic:
    - kind: orchestrator
      seconds: 290
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 115
          output: 1921
          cache_read: 23764594
          cache_write: 29363
          cost: 5.254
        - model: claude-sonnet-5-5
          input: 6
          output: 44
          cache_read: 48039
          cache_write: 29056
          cost: 0.0531
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-21T11:59:01Z
  value: 12.5
  by: planner-S-0474
  at: 2026-10-21T12:07:34Z
forecast:
  duration: 16m
  delivery: 2026-10-22T04:53:00Z
  basis: "flai forecast: median 115 s per unit of size over 27 done remediation stories on claude-opus-5-5 in the medium band, times size 8 (2 criteria, 6 touches); 37th in the pull order with an in-progress limit of 5."
  by: planner-S-0474
  at: 2026-10-21T12:07:34Z
finalized:
  by: alex
  at: 2026-10-21T12:04:23Z
---
# S-0474 TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

## Goal

This story remediates [I-0159](../../../design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md), "TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0159 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0159 is closed with `flai issue close I-0159 --reason` saying what fixed it

## Tasks
- T-1985 TokenEditor marks each value cell's theme with an attribute of its own, not data-theme, and a test fails while any element in the editor carries data-theme
- T-1988 The component library's findings say TokenEditor's cells no longer set data-theme, and I-0159 is closed with what fixed it

## Notes

Cost of delay inputs set by flai from I-0159. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-17T05:32:26Z, 4.3 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the issue's instance and the code: line 371 of `flaiover/src/lib/components/design/TokenEditor.svelte` sets `data-theme` on each value cell. `flaiover/src/routes/layout.css` scopes the dark palette and Tailwind's `dark:` variant to `[data-theme='dark']` on any element, so a dark cell's inputs take the dark field while their text keeps the light page's colour. The swatch paints the token's own CSS value and needs no scope; the attribute is only the tests' hook for finding a cell. So the cell gets an attribute no stylesheet reads, the tests find cells by it, and a new test fails while any element in the editor carries `data-theme`.

Touches. The story declared none, and `flai touches suggest` on `TokenEditor.svelte` found no file changed with it often enough (1 of 2203 commits). All are files; no folder touch is kept.

| Touch | Source | Why |
|-------|--------|-----|
| `flaiover/src/lib/components/design/TokenEditor.svelte` | design | The component I-0159 names, which `design/system/component-library.md` § What the pass found records |
| `flaiover/src/lib/components/design/TokenEditor.svelte.test.ts` | layout | Its `cell()` helper finds a cell by `[data-theme]`; the reproducing test goes here |
| `flaiover/src/routes/design/design.svelte.test.ts` | layout | Its `muted()` helper finds a TokenEditor cell by `[data-theme="light"]` |
| `design/system/component-library.md` | design | Its findings list I-0159 under TokenEditor and record each fix since |
| `design/issues/I-0159-…-cannot-be-read.md` | design | Criterion 2: `flai issue close` writes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

The examples file `design/components/flaiover/src/lib/components/design/TokenEditor.examples.json` is not touched: the fix changes no prop.

Forecast: 16m, delivery 2026-10-22T04:53:00Z, as `flai forecast` gives them with the six touches; it gave 5m on size 2 before the touches were written. It stands. Two of the six touches come from one `flai issue close`, which overstates the size. The flaiover close-out tiers that a three-line change still runs offset that.

Cost of delay: 12.50 USD a week, as `flai cod` works it out from the input flai set from I-0159 (5m lost per 168h cycle at 150 USD an hour). It stands. The defect has one instance, but the operator meets it on the Design page whenever they edit a dark token in the light theme, so the value is a floor.

### Accepted by the orchestrator

- Verified: 011a470c20627226802918ad524600dfd3f9512a
- Batch: 20261021-122058
- Enqueued verified: db0950fe1f1183bcea6ddd5116ec675c1b77ab22
- At: 2026-10-21T12:24:03Z

Verdict: pass. Both criteria are met, every changed file is within the story's touches, and no props changed.
- 1: flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/design.svelte.test.ts, design/system/component-library.md
- 2: design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md, design/issues/summary.md
