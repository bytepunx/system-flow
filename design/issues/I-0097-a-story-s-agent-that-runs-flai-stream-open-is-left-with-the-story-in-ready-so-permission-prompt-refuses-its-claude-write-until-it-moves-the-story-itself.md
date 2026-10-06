---
id: I-0097
title: A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
class: efficiency
status: open
count: 1
cost: 3m
first_reported: 2026-10-06T23:28:35Z
last_reported: 2026-10-06T23:28:35Z
updated: 2026-10-06T23:28:35Z
---

# I-0097 A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Description
A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Instances

### 2026-10-06T23:28:35Z
Story: S-0300.
flai serve started agent-S-0300 on a ready story; the start prompt says to run flai stream open, which opened the branch and the narrative but left S-0300 in ready. The Write of .claude/agents/planner.md was refused: 'S-0300 is ready, not in progress'. Moving S-0300 to in-progress by hand cleared it. S-0274 (flai story start) would fold the move into opening.

## Remediation
