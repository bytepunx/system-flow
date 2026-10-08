---
id: T-1211
type: task
nature: feature
title: The host API reads messages.list and messages.get serve the conversations of a project or a story
status: done
parent: S-0336
owner: alex
created: 2026-10-07T20:16:52Z
updated: 2026-10-08T06:17:27Z
transitions:
  - to: ready
    at: 2026-10-08T06:14:00Z
    by: agent-S-0336
  - to: in-progress
    at: 2026-10-08T06:14:01Z
    by: agent-S-0336
  - to: done
    at: 2026-10-08T06:17:27Z
    by: agent-S-0336
stream: S-0336
tags: [flai]
touches: [flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 206
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 42
      output: 247
      cache_read: 1427881
      cache_write: 77528
      cost: 0.6716
---
# T-1211 The host API reads messages.list and messages.get serve the conversations of a project or a story

## Work

Serve messages to the dashboard, read-only. It waits for nothing in this story.

- `messages.list` with an optional `story` and `all`, and `messages.get` with an `id`, beside `threads.list` in `hostapi.go`: each conversation with its two stories, `about` paths, state, which side it awaits, and its entries. Shape each with `messages.View`, so the reads carry whatever a conversation holds, a share included once S-0334 adds one.
- No write: messages are the agents'.
- Add both to the contract test's list of reads.

## Done when

- Tests cover both reads, the story filter, closed conversations under `all`, and the contract.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
