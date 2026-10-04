---
id: I-0069
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
class: blocker
status: open
count: 3
cost: 10m
first_reported: 2026-10-04T00:50:33Z
last_reported: 2026-10-04T03:17:17Z
updated: 2026-10-04T03:17:17Z
---

# I-0069 A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it

## Description
A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it

## Instances

### 2026-10-04T00:50:33Z
Story: S-0208.
S-0208 T-0785 needed the PreToolUse matcher widened to Edit|Write|NotebookEdit for the planner; the task sub-agent's and the story agent's edits of .claude/settings.json and template/root/.claude/settings.json were refused by the session's permissions. Worked around with a planner-only --settings hook (TH-0096).

### 2026-10-04T00:52:16Z
Story: S-0208.
S-0208 T-0784: creating template/root/.claude/agents/planner.md and .claude/agents/planner.md was refused as a sensitive file for the task sub-agent and the story agent alike; asked the operator on TH-0097.

### 2026-10-04T03:17:17Z
Story: S-0208.
the settings files too: after the operator granted it on TH-0096, Edit of .claude/settings.json and template/root/.claude/settings.json was still refused; asked the operator to write them

## Remediation
