---
id: I-0167
title: task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
class: efficiency
status: open
count: 2
cost: 4m
first_reported: 2026-10-19T15:28:45Z
last_reported: 2026-10-19T20:17:43Z
updated: 2026-10-19T20:17:43Z
---

# I-0167 task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator

## Description
task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator

## Instances

### 2026-10-19T15:28:45Z
Story: S-0398.
S-0398: story_start answered 93,426 characters and each of five task_done calls 83,000 to 84,000. Claude Code saved each one to a file instead of returning it, so every close needed a jq call to find stopped, error, left, and the sync's conflicts. Most of the bulk is the inbox's threads: 33 awaiting the operator, none on S-0398. The pack's parts are capped at 40,000 bytes of JSON (ADR-0104), but these answers have no cap.

### 2026-10-19T20:17:43Z
Story: S-0399.
story_start S-0399 answered 101,139 characters, 32 threads awaiting the operator in its inbox; Claude Code saved it to a file and the agent read the story, worktree, and pack header out of it with jq

## Remediation
