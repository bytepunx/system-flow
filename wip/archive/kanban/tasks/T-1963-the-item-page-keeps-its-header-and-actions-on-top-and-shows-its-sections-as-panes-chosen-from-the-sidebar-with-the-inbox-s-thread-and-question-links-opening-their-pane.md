---
id: T-1963
type: task
nature: feature
title: The item page keeps its header and actions on top and shows its sections as panes chosen from the sidebar, with the Inbox's thread and question links opening their pane
status: done
parent: S-0460
owner: alex
created: 2026-10-21T07:03:54Z
updated: 2026-10-21T11:36:27Z
transitions:
  - to: ready
    at: 2026-10-21T09:24:22Z
    by: agent-S-0460
  - to: in-progress
    at: 2026-10-21T09:24:22Z
    by: agent-S-0460
  - to: done
    at: 2026-10-21T11:36:27Z
    by: agent-S-0460
stream: S-0460
tags: [dashboard]
touches: ["flaiover/src/routes/items/[id]/+page.svelte", "flaiover/src/routes/items/[id]/item.svelte.test.ts"]
after: [T-1962]
usage:
  source: log
  seconds: 885
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 113
      output: 58798
      cache_read: 6662703
      cache_write: 226491
      cost: 3.9889
---
# T-1963 The item page keeps its header and actions on top and shows its sections as panes chosen from the sidebar, with the Inbox's thread and question links opening their pane

## Work

- In `flaiover/src/routes/items/[id]/+page.svelte`, keep the header (title, kind chips, status line), the notice, and the action buttons above everything.
- Replace the two-column grid with `PageSidebar` and one pane at a time, in the story's second criterion's order: Story, Tasks, Questions and threads, Messages, Agent, History, Files, and Planning and usage.
- Build the entries from the item. Leave out a pane with nothing to show, such as Tasks for an item with no children, or Messages for a story with no conversations. Give Questions and threads the count of open ones.
- Open Story on arrival. With `?thread=` or `?question=`, open Questions and threads, and keep today's scroll to the thread or the question.
- Update `item.svelte.test.ts` so that each existing case opens the pane that holds what it checks. Add cases for these:
  - the panes left out;
  - the arrival pane;
  - each inbox link;
  - a pane's hash.
- Waits for T-1962, the sidebar it uses.

## Done when

`item.svelte.test.ts` passes, and on a story with tasks, threads, and messages every section of today's page is reachable from the sidebar.

## Notes
