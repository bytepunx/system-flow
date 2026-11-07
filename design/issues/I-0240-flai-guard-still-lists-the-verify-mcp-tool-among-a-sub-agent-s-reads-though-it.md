---
id: I-0240
title: flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree
class: impression
status: closed
count: 1
cost: 0m
first_reported: 2026-10-31T10:41:25Z
last_reported: 2026-10-31T10:41:25Z
updated: 2026-11-07T06:50:11Z
---

# I-0240 flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Description
flai guard still lists the verify MCP tool among a sub-agent's reads, though it now asks flai serve for a close-out job that commits in the story's worktree

## Instances

### 2026-10-31T10:41:25Z
Story: S-0546.
S-0546 made the MCP verify tool ask flai serve for the close-out job (ADR-0229), which runs scripts/close-out.sh and commits in the story's worktree with FLAI_AGENT set to the caller. guard.MCPReads still lets every sub-agent and strategic agent call it, so a verifier or task sub-agent could start a close-out that commits; the delegation convention now forbids it unless the prompt asks, but the guard does not enforce it. guard.go's waitsOnSubAgents refusal text also still says wait_for_events is only for a thread awaiting the designer, not the verified event.

## Remediation

Story S-0556 remediates this issue, created from it at 2026-11-01T10:44:38Z.
Closed 2026-11-07T06:50:11Z: S-0556: flai guard no longer counts the MCP tool verify among guard.MCPReads, and refuses every sub-agent, and every strategic agent but the orchestrator, that tool and flai verify --job, since the close-out job commits in the story's worktree. The story's agent, whose own calls pass, asks for the job, and the orchestrator holds the tool by name in its role file and guard.MCPOrchestrates. flai verify --last and plain flai verify still pass, and waitsOnSubAgents names the verified event. TestASubAgentRunsFlaiVerifyButRecordsNoIssues and TestOnlyTheStorysAgentAndTheOrchestratorAskForTheCloseOutJob reproduce the cause, and TestTheOrchestratorHoldsVerifyByName holds the orchestrator's call.
