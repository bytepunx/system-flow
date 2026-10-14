---
id: I-0097
title: A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself
class: efficiency
status: closed
count: 1
cost: 3m
first_reported: 2026-10-06T23:28:35Z
last_reported: 2026-10-06T23:28:35Z
updated: 2026-10-14T21:22:40Z
---

# I-0097 A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Description
A story's agent that runs flai stream open is left with the story in ready, so permission_prompt refuses its .claude/ write until it moves the story itself

## Instances

### 2026-10-06T23:28:35Z
Story: S-0300.
flai serve started agent-S-0300 on a ready story; the start prompt says to run flai stream open, which opened the branch and the narrative but left S-0300 in ready. The Write of .claude/agents/planner.md was refused: 'S-0300 is ready, not in progress'. Moving S-0300 to in-progress by hand cleared it. S-0274 (flai story start) would fold the move into opening.

## Remediation

Story S-0304 remediates this issue, created from it at 2026-10-07T01:07:09Z.
Closed 2026-10-14T21:22:40Z: S-0304 fixed it. flai stream open on a ready story now starts it with storystart.Begin, the code flai story start uses: it moves the story to in-progress as FLAI_AGENT, its epic following, or refuses it with exit 4 when the board holds it (ADR-0153). An agent that opens a ready story's stream is never left with it in ready, so permission_prompt asks for its protected write as in any in-progress story. TestStreamOpenStartsAReadyStory reproduces the cause.
