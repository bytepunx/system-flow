---
id: I-0277
title: A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-11-07T06:20:58Z
last_reported: 2026-11-07T06:20:58Z
updated: 2026-11-10T10:21:42Z
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
Closed 2026-11-10T10:21:42Z: S-0668: flai stream state (--context, --decision, --question), the MCP tool stream_state (context, decisions, questions), and the host channel's stream.state now replace ## Context and append ## Decisions and dated ## Open questions bullets, keeping the threads block, through Repo.WriteStreamState. The MCP server runs on the host, so an agent whose wip/ is read-only in a container records them through it. TestStreamStateWritesContextDecisionsAndQuestions in flai/internal/mcpserver/stream_sections_test.go reproduces the case.
