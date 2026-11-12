---
id: T-3139
type: task
nature: remediation
title: Say in the operators guide that the orchestrator in a container verifies and accepts through MCP on the host
status: done
parent: S-0725
owner: alex
created: 2026-11-12T15:51:43Z
updated: 2026-11-12T15:56:24Z
transitions:
  - to: ready
    at: 2026-11-12T15:55:27Z
    by: agent-S-0725
  - to: in-progress
    at: 2026-11-12T15:55:28Z
    by: agent-S-0725
  - to: done
    at: 2026-11-12T15:56:24Z
    by: agent-S-0725
stream: S-0725
tags: [docs]
touches: [docs/operators/index.md]
usage:
  source: log
  seconds: 56
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 1470
      cache_read: 972437
      cache_write: 28350
      cost: 0.4507
---
# T-3139 Say in the operators guide that the orchestrator in a container verifies and accepts through MCP on the host

## Work

It waits for no task: it changes `docs/operators/index.md` alone, which no other task of this story touches.

In § Strategic agents in a container (S-0509), two paragraphs contradict [ADR-0278](../../design/adrs/0278-the-orchestrator-accepts-through-an-mcp-tool-that-runs-flai-accept-on-the-host.md):

- **Where a story's code runs.** It says the orchestrator's own `flai verify` and `flai accept`, run in its shell, run a story's tiers and an acceptance's `fix_command` inside its container. Since ADR-0278 they refuse there, naming the MCP tools `accept` and `verify`: the orchestrator verifies with `verify`, which runs the close-out job in `flai serve` on the host, and accepts with `accept`, which runs `flai accept --by orchestrator` on the host with the host's configuration, so an acceptance's `fix_command` runs on the host too. Rewrite the paragraph to say so, keeping its sentences on the close-out by the story's own harness (ADR-0274) and the merge queue's batch.
- **What it cannot do yet.** It says `flai accept` and `flai verify` in the container see none of the host actions. Say instead that they refuse in the container and what the orchestrator uses in their place, and keep the sentence on a defined role named `plan`, `orchestrate`, `analyze`, or `curate`. Check the claim that the image lacks the test tiers' tools against `flai/images/system-flow/Dockerfile` and the role container's image in `flai/internal/harness/container.go` before keeping it.

Follow `design/system/strategic-agents.md` § In a container (S-0509), which ADR-0278 already updated, and link to it rather than repeat it. Search the rest of `docs/operators/index.md` for another statement that the orchestrator's verify or accept runs in its container, and correct it the same way.

## Done when

- § Strategic agents in a container no longer says the orchestrator's shell `flai verify` or `flai accept` runs tiers or a `fix_command` in its container, and names the MCP tools `verify` and `accept` with ADR-0278.
- `flai test docs/operators/index.md` passes, the markdown lint among its tiers.

## Notes
