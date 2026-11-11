---
id: I-0345
title: wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-10T11:10:23Z
last_reported: 2026-11-10T11:10:23Z
updated: 2026-11-11T22:52:53Z
---

# I-0345 wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Description
wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Instances

### 2026-11-10T11:10:23Z
Story: S-0708.
While S-0708's close-out ran, from 2026-11-10T11:04:00Z to about 11:04:25Z, wip/kanban/tasks/T-3079-... (S-0665's new task) held only the word PLACEHOLDER. Every call S-0708's agent made to wait_for_events failed with "<path>: no front matter" until S-0665's agent finished the creation. A checked creation that leaves a placeholder in place, as ADR-0272 describes, should be one that wait_for_events, inbox, and the other readers skip, not one that fails the whole call. I-0036, closed, recorded the same failure for a read made during a save.

Correction, 2026-11-11, S-0717 (T-3120): no checked creation wrote the placeholder, and no flai code holds the string. S-0665's agent wrote it with Claude Code's `Write` tool, by mistake, and wrote the item back twenty seconds later. Its log, `~/.flai/serve/agents/sf-S-0665-20261110T102103Z.log`, shows each step:

| Time (UTC) | What S-0665's agent did |
|------------|-------------------------|
| 11:03:52.015 | `flai task new --story S-0665` creates T-3079, a whole item; the checked creation ends at 11:03:53 |
| 11:03:59.967 | `Write` of the main checkout's T-3079 file with the content `PLACEHOLDER` |
| 11:04:04.270 | Says "I overwrote T-3079's file by mistake, front matter included", and reads the item back from `.flai-cache/items-wip-kanban-tasks.json` |
| 11:04:19.632 | `Write` of the file again, with its front matter and body |

So the file did not parse from 11:03:59.967 to 11:04:19.632. Any writer can leave such a file, by hand or by tool, so the readers must not fail on it. The fault was the loader's: `Repo.List` failed the whole call on the first file that did not parse.

## Remediation

Story S-0717 remediates this issue, created from it at 2026-11-11T01:43:38Z.
Closed 2026-11-11T22:52:53Z: S-0716 fixed the cause: Repo.List skips a work item file that does not parse (711f3d9c1), so wait_for_events, inbox, and the other readers list the rest. S-0717 found the writer was S-0665's agent's Write tool call, not a checked creation (T-3120), and pinned wait_for_events in one project and across a folder of projects while a task file holds only PLACEHOLDER (T-3121: TestWaitForEventsAnswersWhileATaskFileHoldsOnlyAPlaceholder, TestWaitForEventsAcrossAFolderAnswersWhileATaskFileHoldsOnlyAPlaceholder).
