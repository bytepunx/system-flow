---
id: T-0834
type: task
nature: feature
title: The card menu offers Plan on an open epic or story while the plan host action is on, and starts the planner
status: done
parent: S-0263
owner: alex
created: 2026-10-04T23:39:22Z
updated: 2026-10-04T23:43:00Z
transitions:
  - to: ready
    at: 2026-10-04T23:39:44Z
    by: agent-S-0263
  - to: in-progress
    at: 2026-10-04T23:39:45Z
    by: agent-S-0263
  - to: done
    at: 2026-10-04T23:43:00Z
    by: agent-S-0263
stream: S-0263
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 195
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 11843
      cache_read: 1762493
      cache_write: 68883
      cost: 1.0412
    - model: claude-sonnet-5
      input: 34
      output: 7554
      cache_read: 677839
      cache_write: 60117
      cost: 0.3615
---

# T-0834 The card menu offers Plan on an open epic or story while the plan host action is on, and starts the planner

## Work

Add a `plan` entry to `cardMenu` in `flaiover/src/lib/cardmenu.ts`, labelled **Plan**, offered on an epic or a story by the rule the item page shows `PlanAction` by: the board writable, the item not archived, not done or cancelled, and the `plan` host action on. The board learns whether that action is on from `GET /api/host-agent`, which adds `plan_enabled` from `project.info`'s `host_actions.plan`, false when flai cannot be asked. `CardMenu.svelte` passes it through, and the board's `pickCard` posts `POST /api/items/:id/plan` and puts what flai answered in the notice, as `PlanAction` words it, then asks for the agents again. Waits for nothing: it is the first task.

## Done when

- `cardMenu` lists Plan after the agent entry on an open epic or story while plan is enabled and writable, and not otherwise; `cardmenu.test.ts` covers each case.
- Picking Plan on the board posts to the plan endpoint and the notice says the planner started, or why flai refused; `routes/board/cardmenu.svelte.test.ts` covers both.
- The flaiover tests for those files and `svelte-check` pass.

## Notes
