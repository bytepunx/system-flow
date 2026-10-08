---
id: T-1326
type: task
nature: feature
title: A held card says its holding story's agent was asked, and a card started on a share names the shared paths
status: done
parent: S-0338
owner: alex
created: 2026-10-08T04:32:07Z
updated: 2026-10-08T10:39:17Z
transitions:
  - to: ready
    at: 2026-10-08T10:30:14Z
    by: agent-S-0338
  - to: in-progress
    at: 2026-10-08T10:30:15Z
    by: agent-S-0338
  - to: done
    at: 2026-10-08T10:39:17Z
    by: agent-S-0338
stream: S-0338
tags: [flaiover]
touches: [flaiover/src/lib/activity.ts, flaiover/src/lib/activity.test.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/server/board.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts]
usage:
  source: log
  seconds: 542
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 63
      output: 260
      cache_read: 2981937
      cache_write: 96378
      cost: 1.3787
---
# T-1326 A held card says its holding story's agent was asked, and a card started on a share names the shared paths

## Work

Put it on the card, where the operator watches holds. It waits for T-1324, whose board fields it reads.

- `activity.ts`: the hold's type gains the asking conversation, and a story's activity the share it started on.
- `BoardCard.svelte`: a held card whose holding agent was asked says so beside the hold line and links the conversation; a card started on a share names the shared paths and links the conversation.

## Done when

- Tests cover a held card asked and not asked, a card started on a share, and a card with neither.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
