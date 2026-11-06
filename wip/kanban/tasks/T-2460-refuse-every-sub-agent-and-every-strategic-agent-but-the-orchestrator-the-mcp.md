---
id: T-2460
type: task
nature: remediation
title: Refuse every sub-agent, and every strategic agent but the orchestrator, the MCP verify tool and flai verify --job
status: backlog
parent: S-0556
owner: alex
created: 2026-11-05T01:37:13Z
updated: 2026-11-05T01:37:13Z
transitions: []
stream: S-0556
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go]
after: [T-2459]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 6
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 7
          cache_read: 65969
          cache_write: 1221
          cost: 0.016
---
# T-2460 Refuse every sub-agent, and every strategic agent but the orchestrator, the MCP verify tool and flai verify --job

## Work

This is the fix I-0240 asks for, in `flai/internal/guard/guard.go`:

- remove `verify` from `MCPReads`, and rewrite its comment: the S-0270 reason, that `verify` only stores a result in flai's cache, stopped being true when S-0546 made it ask flai serve for a close-out job that commits (ADR-0229)
- add `--job` to the `without` flags of `verify` in `flagReads`, so that a sub-agent's `flai verify S-nnnn --job` is refused as `--record-issues` is, while `flai verify S-nnnn --last` and a plain `flai verify S-nnnn` still pass; rewrite the comment above `flagReads` to say why
- reword `waitsOnSubAgents`: `wait_for_events` is for the `verified` event of the story's close-out job too, not only for a thread awaiting the designer, as `design/conventions/delegation.md` already says

Tests in `flai/internal/guard/guard_test.go`, each reproducing the cause first:

- turn the test that lets a verifier call MCP `verify` around: a verifier's, an explorer's, and a task sub-agent's MCP `verify` call is refused, naming the tool
- a sub-agent's `flai verify S-0001 --job` and `flai verify S-0001 --job -m x` are refused, and `flai verify S-0001 --last` still passes
- the planner's, the analyzer's, and the curator's own MCP `verify` call is refused, and the story's agent's own call passes

Check the tests that loop over `MCPReads` (`guard_test.go`, `curate_test.go`) still hold without `verify`. This task waits for T-2459, because the orchestrator must hold `verify` by name before it leaves `MCPReads`, and both change `guard.go`.

## Done when

- `verify` is not in `guard.MCPReads`, `flai verify --job` is refused to every sub-agent, and `waitsOnSubAgents` names the `verified` event
- The new tests fail without the change and pass with it, and `flai test` on `flai/internal/guard` passes

## Notes
