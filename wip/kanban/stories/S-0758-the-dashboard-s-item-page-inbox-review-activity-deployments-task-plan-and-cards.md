---
id: S-0758
type: story
nature: feature
title: The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:04:18Z
updated: 2026-11-12T16:04:18Z
transitions: []
tags: [dashboard]
topics: [dashboard]
touches: [flaiover/src/lib/workflow.ts, flaiover/src/lib/workflow.test.ts, "flaiover/src/routes/items/[id]/+page.svelte", flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/lib/touches.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/taskplan.ts, flaiover/src/lib/inbox.svelte.ts, flaiover/src/lib/server/inbox.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/CancelConfirm.svelte, flaiover/src/lib/components/ForgeStatus.svelte, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/WaitTable.svelte, flaiover/src/lib/components/explorer/SourcePicker.svelte, flaiover/src/lib/viz/palette.ts, design/system/flaiover-dashboard.md]
after: [S-0753]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0758 The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name

## Goal

Beyond the board and the charts, about fifteen of flaiover's views and helpers compare an item's status with `'in-progress'`, `'review'`, or `'done'`. They read the state's kind from the columns `board.get` answers instead, so the item page's moves, the inbox, the review page, the activity feed, the Deployments page, the task plan, the cards, and the palette show an operator's own states for what they are.

## Acceptance criteria

- [ ] One module in `flaiover/src/lib` gives a state's kind, colour, and moves from the columns and transitions `board.get` answers, and every view below reads it
- [ ] `routes/items/[id]/+page.svelte`, `routes/ops/deployments/+page.svelte`, `touches.ts`, `activity.ts`, `taskplan.ts`, `inbox.svelte.ts`, `server/inbox.ts`, `BoardCard.svelte`, `CancelConfirm.svelte`, `ForgeStatus.svelte`, `Review.svelte`, `InboxView.svelte`, `WaitTable.svelte`, `explorer/SourcePicker.svelte`, and `viz/palette.ts` compare no status with a state's name
- [ ] The item page offers only the moves its transitions allow the operator, and asks for a reason where the transition needs one
- [ ] Tests render the item page, a card, and the review page for a story in a state added between `in-progress` and `review`
- [ ] With the baseline, every view looks and behaves as today and every existing flaiover test passes
- [ ] `design/system/flaiover-dashboard.md` says these views read kinds from flai

## Tasks

## Notes

- Found by a search for quoted state names under `flaiover/src` while planning E-0034; the board's own lists are S-0753's and the charts' S-0754's.
