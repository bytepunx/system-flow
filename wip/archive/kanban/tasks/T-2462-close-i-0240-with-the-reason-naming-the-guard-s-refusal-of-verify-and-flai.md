---
id: T-2462
type: task
nature: remediation
title: Close I-0240 with the reason naming the guard's refusal of verify and flai verify --job
status: done
parent: S-0556
owner: alex
created: 2026-11-05T01:37:33Z
updated: 2026-11-07T06:50:17Z
transitions:
  - to: ready
    at: 2026-11-07T06:50:06Z
    by: agent-S-0556
  - to: in-progress
    at: 2026-11-07T06:50:08Z
    by: agent-S-0556
  - to: done
    at: 2026-11-07T06:50:17Z
    by: agent-S-0556
stream: S-0556
tags: [flai]
touches: [design/issues/I-0240-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md, design/issues/summary.md]
after: [T-2461]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1370
      cache_read: 553365
      cache_write: 10002
      cost: 0.213
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 65968
          cache_write: 1221
          cost: 0.016
---
# T-2462 Close I-0240 with the reason naming the guard's refusal of verify and flai verify --job

## Work

Close the issue with `flai issue close I-0240 --reason "<reason>"` in the story's worktree, the reason saying what fixed it: `flai guard` no longer counts the MCP tool `verify` among the reads every sub-agent and strategic agent may call, refuses a sub-agent `flai verify --job`, and lets the story's agent and the orchestrator, whose role file now names it, ask for the close-out job; and `waitsOnSubAgents` names the `verified` event. Name the tests that reproduce the cause. flai regenerates `design/issues/summary.md`. This task waits for T-2461, so that the reason describes what was built and documented.

## Done when

- I-0240 is closed, its reason naming the fix and its tests
- `flai check --strict` passes

## Notes
