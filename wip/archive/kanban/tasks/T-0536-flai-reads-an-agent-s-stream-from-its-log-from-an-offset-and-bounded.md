---
id: T-0536
type: task
nature: feature
title: flai reads an agent's stream from its log, from an offset and bounded
status: done
parent: S-0142
owner: alex
created: 2026-09-29T05:44:46Z
updated: 2026-09-29T05:47:40Z
transitions:
  - to: ready
    at: 2026-09-29T05:44:51Z
    by: agent-S-0142
  - to: in-progress
    at: 2026-09-29T05:44:51Z
    by: agent-S-0142
  - to: done
    at: 2026-09-29T05:47:40Z
    by: agent-S-0142
stream: S-0142
tags: []
touches: [flai/internal/serve]
---
# T-0536 flai reads an agent's stream from its log, from an offset and bounded

## Work

Add `serve.Stream` in `flai/internal/serve`: for a story's newest run, read its log from a byte offset (the tail when none is given), take complete lines only, turn Claude Code's `stream-json` events into short entries (text, thinking, tool call, tool result, result, system), keep other output as plain lines, and bound what one read returns (bytes read, entries, characters per entry). It says whether the agent still runs and the offset to ask from next.

## Done when

- Behaviour tests cover a fresh read, a read from an offset, a partial last line, a log larger than the bound, non-JSON output, an unknown story, and a run with no log.
- `make test` and lint pass.

## Notes
