---
id: T-1840
type: task
nature: feature
title: flai thread reply and flai thread confirm write the operator's answer as a digest entry and name it in their result
status: done
parent: S-0433
owner: alex
created: 2026-10-16T06:07:52Z
updated: 2026-10-24T20:58:12Z
transitions:
  - to: ready
    at: 2026-10-24T20:35:52Z
    by: agent-S-0433
  - to: in-progress
    at: 2026-10-24T20:35:52Z
    by: agent-S-0433
  - to: done
    at: 2026-10-24T20:58:12Z
    by: agent-S-0433
stream: S-0433
tags: [cli]
touches: [flai/cmd/thread.go, flai/cmd/thread_test.go, docs/users/flai-reference.md]
after: [T-1837]
usage:
  source: log
  seconds: 1340
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 10925
      cache_read: 1688396
      cache_write: 93282
      cost: 1.0998
---
# T-1840 flai thread reply and flai thread confirm write the operator's answer as a digest entry and name it in their result

## Work

In `flai/cmd/thread.go`, after `reply` and `confirm` save the thread, call `digest.OnAnswer` and print what it wrote or updated, or its warning on standard error; `--json` adds `digest` with the entry's ID and whether it was written or updated.

Waits for T-1837. Runs beside the MCP and host API tasks, which touch other files.

## Done when

`thread_test.go` covers an operator's reply that writes an entry, a confirmation, an agent's reply that writes none, and a warning that leaves the reply saved; `flai test flai/cmd` passes.

## Notes

The author is `threadAuthor`'s; the operator is the manifest's `owner` (TH-0476 question 5).
