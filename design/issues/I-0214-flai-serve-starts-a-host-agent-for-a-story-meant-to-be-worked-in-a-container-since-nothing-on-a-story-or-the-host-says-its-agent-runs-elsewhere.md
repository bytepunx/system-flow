---
id: I-0214
title: flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere
class: defect
status: open
count: 2
cost: 8m
first_reported: 2026-10-28T12:27:00Z
last_reported: 2026-10-28T12:32:02Z
updated: 2026-10-28T12:52:40Z
---

# I-0214 flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere

## Description
flai serve starts a host agent for a story meant to be worked in a container, since nothing on a story or the host says its agent runs elsewhere

## Instances

### 2026-10-28T12:27:00Z
Story: S-0505.
S-0505 T-2166: S-0513 was made to be worked by hand in flai-agent-claude-code:local, and its notes said so. agent-S-0505 moved it to ready and started it on the host with flai story start at 12:22:46Z; flai serve started agent-S-0513 on the host within a second all the same. That agent read the notes, stood down without a change, and messaged S-0505 (MS-0200). Moving the story to ready and letting the container agent call story_start was not possible either: flai serve takes a ready story first, and flai guard refuses S-0505's agent FLAI_AGENT=agent-S-0513 (ADR-0143). A story's agent harness is the only place to say where it runs, and only flai serve's adapters read it. Remediation: the container adapter (S-0508) makes the container a harness flai serve starts, so the hand-run case goes; until then a story needs a way to say flai serve must not start its agent.

### 2026-10-28T12:32:02Z
Story: S-0505.
S-0505 T-2166: the first host agent for S-0513 ended (exit 0) with the story in progress, and flai serve started a second as automatic restart 1 of 2 (ADR-0108) at about 12:26Z, while the container agent was working the story. It too stood down without a change (MS-0200). A third would have been asked of the operator; the container agent's move to review at 12:30:04Z ended the restarts.

## Remediation

Story S-0516 remediates this issue, created from it at 2026-10-28T12:52:40Z.
