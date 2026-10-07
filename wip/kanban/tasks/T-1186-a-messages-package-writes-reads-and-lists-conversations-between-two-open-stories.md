---
id: T-1186
type: task
nature: feature
title: A messages package writes, reads, and lists conversations between two open stories
status: in-progress
parent: S-0330
owner: alex
created: 2026-10-07T20:14:12Z
updated: 2026-10-07T20:26:25Z
transitions:
  - to: ready
    at: 2026-10-07T20:26:24Z
    by: agent-S-0330
  - to: in-progress
    at: 2026-10-07T20:26:25Z
    by: agent-S-0330
stream: S-0330
tags: [flai]
touches: [flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, wip/messages/README.md]
after: [T-1185]
---
# T-1186 A messages package writes, reads, and lists conversations between two open stories

## Work

Build the store T-1185's ADR decides, in a new package `flai/internal/messages`. It waits for T-1185, because the ADR settles the file format and the rules.

- `Send`, `Reply`, `Close`, `Get`, `List`, and `For(story)`, with the next ID, front matter, and dated entries, written atomically and linted as flai lints what it writes in `wip/` (ADR-0061).
- `Send` refuses a story that does not exist, and one not in progress or in review, on either side, naming its state.
- A conversation of a story no longer open reads as closed; `Awaiting(story)` says whether the last entry is the other story's.
- `wip/messages/README.md` says what the folder holds, if the ADR keeps messages there.

## Done when

- Tests cover sending, replying, listing per story, the refusals, and a conversation that reads as closed when a story leaves the open columns.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
