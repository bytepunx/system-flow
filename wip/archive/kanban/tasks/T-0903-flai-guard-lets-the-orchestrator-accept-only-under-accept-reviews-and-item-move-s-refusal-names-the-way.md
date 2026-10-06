---
id: T-0903
type: task
nature: feature
title: flai guard lets the orchestrator accept only under accept_reviews, and item_move's refusal names the way
status: done
parent: S-0221
owner: alex
created: 2026-10-05T04:47:37Z
updated: 2026-10-06T11:24:16Z
transitions:
  - to: ready
    at: 2026-10-06T11:18:32Z
    by: agent-S-0221
  - to: in-progress
    at: 2026-10-06T11:18:33Z
    by: agent-S-0221
  - to: done
    at: 2026-10-06T11:24:16Z
    by: agent-S-0221
stream: S-0221
tags: [flai]
touches: [flai/internal/guard, flai/cmd/guard.go, flai/internal/mcpserver]
after: [T-0898]
usage:
  source: log
  seconds: 343
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 57
      output: 19172
      cache_read: 2943395
      cache_write: 82285
      cost: 1.455
---
# T-0903 flai guard lets the orchestrator accept only under accept_reviews, and item_move's refusal names the way

## Work

S-0218 gives `flai guard` the orchestrator's permission table. To it, add the rule for acceptance:

- With `accept_reviews` on, the orchestrator may run `flai accept <S-nnnn>` (and `flai move <S-nnnn> done`, which runs the same flow), but only with `--by orchestrator`.
- With it off, or with any other `--by`, the call is refused. The refusal names `orchestration.permissions.accept_reviews`, and S-0218's log of refusals records it.
- A story agent, a sub-agent, and the planner are refused as they are today.

In `flai/internal/mcpserver/server.go`, `item_move` still refuses to move a story to done. The orchestrator accepts through `flai accept`, which runs the preview and records the evidence, so `item_move` stays the operator's no. When the caller is the orchestrator, the refusal's message names `flai accept --by orchestrator` and the permission.

Add guard tests beside the existing ones: a permitted accept, a refusal with the permission off, a refusal with another `--by`, and a refusal for the planner.

This task waits for T-0898, which fixes the rule.

## Done when

- `flai guard` allows the orchestrator's `flai accept --by orchestrator` only with `accept_reviews` on, and refuses everything else naming the permission.
- `item_move` to done by the orchestrator is refused with a message naming `flai accept` and the permission.
- The guard tests and `scripts/flai-test.sh` pass.

## Notes
