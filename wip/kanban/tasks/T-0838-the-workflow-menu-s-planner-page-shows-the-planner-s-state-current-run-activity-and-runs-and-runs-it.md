---
id: T-0838
type: task
nature: feature
title: The Workflow menu's Planner page shows the planner's state, current run, activity, and runs, and runs it
status: done
parent: S-0259
owner: alex
created: 2026-10-05T00:07:53Z
updated: 2026-10-05T00:21:54Z
transitions:
  - to: ready
    at: 2026-10-05T00:08:20Z
    by: agent-S-0259
  - to: in-progress
    at: 2026-10-05T00:13:54Z
    by: agent-S-0259
  - to: done
    at: 2026-10-05T00:21:54Z
    by: agent-S-0259
stream: S-0259
tags: []
touches: [flaiover/src/routes/workflow, flaiover/src/lib/components, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/usage.ts]
after: [T-0837]
usage:
  source: log
  seconds: 480
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 92
      output: 1249
      cache_read: 4346355
      cache_write: 141469
      cost: 1.8255
---
# T-0838 The Workflow menu's Planner page shows the planner's state, current run, activity, and runs, and runs it

## Work

`/workflow/planner`, a page in the Workflow tier of the site menu (`flaiover/src/lib/sitemenu.ts`), built on what T-0837 answers:

- The `plan` host action's state: on, or off with the command that turns it on (`flai serve enable plan`, on the host, in the project).
- The current run, if any: who, with what harness and model, since when, for which item and why (its trigger), and its stream in the pane story pages have (`AgentStream.svelte`, told to read a planner run's stream).
- The activity document's totals (accrued cost, accrued seconds, activities, last run) and its entries newest first, each with its summary, trigger, items linked, seconds, and cost.
- The past runs newest first, each with its item, start, end, cost, and outcome.
- A form taking an epic's or a story's ID that runs the planner through `POST /api/items/[id]/plan` and says what flai answered; while the action is off it is disabled and says why.
- It loads again when a narrative changes, which `wip/agents/planner.md` is (`follow(['narrative'])`), when a work item changes, and when flai serve says a planner run started or ended.

It waits for T-0837, whose routes and `planner.ts` it reads.

## Done when

- [x] The page and its component show each part above, with component tests for the state, the runs, the entries, the live reload, and the form, on and off
- [x] The site menu lists Planner in the Workflow tier, with its test
- [x] `npx vitest run` passes for the files it changed, and `npm run check` reports nothing in them

## Notes
