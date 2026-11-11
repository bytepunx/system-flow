---
id: S-0700
type: story
nature: remediation
title: The new item form's lint refusal is easy to miss, so a refused epic looks like a hang
status: ready
owner: alex
created: 2026-11-10T08:55:31Z
updated: 2026-11-11T06:15:20Z
transitions:
  - to: ready
    at: 2026-11-11T02:31:27Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/issues/I-0263-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md, design/issues/summary.md]
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
      seconds: 41
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 178
          output: 678
          cache_read: 9740104
          cache_write: 435662
          cost: 0.685
        - model: claude-opus-5-5
          input: 185
          output: 12947
          cache_read: 12479748
          cache_write: 343216
          cost: 4.2776
    - kind: orchestrator
      seconds: 38
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 20
          cache_read: 416086
          cache_write: 23811
          cost: 0.1206
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-10T08:55:31Z
  value: 50
  by: planner-S-0700
  at: 2026-11-10T10:27:56Z
forecast:
  duration: 45m
  delivery: 2026-11-11T07:44:00Z
  basis: "Its own forecast of 45m; 3rd in the pull order with an in-progress limit of 3, behind S-0686, S-0691, S-0693, S-0722, S-0696 and S-0697."
  by: flai
  at: 2026-11-11T06:15:20Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:40Z
---
# S-0700 The new item form's lint refusal is easy to miss, so a refused epic looks like a hang

## Goal

This story remediates [I-0263](../../../design/issues/I-0263-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md), "The new item form's lint refusal is easy to miss, so a refused epic looks like a hang". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0263 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0263 is closed with `flai issue close I-0263 --reason` saying what fixed it

## Tasks
- T-2966 Scroll to and focus the new item form's refusal, and mark the Create button refused
- T-2968 Say in the dashboard design and the users guide that a refusal on the new item page is scrolled to and focused
- T-2971 Close I-0263 with the fix as its reason

## Notes

Cost of delay inputs set by flai from I-0263. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-06T21:31:26Z, 3.5 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0263's one instance: `flaiover/src/lib/components/NewItemForm.svelte` draws its `refusal` alert between the body and preview pair (each `h-[55vh]`) and the Create button, and does nothing else when `POST /api/items` answers 422. A submit from the title field, with Enter, leaves the operator at the top of the form, where nothing changes: the button's `Creating…` and the alert are both below the fold. After any refusal the form scrolls the alert into view and moves focus to it (`tabindex="-1"`), and the Create button says the item was refused, with the alert as its description, until the text changes or the next submit. Nothing about flai's answer, the findings, or the props changes.

Touches:

| Path | Source | Why |
|------|--------|-----|
| `flaiover/src/lib/components/NewItemForm.svelte` | layout | the form and its refusal alert |
| `flaiover/src/lib/components/NewItemForm.svelte.test.ts` | layout | the test that reproduces the missed refusal |
| `design/system/flaiover-dashboard.md` | design | the `/new` row says what a refusal does |
| `docs/users/flaiover.md` | design | the new item paragraph says what a refusal does |
| `design/issues/I-0263-the-new-item-form-s-lint-refusal-is-easy-to-miss-so-a-refused-epic-looks-like-a.md` | declared (goal) | closed by criterion 2 |
| `design/issues/summary.md` | declared (goal) | regenerated when I-0263 closes |

`flai touches suggest` listed only hub files that change with most commits, such as `design/system/flai-cli.md` and `docs/users/flai.md`; none bears on this form, so none was added. The component's examples, `design/components/flaiover/src/lib/components/NewItemForm.examples.json`, are left out: the fix adds no prop.

Figures:

- Forecast 45m, adjusted from flai's 18m (134 s per unit of size times size 8, over 54 done remediation stories). The close-out runs the flaiover tier for a change under `flaiover/src`, and S-0669's flaiover close-out today took 20m44s on its own; the median does not weigh that. Delivery moved by the same 27m.
- Cost of delay 50 USD a week, as `flai cod` gives it: 20m lost per 168h cycle at 150 USD an hour. It stands: one instance, and the cost is the operator's time spent thinking the form hung.

No `after`: the story reads no other story's output.
