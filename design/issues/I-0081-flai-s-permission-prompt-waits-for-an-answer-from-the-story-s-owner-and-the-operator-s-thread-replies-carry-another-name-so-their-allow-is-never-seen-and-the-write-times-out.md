---
id: I-0081
title: flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out
class: defect
status: open
count: 2
cost: 45m
first_reported: 2026-10-05T07:45:13Z
last_reported: 2026-10-06T12:38:18Z
updated: 2026-10-06T12:38:18Z
---

# I-0081 flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

## Description
flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

## Instances

### 2026-10-05T07:45:13Z
Story: S-0218.
S-0218's owner is arobson; the operator replied allow on TH-0158 as alex (system-flow.yaml's owner) at 07:21Z; awaitAnswer (flai/internal/mcpserver/permission.go:247) accepts only the owner's entry, so the Write of template/root/.claude/agents/orchestrator.md waited until Claude Code's 1800s MCP idle timeout aborted it

### 2026-10-06T12:38:18Z
Story: S-0222.
T-0895's sub-agent Edit on .claude/agents/orchestrator.md and then template/root/.claude/agents/orchestrator.md each held thirty minutes in the installed flai's permission_prompt: the operator replied allow as alex, the prompt took only arobson; the operator restarted the agent, and the definition waits for the operator to copy it in (TH-0183, TH-0185)

## Remediation

Story S-0284 remediates this issue, created from it at 2026-10-06T03:45:19Z.
