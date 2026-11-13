---
id: S-0758
type: story
nature: feature
title: The dashboard's item page, inbox, review, activity, deployments, task plan, and cards read a state's kind from flai instead of comparing its name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:04:18Z
updated: 2026-11-13T12:39:39Z
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
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 11
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 14
          cache_read: 121380
          cache_write: 20507
          cost: 0.0405
draft: true
forecast:
  duration: 44m
  delivery: 2026-11-13T23:18:00Z
  basis: "Its own forecast of 44m; 29th in the pull order with an in-progress limit of 3, behind S-0740, S-0767, S-0760, S-0745, S-0759, S-0751, S-0757, S-0661, S-0663, S-0662, S-0241, S-0730, S-0736, S-0737, S-0738, S-0739, S-0741, S-0742, S-0743, S-0744, S-0746, S-0747, S-0748, S-0749, S-0750, S-0752, S-0753, S-0754, S-0755 and S-0756."
  by: flai
  at: 2026-11-13T12:39:39Z
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

- T-3252 One flaiover module gives a state's kind, colour, and moves from the columns and transitions board.get answers
- T-3255 The item page and the Deployments page read a state's kind and offer only the moves the operator may make
- T-3256 The touches, activity, task plan, inbox, and palette helpers read a state's kind
- T-3257 The card, cancel, forge status, review, inbox, wait table, and source picker components read a state's kind, tested on a state of the workflow's own
- T-3258 flaiover-dashboard.md says the dashboard's views read a state's kind from flai

## Notes

- Found by a search for quoted state names under `flaiover/src` while planning E-0034; the board's own lists are S-0753's and the charts' S-0754's.

### Planning

- An addition to the epic's proposed order, which named the board and the charts but not these views. Proposed on the plan's thread.
- Touches, all files:
  - `flaiover/src/lib/workflow.ts`, `workflow.test.ts`: layout, a new module.
  - The fifteen files criterion 2 names: a search for `'in-progress'` and `'review'` under `flaiover/src`, the board's and the charts' files left out.
  - `design/system/flaiover-dashboard.md`: design (criterion 6).
  - The new component tests T-3257 adds sit beside their components and are named when that task writes them; `task_done` widens the touches then.
- `after`: S-0753. The module reads the columns, kinds, and transitions S-0753 adds to `board.get`.
- Tasks: T-3252 (the module) first; then T-3255 (pages), T-3256 (helpers), T-3257 (components), and T-3258 (design) together, each reading the module and sharing no path with another.
- Forecast 44m, as `flai forecast` gives it for size 24.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
