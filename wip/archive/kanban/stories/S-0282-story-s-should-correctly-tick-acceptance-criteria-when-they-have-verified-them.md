---
id: S-0282
type: story
nature: remediation
title: Story's should correctly tick acceptance criteria when they have verified them
status: done
owner: alex
created: 2026-10-06T02:59:32Z
updated: 2026-10-06T05:58:56Z
transitions:
  - to: ready
    at: 2026-10-06T02:59:33Z
    by: alex
  - to: in-progress
    at: 2026-10-06T03:45:39Z
    by: agent-S-0282
  - to: review
    at: 2026-10-06T05:37:33Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T05:58:56Z
    by: alex
tags: [cli]
touches: [flai/cmd, flai/internal/itemedit, flai/internal/workitem, flai/internal/guard, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/harness, "flaiover/src/routes/api/items/[id]/criteria", flaiover/src/lib/server/agent.ts, design/adrs, design/system/flai-cli.md, design/system/workflow.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flai-reference.md, design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/CHANGELOG.md, docs/operators/settings.md, docs/users/conventions.md, template/template.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6746
  models:
    - model: claude-opus-5-5
      input: 468
      output: 137800
      cache_read: 22071683
      cache_write: 649633
      cost: 11.0834
    - model: claude-sonnet-5-5
      input: 32
      output: 7378
      cache_read: 571723
      cache_write: 106551
      cost: 0.4546
cost_of_delay:
  inputs:
    time_lost_per_cycle: 30m
    by: alex
    at: 2026-10-06T02:59:32Z
---
# S-0282 Story's should correctly tick acceptance criteria when they have verified them

## Goal

When an agent's story is validating completed criteria or a task sub-agent is completing criteria, it should correctly tick that acceptance criteria as completed on the story. flai should expose a command in the CLI, HTTP, and MCP for ticking acceptance criteria to prevent the AI's from needing to edit the file or the specific syntax change to make.

## Acceptance criteria
- [x] criteria are ticked consistently when they are completed, regardless by which agent, through flai

## Tasks
- T-0991 itemedit reads a body's acceptance criteria and ticks or unticks them by number
- T-0992 flai move warns when a story goes to review with unticked criteria
- T-0993 flai criteria lists, ticks, and unticks a story's acceptance criteria
- T-0994 The MCP tool criteria_tick ticks and unticks a story's acceptance criteria
- T-0995 The host API action item.criteria and the dashboard's HTTP route tick a story's criteria
- T-0996 An ADR, the design, docs, conventions, and flai serve's prompt say criteria are ticked through flai once verified

## Notes
