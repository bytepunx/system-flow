---
id: S-0201
type: story
nature: feature
title: "The story page shows [Draft] and finalizes a draft story, and cards mark drafts"
status: done
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T18:30:56Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:01Z
    by: alex
  - to: in-progress
    at: 2026-10-03T07:32:28Z
    by: agent-S-0201
  - to: review
    at: 2026-10-03T18:03:16Z
    by: agent-S-0201
  - to: done
    at: 2026-10-03T18:30:56Z
    by: alex
tags: [dashboard, flai]
touches: [flaiover/src/routes/items, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/routes/api/items, flai/internal/hostapi, flai/internal/search, flai/internal/workitem, flai/internal/itemedit, flai/cmd/edit.go, flai/cmd/items.go, flai/cmd/planning.go, flai/cmd/planning_test.go, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/ItemEditor.svelte, flaiover/src/lib/components/ItemEditor.svelte.test.ts, flaiover/src/lib/server/board.ts, flaiover/src/lib/server/repo-channel.test.ts, flaiover/src/lib/server/search.ts, flaiover/src/lib/server/search.test.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/search, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, design/adrs, design/issues, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/work-hierarchy.md]
after: [S-0199]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3771
  models:
    - model: claude-haiku-4-5-20251001
      input: 146
      output: 8953
      cache_read: 879957
      cache_write: 89951
      cost: 0.2453
    - model: claude-opus-5-5
      input: 618
      output: 167378
      cache_read: 36878817
      cache_write: 957422
      cost: 16.8698
    - model: claude-sonnet-5
      input: 86
      output: 21280
      cache_read: 2532418
      cache_write: 181244
      cost: 1.1726
---
# S-0201 The story page shows [Draft] and finalizes a draft story, and cards mark drafts

## Goal

Stories the planner or the issue step writes carry `draft: true`. The operator needs to see which stories are drafts and to accept one as their own with one action.

## Acceptance criteria
- [x] The story page shows a `[Draft]` indicator beside the title when `draft` is true, and a **Finalize** button that clears it through a hostapi write (`item.finalize`), removing the indicator without a reload
- [x] Board cards, the items list, and search hits show a draft marker (text, not colour alone)
- [x] Finalizing records who did it in the item's transitions or `updated` as an edit by the operator, so the metrics can tell planner drafts accepted from drafts rewritten
- [x] The edit form warns before a draft is moved to ready and offers Finalize there
- [x] `design/system/flaiover-dashboard.md` and the user guide describe it; tests cover the indicator, the button, and the card marker

## Tasks
- T-0753 Board cards and search hits carry a story's draft flag
- T-0754 The story page shows [Draft] beside the title and a Finalize button that clears it
- T-0755 Board cards, the overview's item list, and search hits mark a draft story
- T-0756 Finalizing a draft records who finalized it, and the hostapi op item.finalize runs it
- T-0757 The edit form warns that a draft cannot go to ready and offers Finalize
- T-0758 The design and the user guide describe the draft marker, Finalize, and who finalized

## Notes

- Who finalized is recorded in a `finalized: {by, at}` block on the story, not in `transitions` or `updated`: a finalize is not a move, and `updated` has no author (TH-0083: the designer took the recommendation; ADR-0077). `item.finalize` stamps the dashboard's owner.
- "The items list" is the epic page's list of its stories: the overview lists counts by type and state, not items.
- The edit form does not move items; it warns on a draft story that it cannot go to ready until finalized, and finalizes in place.
- MCP `item_get` does not return `finalized` yet: `flai/internal/mcpserver` is S-0200's while it is in progress. `flai show --json` returns it.
