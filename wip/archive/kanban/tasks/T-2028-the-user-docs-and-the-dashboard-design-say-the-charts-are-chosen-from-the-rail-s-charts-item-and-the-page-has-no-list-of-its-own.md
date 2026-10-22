---
id: T-2028
type: task
nature: improvement
title: The user docs and the dashboard design say the charts are chosen from the rail's Charts item and the page has no list of its own
status: done
parent: S-0484
owner: alex
created: 2026-10-21T21:00:45Z
updated: 2026-10-22T20:35:28Z
transitions:
  - to: ready
    at: 2026-10-22T20:19:33Z
    by: agent-S-0484
  - to: in-progress
    at: 2026-10-22T20:31:23Z
    by: agent-S-0484
  - to: done
    at: 2026-10-22T20:35:28Z
    by: agent-S-0484
stream: S-0484
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-2024]
usage:
  source: log
  seconds: 245
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 15127
      cache_read: 2367471
      cache_write: 84879
      cost: 1.2929
---
# T-2028 The user docs and the dashboard design say the charts are chosen from the rail's Charts item and the page has no list of its own

## Work

In `docs/users/flaiover.md`, § The header: say that Charts in the Flow group expands to list the charts under flow, usage, planning, and strategic, that clicking one opens it, that the chart shown is highlighted, and that the browser remembers whether Charts is open as it does the groups; § Charts: replace the sentences on the list beside the chart and the row above it on a narrow window with the rail. In `design/system/flaiover-dashboard.md` § Views: the rail's groups table row for Flow (Charts expands to its charts, grouped), the rail behaviour table (a page with children is a button; `locate`'s `child`; `charts` in `flaiover-nav-rail`), the page sidebar paragraph (the charts page no longer uses it; S-0461's link entries stay available), and the `/charts/*` rows where they mention the sidebar. Name `viz/catalog.ts` as where the catalogue lives. Bump each document's `updated`.

Fifth layer: waits for T-2026, which settles what the page shows.

## Done when

Both documents describe the rail and the page as built, the markdown lint and `flai check --strict` pass on them.

## Notes

Drafted by the planner for E-0028. S-0485 changes the same two documents for the page's order and waits for this story.
