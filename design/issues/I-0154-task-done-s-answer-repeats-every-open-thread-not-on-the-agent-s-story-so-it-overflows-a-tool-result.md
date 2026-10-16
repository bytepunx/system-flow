---
id: I-0154
title: task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-10-16T06:05:59Z
last_reported: 2026-10-16T06:05:59Z
updated: 2026-10-16T06:05:59Z
---

# I-0154 task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result

## Description
task_done's answer repeats every open thread not on the agent's story, so it overflows a tool result

## Instances

### 2026-10-16T06:05:59Z
Story: S-0413.
S-0413, 2026-10-16: each task_done answer carried the inbox with 14 threads `awaiting: you` that are on other epics and stories (E-0021 to E-0025, S-0407, S-0426, S-0427, system-flow.yaml), whole, on every call. The first answer ran to 72,938 characters and Claude Code saved it to a file instead of returning it, so the agent had to extract the outcome with jq; later answers were the same size in context. The threads await the operator, not agent-S-0413. A story's agent needs the threads on its own story and its tasks, and its messages.

## Remediation
