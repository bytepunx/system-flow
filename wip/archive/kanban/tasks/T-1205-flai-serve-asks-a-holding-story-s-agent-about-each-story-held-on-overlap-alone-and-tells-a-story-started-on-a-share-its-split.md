---
id: T-1205
type: task
nature: feature
title: flai serve asks a holding story's agent about each story held on overlap alone, and tells a story started on a share its split
status: done
parent: S-0334
owner: alex
created: 2026-10-07T20:16:17Z
updated: 2026-10-08T10:20:15Z
transitions:
  - to: ready
    at: 2026-10-08T10:11:03Z
    by: agent-S-0334
  - to: in-progress
    at: 2026-10-08T10:11:04Z
    by: agent-S-0334
  - to: done
    at: 2026-10-08T10:20:04Z
    by: agent-S-0334
stream: S-0334
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/start.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/workitem/share.go]
after: [T-1203]
usage:
  source: log
  seconds: 540
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 59
      output: 25687
      cache_read: 4207791
      cache_write: 124340
      cost: 2.106
---
# T-1205 flai serve asks a holding story's agent about each story held on overlap alone, and tells a story started on a share its split

## Work

Send the ask and pass the share on. It waits for T-1203, whose `AskHold` and shares it calls. It shares no path with T-1206 and runs beside it.

- At each look, for each ready story held on overlap alone, call `AskHold` for each story in progress that holds it (`Holds.OverlapsBy`); it asks once per pair while the held story stays in ready.
- When flai serve starts a story whose overlap a share cleared, the prompt `harness.Prompt` builds names the conversation, the shared paths, and the split.
- A request that cannot be written is logged and the look goes on.

## Done when

- Tests cover the request sent once per pair, no request for a story held by `after` or `no-touches`, and the prompt naming the share.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
