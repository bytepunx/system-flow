---
id: T-0768
type: task
nature: feature
title: The new-item and edit forms have a collapsed Cost of delay panel
status: done
parent: S-0204
owner: alex
created: 2026-10-03T18:35:38Z
updated: 2026-10-03T18:47:24Z
transitions:
  - to: ready
    at: 2026-10-03T18:36:18Z
    by: agent-S-0204
  - to: in-progress
    at: 2026-10-03T18:36:19Z
    by: agent-S-0204
  - to: done
    at: 2026-10-03T18:47:24Z
    by: agent-S-0204
stream: S-0204
tags: []
touches: [flaiover/src/lib/components/CostOfDelayPanel.svelte, flaiover/src/lib/components/CostOfDelayPanel.svelte.test.ts, flaiover/src/lib/components/NewItemForm.svelte, flaiover/src/lib/components/NewItemForm.svelte.test.ts, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/components/ItemEditor.svelte.test.ts, flaiover/src/routes/new, flaiover/src/routes/edit, flaiover/src/routes/api/items/+server.ts, "flaiover/src/routes/api/items/[id]/edit", flaiover/src/lib/costofdelay.ts, flaiover/src/lib/costofdelay.test.ts]
usage:
  source: log
  seconds: 665
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 1518
      cache_read: 3969520
      cache_write: 110055
      cost: 1.6536
---
# T-0768 The new-item and edit forms have a collapsed Cost of delay panel

## Work

A `CostOfDelayPanel` component, collapsed by default, with revenue per week and penalty per week labelled in the project's currency and time lost per cycle as a duration, each optional; collapsed, it shows a one-line summary when any value is set. `NewItemForm` and `ItemEditor` use it for epics and stories. The edit form sends `cost_of_delay` with only the keys that changed, an emptied key as an empty string, so clearing all three removes the inputs and keeps the planner's value; the new form sends the keys given to `item.new`. The edit form takes the currency from the edit view (`currency`); the new form from the manifest's `planning.currency`, default USD. Waits for nothing: it shares no path with the flai task, and its tests mock the API.

## Done when

- The panel is collapsed by default in both forms for epics and stories, absent for tasks, with the summary line when collapsed and set.
- Tests cover saving inputs from each form, clearing all three (the request empties the three keys and leaves the value), and the summary line.

## Notes
