---
id: T-2297
type: task
nature: feature
title: inbox.CatchUp reports a thread event when a thread comes to await the reader, or for the orchestrator the operator, naming the thread, its item, and whom it awaits under ADR-0181
status: done
parent: S-0535
owner: alex
created: 2026-10-28T21:38:39Z
updated: 2026-10-28T23:23:41Z
transitions:
  - to: ready
    at: 2026-10-28T22:45:33Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T22:45:34Z
    by: system-flow
  - to: done
    at: 2026-10-28T23:23:41Z
    by: agent-S-0535
stream: S-0535
tags: [flai, events, inbox, threads]
touches: [flai/internal/inbox/threadevents.go, flai/internal/inbox/threadevents_test.go, flai/internal/inbox/inbox.go, flai/internal/inbox/inbox_test.go, flai/internal/mcpserver/cursor.go, flai/internal/mcpserver/server_test.go]
after: [T-2292]
usage:
  source: log
  seconds: 2287
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 84
      output: 28542
      cache_read: 5340332
      cache_write: 236451
      cost: 2.969
---
# T-2297 inbox.CatchUp reports a thread event when a thread comes to await the reader, or for the orchestrator the operator, naming the thread, its item, and whom it awaits under ADR-0181

## Work

`inbox.Summarize` (`flai/internal/inbox/inbox.go`) works out whom an open thread awaits under ADR-0181: `you`, `operator`, or `other`. Nothing records the moment a thread comes to await someone. Threads are written by the MCP server, the CLI, the dashboard, and flai serve, so a log written at each writer would miss some. Derive the event from the thread files instead, as `CatchUp` derives messages.

- Add `flai/internal/inbox/threadevents.go`: for each open thread with an entry newer than the reader's cursor, work out what it awaits with `Summarize`'s rule. Report a `thread` event when it awaits the reader (`you`), or, in the orchestrator's role, the operator. Give the thread ID, its item or path, `to` the awaiting value, whom it awaits by name, and the last entry's author and time.
- Report a thread once per new entry that makes it await the reader. Report nothing for an entry by the reader, a resolved thread, or a thread that awaits another.
- Merge the events into `CatchUp` under `MaxEvents`, and give `Describe` a case.

It waits for T-2292, which adds the other kinds to `inbox.go` first. It shares no path with T-2294, T-2295, or T-2296, so the four run together.

## Done when

- Tests show one `thread` event for each of these: a reply on a thread the reader takes part in; a new thread on the reader's own story; and, read in the orchestrator's role, a planner's question that awaits the operator.
- Tests show no event for the reader's own entry, a resolved thread, a thread awaiting another agent, or a second read after the cursor moves.
- `flai test` on the changed paths passes.

## Notes

T-2293's `awaits-orchestrator` matches this event when it awaits the orchestrator or the operator.
