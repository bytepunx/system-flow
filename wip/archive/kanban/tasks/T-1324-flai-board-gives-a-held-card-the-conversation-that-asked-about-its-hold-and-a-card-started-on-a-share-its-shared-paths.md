---
id: T-1324
type: task
nature: feature
title: flai board gives a held card the conversation that asked about its hold, and a card started on a share its shared paths
status: done
parent: S-0338
owner: alex
created: 2026-10-08T04:31:58Z
updated: 2026-10-08T10:38:57Z
transitions:
  - to: ready
    at: 2026-10-08T10:30:13Z
    by: agent-S-0338
  - to: in-progress
    at: 2026-10-08T10:30:13Z
    by: agent-S-0338
  - to: done
    at: 2026-10-08T10:38:57Z
    by: agent-S-0338
stream: S-0338
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/workitem/share.go, flai/internal/workitem/share_test.go, flaiover/src/lib/activity.test.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/components/BoardCard.svelte, flaiover/src/lib/components/BoardCard.svelte.test.ts, flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts, flaiover/src/lib/server/board.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts]
usage:
  source: log
  seconds: 524
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 16483
      cache_read: 3097111
      cache_write: 91467
      cost: 1.5404
---
# T-1324 flai board gives a held card the conversation that asked about its hold, and a card started on a share its shared paths

## Work

Give the dashboard what the card shows, from flai, so the board, `board.get`, and the MCP `board` agree. It waits for nothing in this story.

- A ready card held on overlap carries the conversation flai serve opened to ask its holding story's agent (S-0334), beside `held: {code, reason}`.
- A card of a story started on a share carries the shared paths, the split, and the conversation.
- `workitem` cannot import `messages`, which imports it: take the asks and the shares as S-0334 passes shares into `Holds`, and widen this task's touches to each caller that has to pass them.

## Done when

- Tests cover a held card that was asked, one that was not, a card started on a share, and a share that ended.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
