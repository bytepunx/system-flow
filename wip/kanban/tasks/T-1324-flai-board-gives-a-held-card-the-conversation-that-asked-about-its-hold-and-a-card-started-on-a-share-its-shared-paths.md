---
id: T-1324
type: task
nature: feature
title: flai board gives a held card the conversation that asked about its hold, and a card started on a share its shared paths
status: backlog
parent: S-0338
owner: alex
created: 2026-10-08T04:31:58Z
updated: 2026-10-08T04:31:58Z
transitions: []
stream: S-0338
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go]
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
