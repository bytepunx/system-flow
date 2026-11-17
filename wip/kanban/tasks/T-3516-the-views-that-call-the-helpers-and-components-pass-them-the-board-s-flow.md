---
id: T-3516
type: task
nature: feature
title: The views that call the helpers and components pass them the board's flow
status: in-progress
parent: S-0758
owner: alex
created: 2026-11-17T08:12:42Z
updated: 2026-11-17T08:12:55Z
transitions:
  - to: ready
    at: 2026-11-17T08:12:55Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T08:12:55Z
    by: system-flow
stream: S-0758
tags: []
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/cardmenu.svelte.test.ts, flaiover/src/routes/design/components/+page.svelte, flaiover/src/routes/design/components/components.svelte.test.ts, flaiover/src/lib/cardmenu.ts, flaiover/src/lib/cardmenu.test.ts, flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/components/ActivityView.svelte, flaiover/src/lib/components/ActivityView.svelte.test.ts, flaiover/src/routes/activity/+page.svelte, "flaiover/src/routes/docs/[...path]/+page.svelte", flaiover/src/lib/components/DocEditor.svelte, flaiover/src/lib/components/DocEditor.svelte.test.ts, flaiover/src/lib/edit.test.ts, flaiover/src/lib/components/DeploymentRow.svelte, "flaiover/src/routes/items/[id]/+page.svelte"]
after: [T-3255, T-3256, T-3257]
---

# T-3516 The views that call the helpers and components pass them the board's flow

## Work

Waits for T-3255, T-3256, and T-3257: they gave `touching`, the agent-action helpers, `BoardCard`, `CancelConfirm`, and `openStories` a flow, and their callers in other files do not pass one yet, so with the empty default those views offer nothing a state's kind decides. Make the flow required where a caller must give it, so svelte-check names every caller, and pass the board's flow from each: the board page, the components page, the card menu (`cardmenu.ts`, which also compares done and cancelled itself), `StoryAgent`, `ActivityView`, the docs page, `DocEditor`, the item page, and `DeploymentRow`, which names columns by state.

## Done when

No caller relies on an empty flow, none of these files compares a status with a state's name, and the flaiover tests pass.

## Notes
