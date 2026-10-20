---
id: I-0167
title: task_done and story_start answers exceed Claude Code's tool output limit, because the inbox they carry lists every thread awaiting the operator
class: efficiency
status: closed
count: 4
cost: 3m
first_reported: 2026-10-19T15:28:45Z
last_reported: 2026-10-20T00:35:02Z
updated: 2026-10-20T19:54:47Z
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

### 2026-10-19T22:24:12Z
Story: S-0401.
S-0401: story_start answered 87,568 characters, past Claude Code's limit; the pack was only its header and the inbox carried 31 threads awaiting the operator, so the answer was read from the saved file with jq.

### 2026-10-20T00:35:02Z
Story: S-0404.
S-0404: story_start's answer was 87,106 characters with 32 threads awaiting in its inbox and the pack only as a header, so Claude Code saved it to a file and the agent read it back with jq.

## Remediation

Story S-0450 remediates this issue, created from it at 2026-10-20T18:53:36Z.
Closed 2026-10-20T19:54:47Z: Fixed by S-0450 (ADR-0179): story_start and task_done, over MCP and with --json, now carry an inbox that lists only the threads on the story and its tasks and counts the rest in threads_omitted, and each whole answer is fitted to 40,000 bytes of JSON, cutting the oldest changes, then the pack to its header, then the listed threads' last entries. Tests reproduce S-0398's instance (33 threads elsewhere, task_done) and S-0399's (32 threads elsewhere, story_start). S-0450's own task_done answered 54,219 bytes before and 2,873 after, with 33 threads omitted.
