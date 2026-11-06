---
id: D-0446
title: "S-0546 kept the MCP verify tool among flai guard's reads, since a call at a verified commit answers from the record and a story in review has nothing to commit; I-0240 records that a sub-agent can still start a committing close-out"
kind: decision
sources: [S-0546, I-0240]
topics: [mcp, orchestration, flai, testing]
paths: [flai/internal/guard, flai/internal/mcpserver/verify.go, design/conventions/delegation.md]
date: 2026-11-01
status: active
absorbed_by: ""
asked_again: []
cited: [S-0556, T-2459, TH-0827]
---

# D-0446 S-0546 kept the MCP verify tool among flai guard's reads, since a call at a verified commit answers from the record and a story in review has nothing to commit; I-0240 records that a sub-agent can still start a committing close-out

S-0546 made the MCP tool `verify` ask flai serve for the close-out job (ADR-0229), which runs `scripts/close-out.sh` and commits in the story's worktree. It left the tool in `flai guard`'s list of reads, which every sub-agent and strategic agent may call. Its narrative gives the reason (wip/archive/agents/S-0546.md, Decisions):

> The MCP `verify` tool stays among the guard's reads: a sub-agent or the orchestrator asking at a verified commit gets the record, and a story in review has nothing to commit.

The story's agent also filed I-0240 against that choice: "guard.MCPReads still lets every sub-agent and strategic agent call it, so a verifier or task sub-agent could start a close-out that commits; the delegation convention now forbids it unless the prompt asks, but the guard does not enforce it." The same issue notes that the refusal text of `waitsOnSubAgents` in `guard.go` still says `wait_for_events` is only for a thread awaiting the designer, and leaves out the `verified` event.

S-0556 remediates I-0240. Whoever narrows the guard there should keep the two calls S-0546 kept the tool for. One is the orchestrator's verify of a story in review. The other is a reader asking at a commit the record holds, which runs nothing.

Sources: S-0546, I-0240.
