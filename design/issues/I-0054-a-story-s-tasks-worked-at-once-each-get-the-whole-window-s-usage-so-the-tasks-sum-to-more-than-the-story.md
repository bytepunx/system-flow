---
id: I-0054
title: A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-01T13:02:23Z
last_reported: 2026-10-01T13:02:23Z
updated: 2026-10-01T13:02:23Z
---

# I-0054 A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story

## Description
A story's tasks worked at once each get the whole window's usage, so the tasks sum to more than the story

## Instances

### 2026-10-01T13:02:23Z
S-0176 (2026-10-01): T-0679 and T-0680 ran at once from 11:57Z; flai serve agent usage gave each the same 13.4M tokens and 5.08 US dollars, and the six tasks summed to about 22 dollars against the run's 17.79. ADR-0051 apportions a run to a task by the time it was in progress, which assumes one task at a time; with parallel task sub-agents it should apportion by the calls each sub-agent made (parent_tool_use_id) or split a shared window.

## Remediation
