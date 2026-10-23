---
id: D-0129
title: "The stale-value note needs the inputs and the value stamped apart: yes, separate stamps"
kind: answer
sources: [TH-0090]
topics: [orchestration, planning, analysis, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/routes/new, flaiover/src/routes/edit, flaiover/src/routes/items, flaiover/src/routes/api/items/+server.ts, flaiover/src/lib/components/CostOfDelayPanel.svelte, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/costofdelay.ts, flaiover/src/lib/planning.ts, flaiover/src/lib/planning.test.ts, flai/internal/hostapi, flai/cmd, flai/internal/workitem, flai/internal/itemedit, design/adrs, design/system/work-hierarchy.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md, "flaiover/src/routes/api/items/[id]/edit", flaiover/src/lib/costofdelay.test.ts, flaiover/src/lib/components/CostOfDelayPanel.svelte.test.ts, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/components/ItemEditor.svelte.test.ts, design/issues/I-0066-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md, design/issues/summary.md, docs/operators/settings.md, flai/internal/issues, flai/internal/mcpserver, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0129 The stale-value note needs the inputs and the value stamped apart: yes, separate stamps

The operator, alex, answered TH-0090 at 2026-10-03T19:15:29Z.

What was asked:

```text
S-0204 asks the item page to show the planner's `value` with its `by` and `at`, and to note when the value is stale (inputs newer than the value). ADR-0074 gives `cost_of_delay` one `by` and `at`, the block's last editor: once the operator saves inputs, the planner's stamp is overwritten, so neither the value's `by`/`at` nor "inputs newer than the value" can be told.
```

The answer:

```text
yes, separate stamps
```

Source: TH-0090.
