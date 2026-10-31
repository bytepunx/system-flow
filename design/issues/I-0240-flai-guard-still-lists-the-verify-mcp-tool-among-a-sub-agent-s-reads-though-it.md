---
id: I-0240
title: flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree
class: impression
status: open
count: 1
cost: 0m
first_reported: 2026-10-31T10:41:25Z
last_reported: 2026-10-31T10:41:25Z
updated: 2026-10-31T10:41:25Z
---

# I-0240 flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Description
flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Instances

### 2026-10-31T10:41:25Z
Story: S-0546.
S-0546 made the MCP verify tool ask flai serve for the close-out job (ADR-0229), which runs scripts/close-out.sh and commits in the story's worktree with FLAI_AGENT set to the caller. guard.MCPReads still lets every sub-agent and strategic agent call it, so a verifier or task sub-agent could start a close-out that commits; the delegation convention now forbids it unless the prompt asks, but the guard does not enforce it. guard.go's waitsOnSubAgents refusal text also still says wait_for_events is only for a thread awaiting the designer, not the verified event.

## Remediation
