---
id: I-0103
title: A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout
class: efficiency
status: closed
count: 1
cost: 30m
first_reported: 2026-10-07T04:55:43Z
last_reported: 2026-10-07T04:55:43Z
updated: 2026-10-08T06:02:09Z
---

# I-0103 A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout

## Description
A story agent's .claude/ write waits thirty minutes on an unanswered permission thread, then fails on Claude Code's MCP idle timeout

## Instances

### 2026-10-07T04:55:43Z
Story: S-0270.
S-0270's agent wrote `.claude/agents/verifier.md` through `permission_prompt`, which was allowed. Its next write, `template/root/.claude/agents/verifier.md`, opened a permission thread that nobody answered. After 1800 s Claude Code aborted the call ("MCP server flai tool permission_prompt sent no response or progress for 1800s"), and the write failed. The agent then staged the remaining three files in the worktree's `.flai-cache/` and opened TH-0245 with cp commands. permission_prompt sends no progress while it waits, so Claude Code's idle timeout ends it at thirty minutes, and the agent cannot tell beforehand whether the operator is there.

## Remediation

Story S-0309 remediates this issue, created from it at 2026-10-07T06:48:44Z.
Closed 2026-10-08T06:02:09Z: Fixed by S-0309 (ADR-0124). permission_prompt holds a protected write at most four minutes, below Claude Code's MCP idle timeout. When no answer comes in that time, or the session ends, it refuses the write naming the thread and leaves the thread open. When the agent makes the same write again, the answer given since is taken from that thread, or it waits again. The start prompt and delegation.md now tell the agent to retry the write once the thread is answered, instead of staging files in .flai-cache/ with cp commands. TestPermissionPromptRefusesUnansweredWithinTheBoundAndLeavesTheThreadOpen reproduces the cause.
