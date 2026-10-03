---
id: T-0764
type: task
nature: feature
title: The card menu's entries come from one pure function, and the story page's agent actions share its rule
status: in-progress
parent: S-0202
owner: alex
created: 2026-10-03T18:32:18Z
updated: 2026-10-03T18:32:55Z
transitions:
  - to: ready
    at: 2026-10-03T18:32:54Z
    by: agent-S-0202
  - to: in-progress
    at: 2026-10-03T18:32:55Z
    by: agent-S-0202
stream: S-0202
tags: []
touches: [flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/StoryAgent.svelte]
---

# T-0764 The card menu's entries come from one pure function, and the story page's agent actions share its rule

## Work

Add `flaiover/src/lib/cardmenu.ts`: `cardMenu(card, ctx)` returns the card menu's entries in order, each `{ action, label }`: `open` (Open), `finalize` (Finalize: a draft story, writable, not archived), the agent action (Start agent or Retry, as the story page offers them), `block` (Block…) or `unblock` (Unblock) on an open item, `cancel` (Cancel…) on an item not done or cancelled. The context gives `writable` and what the story's agent is doing and whether the agent host action is on.

Move the rule `StoryAgent.svelte` uses for Start agent, Retry, and Start agent here into one exported function in `$lib/activity.ts` (`agentAction`), so the page and the menu offer the same action for the same story; `StoryAgent.svelte` keeps its own one-shot `retried` state and calls it. Waits for nothing: it is the first layer.

## Done when

- `cardMenu` and `agentAction` have behaviour tests in `cardmenu.test.ts` and `activity.test.ts`: Finalize on a writable draft story only, never on a finalized story, a task, an epic, an archived card, or a read-only board; Block/Unblock and Cancel by state; Start agent and Retry as `StoryAgent` gave them before.
- `StoryAgent.svelte.test.ts` still passes unchanged.

## Notes
