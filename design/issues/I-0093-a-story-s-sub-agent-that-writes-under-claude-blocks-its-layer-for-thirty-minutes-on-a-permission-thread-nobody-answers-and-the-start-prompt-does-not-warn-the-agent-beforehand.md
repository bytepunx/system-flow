---
id: I-0093
title: A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand
class: blocker
status: closed
count: 1
cost: 6m
first_reported: 2026-10-06T20:53:14Z
last_reported: 2026-10-06T20:53:14Z
updated: 2026-10-06T22:53:01Z
---

# I-0093 A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Description
A story's sub-agent that writes under .claude/ blocks its layer for thirty minutes on a permission thread nobody answers, and the start prompt does not warn the agent beforehand

## Instances

### 2026-10-06T20:53:14Z
Story: S-0223.
T-0957's sub-agent called Write on template/root/.claude/agents/analyzer.md and Edit on template/root/.claude/settings.json. The installed flai's permission_prompt opened TH-0192 and held the call. The owner was away, so the layer waited until the board watcher stopped and restarted the session about six minutes later. The warning on TH-0190 reached the agent after it had launched the layer. The fix was to write the files into .flai-cache/S-0223/ and ask alex to cp them in on TH-0196.

## Remediation

Story S-0299 remediates this issue, created from it at 2026-10-06T21:00:33Z.
Closed 2026-10-06T22:53:01Z: S-0299: flai guard refuses a story's sub-agent an Edit, MultiEdit, Write, or NotebookEdit of a file in a .claude/ folder at once while auto-approve is off (T-1035), the settings now run the guard before a story session's file writes (T-1053), and the story agent's start prompt says before any layer that such a write is its own, made after the layer or handed over with cp commands on a thread (T-1036); ADR-0102.
