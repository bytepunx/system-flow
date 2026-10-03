---
id: S-0203
type: story
nature: improvement
title: A story created from an issue is a draft and carries the issue's cost of delay inputs
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:12Z
updated: 2026-10-03T05:34:14Z
transitions:
  - to: ready
    at: 2026-10-03T05:34:14Z
    by: alex
tags: [flai]
touches: [flai/internal/issues, flai/cmd/issue.go, flai/internal/mcpserver]
after: [S-0199, S-0198]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0203 A story created from an issue is a draft and carries the issue's cost of delay inputs

## Goal

S-0198 gives flai a step that creates a remediation or improvement story from an issue. Such a story is an agent's draft until the operator finalizes it, and when the issue records an impact (the analyzer's, or a `cost`), the story should carry it as cost of delay inputs so the planner and orchestrator can rank it.

## Acceptance criteria
- [ ] The story `flai issue story I-nnnn` (and the MCP equivalent, and S-0198's warning flow) creates has `draft: true`
- [ ] When the issue carries `cost` (time per occurrence) and `count`, the story's `cost_of_delay.inputs.time_lost_per_cycle` is set from them, with `by: flai` and the derivation in the story's Notes; an issue the analyzer wrote with an `impact` section carries its revenue or penalty figures over
- [ ] The story links the issue and the issue's Remediation section links the story
- [ ] Tests cover the draft flag and the inputs carried over; the user guide and `design/system/continuous-improvement.md` say so

## Tasks

## Notes
