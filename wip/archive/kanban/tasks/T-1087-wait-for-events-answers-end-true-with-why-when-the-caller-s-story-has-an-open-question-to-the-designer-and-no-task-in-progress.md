---
id: T-1087
type: task
nature: improvement
title: "wait_for_events answers end: true with why when the caller's story has an open question to the designer and no task in progress"
status: done
parent: S-0272
owner: alex
created: 2026-10-06T22:52:52Z
updated: 2026-10-07T00:48:59Z
transitions:
  - to: ready
    at: 2026-10-07T00:42:40Z
    by: agent-S-0272
  - to: in-progress
    at: 2026-10-07T00:42:40Z
    by: agent-S-0272
  - to: done
    at: 2026-10-07T00:48:59Z
    by: agent-S-0272
stream: S-0272
tags: [cli, mcp]
touches: [flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go]
usage:
  source: log
  seconds: 379
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 13523
      cache_read: 1669411
      cache_write: 64072
      cost: 0.974
    - model: claude-sonnet-5-5
      input: 18
      output: 4011
      cache_read: 327010
      cache_write: 51933
      cost: 0.2354
---
# T-1087 wait_for_events answers end: true with why when the caller's story has an open question to the designer and no task in progress

## Work

Criterion 1, the MCP side. In `flai/internal/mcpserver/server.go`, `waitForEvents` (about line 867) and `WaitOut` (about line 801):

- Add `end` (bool) and `why` (string, omitted when empty) to `WaitOut`.
- After the catch-up, so that events already behind the cursor, the answer among them, are still reported first, check the caller's story. It is the one in `FLAI_STORY`, which flai serve sets for a story's agent (`serve/agents.go`, about line 739). Answer at once with `end: true` and a `why` naming the thread when two things hold. The story has a thread open to the designer, one that awaits someone other than this agent. No task of the story is in progress. Then hold for nothing.
- Keep the timeout as it is in every other case. That covers an agent with no `FLAI_STORY`, such as one run by hand, which flai serve would not start again. It also covers the planner, the orchestrator and the analyzer (`FLAI_ROLE`), and a story with a task in progress.
- Say the same in the tool's description and in the server's instructions (about lines 175 and 180).

Tests go in `server_test.go`, beside `TestWaitForEvents*`:

- end on an open question with no task in progress
- no end with a task in progress
- no end with no `FLAI_STORY`
- no end for a strategic role
- events already behind the cursor come before `end`

It waits for no task: it shares no path with the others in layer 1.

## Done when

- `wait_for_events` answers `end: true` and `why` in the case above at once, and holds to its timeout otherwise.
- The tool's description and the server's instructions say when it ends.
- The new tests pass with `scripts/flai-test.sh`.

## Notes

Drafted by the planner. `flai/internal/guard`'s refusal of `wait_for_events` while a sub-agent runs (S-0285) is unchanged.
