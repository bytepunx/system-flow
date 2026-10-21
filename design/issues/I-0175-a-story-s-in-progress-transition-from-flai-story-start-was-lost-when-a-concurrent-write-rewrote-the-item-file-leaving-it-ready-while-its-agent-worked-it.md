---
id: I-0175
title: A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-21T07:37:23Z
last_reported: 2026-10-21T07:37:23Z
updated: 2026-10-21T11:59:43Z
---

# I-0175 A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it

## Description
A story's in-progress transition from flai story start was lost when a concurrent write rewrote the item file, leaving it ready while its agent worked it

## Instances

### 2026-10-21T07:37:23Z
Story: S-0463.
story_start moved S-0463 to in-progress at about 07:33:32Z (its narrative opened at 07:33:39Z) and answered status in-progress. At 07:36 the item file in the main checkout said status ready, with only the orchestrator's 07:33:31Z ready transition, and the inbox listed S-0463 as ready to pull while its agent had closed T-1977. .flai-cache/edits.jsonl records a forecast write by flai at 07:33:45Z and the file's updated stamp says 07:36:36Z (a usage write); one of them, or both, read the file before the move and wrote it back after. The agent moved S-0463 to in-progress again at 07:37Z, so the transition's time is about four minutes late and another agent could have pulled the story meanwhile. S-0463 changes nothing under flai/internal/workitem or the writers.

## Remediation

Story S-0479 remediates this issue, created from it at 2026-10-21T11:59:43Z.
