---
id: I-0069
title: A story agent started by flai serve cannot edit .claude/settings.json, so a story that changes the guard hook's matcher cannot ship it
class: blocker
status: open
count: 6
cost: 8m
first_reported: 2026-10-04T00:50:33Z
last_reported: 2026-10-05T01:18:40Z
updated: 2026-10-05T01:18:40Z
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

### 2026-10-04T04:11:00Z
Story: S-0209.
T-0795 had to change template/root/.claude/agents/planner.md and .claude/agents/planner.md; Claude Code refused the Edit as a sensitive file, so the contents went to the operator on TH-0101 to copy.

### 2026-10-04T20:07:01Z
Story: S-0210.
S-0210's agent and its task sub-agent could not edit template/root/.claude/agents/planner.md or .claude/agents/planner.md (Claude Code calls them sensitive); asked the operator to paste the file on TH-0109.

### 2026-10-05T01:18:40Z
Story: S-0266.
T-0850's rule for the verifier's definition was refused as a sensitive file, to the story's agent and to its task sub-agent, in both .claude/agents/verifier.md and the template's template/root/.claude/agents/verifier.md; asked the operator to paste the whole file on TH-0120.

## Remediation
