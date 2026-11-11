---
id: I-0344
title: One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them
class: defect
status: closed
count: 1
cost: 3m
first_reported: 2026-11-10T11:05:08Z
last_reported: 2026-11-10T11:05:08Z
updated: 2026-11-11T15:15:11Z
---

# I-0344 One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them

## Description
One malformed work item file in the main checkout's wip makes every flai read of the project fail, a story's close-out and wait_for_events among them

## Instances

### 2026-11-10T11:05:08Z
Story: S-0667.
S-0667, at 11:04:19Z on 2026-11-10: wip/kanban/tasks/T-3079-teststreamstatewritescontextdecisionsandquestions-runs-in-parallel-without.md, a task of S-0665's agent, was an untracked 11-byte file holding `PLACEHOLDER`, with no front matter. flai's Go code has no such string, so what wrote it is not known. For the seconds it lasted, S-0667's close-out job stopped at check with "flai check could not run: …T-3079-…md: no front matter". wait_for_events failed three times with the same error, and stream_state wrote the narrative but failed to write wip/agents/index.md. Six seconds later the file held the task (1,581 bytes), and the close-out passed when asked again. One other story's file broke every agent's reads of the project, where a check could report it as one finding on that item and go on.

## Remediation

Story S-0716 remediates this issue, created from it at 2026-11-11T01:40:46Z.
Closed 2026-11-11T15:15:11Z: S-0716: Repo.List leaves out an item file that cannot be read or does not parse and lists the rest, and Repo.Unreadable names it, so flai check, the close-out, task done, inbox, wait_for_events, wait_for_work, and stream_state's index write go on past it; flai check reports the file once as item.unreadable, a note outside the story whose file it is. Tests reproduce it with a task file holding PLACEHOLDER.
