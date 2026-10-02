---
id: S-0201
type: story
nature: feature
title: "The story page shows [Draft] and finalizes a draft story, and cards mark drafts"
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-02T11:54:38Z
transitions: []
tags: [dashboard, flai]
touches: [flaiover/src/routes/items, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/routes/api/items, flai/internal/hostapi]
after: [S-0199]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0201 The story page shows [Draft] and finalizes a draft story, and cards mark drafts

## Goal

Stories the planner or the issue step writes carry `draft: true`. The operator needs to see which stories are drafts and to accept one as their own with one action.

## Acceptance criteria
- [ ] The story page shows a `[Draft]` indicator beside the title when `draft` is true, and a **Finalize** button that clears it through a hostapi write (`item.finalize`), removing the indicator without a reload
- [ ] Board cards, the items list, and search hits show a draft marker (text, not colour alone)
- [ ] Finalizing records who did it in the item's transitions or `updated` as an edit by the operator, so the metrics can tell planner drafts accepted from drafts rewritten
- [ ] The edit form warns before a draft is moved to ready and offers Finalize there
- [ ] `design/system/flaiover-dashboard.md` and the user guide describe it; tests cover the indicator, the button, and the card marker

## Tasks

## Notes
