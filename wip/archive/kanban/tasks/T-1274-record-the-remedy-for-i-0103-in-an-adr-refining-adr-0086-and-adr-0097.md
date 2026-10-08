---
id: T-1274
type: task
nature: improvement
title: Record the remedy for I-0103 in an ADR refining ADR-0086 and ADR-0097
status: done
parent: S-0309
owner: alex
created: 2026-10-07T23:27:16Z
updated: 2026-10-08T05:53:09Z
transitions:
  - to: ready
    at: 2026-10-08T05:52:46Z
    by: agent-S-0309
  - to: in-progress
    at: 2026-10-08T05:52:46Z
    by: agent-S-0309
  - to: done
    at: 2026-10-08T05:53:09Z
    by: agent-S-0309
stream: S-0309
tags: [adr, permission-prompt]
touches: [design/adrs, design/adrs/README.md]
usage:
  source: log
  seconds: 23
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1821
      cache_read: 257558
      cache_write: 11765
      cost: 0.1648
---
# T-1274 Record the remedy for I-0103 in an ADR refining ADR-0086 and ADR-0097

## Work

The story asks for a remedy proposed from I-0103's instance before it is built. The plan's thread on S-0309 proposes one; confirm it is settled there, or put it to the operator on that thread, before writing the ADR. This task waits for nothing: it is the first layer, and every other task builds what it decides.

The proposed remedy, which this ADR records once settled:

1. `permission_prompt` holds a call for at most a few minutes (proposed: 4, below Claude Code's idle timeout of 5 minutes over HTTP and 30 over stdio), so no transport drops the call.
2. When that time passes, or the session ends, unanswered, it refuses the write with a reason naming the thread, and leaves the thread open instead of settling it as refused.
3. A later request for the same story, tool, path, and input finds that thread: an allow already on it lets the write through at once and settles the thread; a refusal refuses it; no answer yet holds the call again, bounded the same way. A request with other input opens a thread of its own.
4. The open thread is the agent's own question on its story, so `wait_for_events` already ends a `flai serve` agent with `end: true` when nothing else is left (`endWhy` in `flai/internal/mcpserver/server.go`), and flai serve starts it again on the answer, when it repeats the write.

Write the ADR under `design/adrs/` with the next free number, `refines: [ADR-0086, ADR-0097]`, its Alternatives considered (progress notifications through `keepAlive` in `flai/internal/mcpserver/timing.go`, which hold the agent blocked for as long as the operator is away and help only when Claude Code sends a progress token; a manifest setting for the bound), and add it to `design/adrs/README.md`.

## Done when

- The remedy is settled on S-0309's plan thread or a thread of its own, and the ADR records it with status accepted.
- `design/adrs/README.md` lists the new ADR.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

The touch `design/adrs/` is a folder because the ADR's number and slug are not known until it is written.
