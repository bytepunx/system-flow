---
id: I-0345
title: wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-10T11:10:23Z
last_reported: 2026-11-10T11:10:23Z
updated: 2026-11-11T01:43:38Z
---

# I-0345 wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Description
wait_for_events fails with no front matter while another story's checked creation holds a PLACEHOLDER task file

## Instances

### 2026-11-10T11:10:23Z
Story: S-0708.
While S-0708's close-out ran, from 2026-11-10T11:04:00Z to about 11:04:25Z, wip/kanban/tasks/T-3079-... (S-0665's new task) held only the word PLACEHOLDER. Every call S-0708's agent made to wait_for_events failed with "<path>: no front matter" until S-0665's agent finished the creation. A checked creation that leaves a placeholder in place, as ADR-0272 describes, should be one that wait_for_events, inbox, and the other readers skip, not one that fails the whole call. I-0036, closed, recorded the same failure for a read made during a save.

## Remediation

Story S-0717 remediates this issue, created from it at 2026-11-11T01:43:38Z.
