---
id: I-0277
title: A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-11-07T06:20:58Z
last_reported: 2026-11-07T06:20:58Z
updated: 2026-11-10T08:34:42Z
---

# I-0277 A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps

## Description
A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps

## Instances

### 2026-11-07T06:20:58Z
Story: S-0556.
S-0556 runs on harness container. `wip/agents/S-0556.md` is read-only in the container (ADR-0253 mounts the main checkout's wip read-only), but work-management.md, tooling.md, and T-2459's Done when require `## Decisions` to be written by hand. The MCP tool stream_state writes only `## Current state` and `## Next steps`, and task_done writes only the log. S-0556 recorded its decisions in T-2459's `## Notes` through item_edit, with a pointer in Current state, so T-2459's third Done when is met there rather than in the narrative.

## Remediation

Story S-0668 remediates this issue, created from it at 2026-11-10T08:34:42Z.
