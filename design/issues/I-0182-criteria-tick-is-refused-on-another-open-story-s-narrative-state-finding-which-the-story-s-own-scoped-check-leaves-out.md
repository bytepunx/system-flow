---
id: I-0182
title: criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out
class: defect
status: open
count: 1
cost: 2m
first_reported: 2026-10-21T12:14:14Z
last_reported: 2026-10-21T12:14:14Z
updated: 2026-10-21T12:14:14Z
---

# I-0182 criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out

## Description
criteria_tick is refused on another open story's narrative.state finding, which the story's own scoped check leaves out

## Instances

### 2026-10-21T12:14:14Z
Story: S-0474.
S-0474: criteria_tick S-0474 [1] at about 12:13Z was refused with "flai check has 2 finding(s) with this change", both narrative.state on wip/agents/S-0471.md, a story its agent had just started and not yet written. task_done's check scoped to S-0474 passed at the same time, and work-management.md says a narrative.state on another open story's narrative is left out. The tick went through once S-0471's agent had written its state.

## Remediation
