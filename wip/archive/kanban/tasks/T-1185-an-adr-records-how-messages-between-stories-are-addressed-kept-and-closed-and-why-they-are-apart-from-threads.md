---
id: T-1185
type: task
nature: feature
title: An ADR records how messages between stories are addressed, kept, and closed, and why they are apart from threads
status: done
parent: S-0330
owner: alex
created: 2026-10-07T20:14:06Z
updated: 2026-10-07T20:26:04Z
transitions:
  - to: ready
    at: 2026-10-07T20:25:34Z
    by: agent-S-0330
  - to: in-progress
    at: 2026-10-07T20:25:34Z
    by: agent-S-0330
  - to: done
    at: 2026-10-07T20:26:04Z
    by: agent-S-0330
stream: S-0330
tags: [flai]
touches: [design/adrs]
usage:
  source: log
  seconds: 30
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 8
      output: 3555
      cache_read: 470214
      cache_write: 15415
      cost: 0.2579
---
# T-1185 An ADR records how messages between stories are addressed, kept, and closed, and why they are apart from threads

## Work

Write the ADR with `flai adr new` before any code, since every later task builds what it decides. It waits for nothing.

- **Addressing.** A conversation is between two stories, not two agent names: whichever agent works a story reads its messages, so a restart or a new session loses nothing (I-0037).
- **Storage.** Recommended: one markdown file per conversation under `wip/messages/`, with its own ID prefix, front matter naming the two stories, `about` paths, and state, and dated entries as threads have. Durable state stays in files (ADR-0020). Weigh it against threads with a `to` field, and against a log under `.flai-cache`.
- **Apart from threads.** Messages never count as awaiting the operator, so the operator's inbox stays the operator's.
- **Closing.** A conversation reads as closed once either story is not in progress or in review; acceptance and `flai archive` write the closing entry, as ADR-0109 does for threads.
- **Refusals.** A send to or from a story that is not open is refused.

## Done when

- The ADR is accepted in `design/adrs` with its decision, its alternatives, and its consequences.
- `flai check --strict` passes.

## Notes
