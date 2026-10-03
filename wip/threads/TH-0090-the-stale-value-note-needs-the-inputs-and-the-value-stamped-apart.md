---
id: TH-0090
title: The stale-value note needs the inputs and the value stamped apart
anchor:
  path: wip/kanban/stories/S-0204-the-new-item-and-edit-forms-take-cost-of-delay-inputs-in-a-collapsed-panel.md
  item: S-0204
status: open
participants: [agent-S-0204]
created: 2026-10-03T18:34:44Z
updated: 2026-10-03T18:34:44Z
---

# TH-0090 The stale-value note needs the inputs and the value stamped apart

On wip/kanban/stories/S-0204-the-new-item-and-edit-forms-take-cost-of-delay-inputs-in-a-collapsed-panel.md.

## Entries

### 2026-10-03T18:34:44Z agent-S-0204
S-0204 asks the item page to show the planner's `value` with its `by` and `at`, and to note when the value is stale (inputs newer than the value). ADR-0074 gives `cost_of_delay` one `by` and `at`, the block's last editor: once the operator saves inputs, the planner's stamp is overwritten, so neither the value's `by`/`at` nor "inputs newer than the value" can be told.

My recommendation: `inputs` get their own `by` and `at`, stamped when an input changes; the block's `by` and `at` become the value's, stamped only when `value` changes. The value is stale when it exists and `inputs.at` is later than `at`. A new ADR records it (amending ADR-0074's "one by and at per block"), with work-hierarchy.md, validation (a block with inputs needs inputs.by/at; one with a value needs by/at), and front-matter-fields.txt. A file written by today's flai (one stamp) reads as: the stamp belongs to the value when there is a value, else to the inputs.

Clearing all three inputs keeps the value (a value alone is legitimate: an epic's apportioned share) with no stale note, since nothing newer is left; the planner recomputes on its next run.

Alternative: keep one stamp and show the note whenever the block's `by` is not the planner. That needs no schema change but is wrong whenever anyone else sets a value, and loses the planner's stamp.

Separately, and decided unless you say otherwise: the new-item form saves cost of delay in the same commit that makes the item, by giving `flai story new`/`flai epic new` the same `--revenue-per-week`, `--penalty-per-week`, `--time-lost-per-cycle` flags as `flai edit` and hostapi `item.new` a `cost_of_delay` (create.go already takes one since S-0203). That reaches `flai/cmd` beyond the story's declared touches.

Shall I go with the separate stamps?
