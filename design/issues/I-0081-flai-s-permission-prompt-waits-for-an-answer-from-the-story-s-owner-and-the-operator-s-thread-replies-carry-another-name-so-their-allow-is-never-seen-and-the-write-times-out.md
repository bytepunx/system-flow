---
id: I-0081
title: flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out
class: defect
status: open
count: 1
cost: 30m
first_reported: 2026-10-05T07:45:13Z
last_reported: 2026-10-05T07:45:13Z
updated: 2026-10-05T07:45:13Z
---

# I-0081 flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

## Description
flai's permission_prompt waits for an answer from the story's owner, and the operator's thread replies carry another name, so their allow is never seen and the write times out

## Instances

### 2026-10-05T07:45:13Z
Story: S-0218.
S-0218's owner is arobson; the operator replied allow on TH-0158 as alex (system-flow.yaml's owner) at 07:21Z; awaitAnswer (flai/internal/mcpserver/permission.go:247) accepts only the owner's entry, so the Write of template/root/.claude/agents/orchestrator.md waited until Claude Code's 1800s MCP idle timeout aborted it

## Remediation
