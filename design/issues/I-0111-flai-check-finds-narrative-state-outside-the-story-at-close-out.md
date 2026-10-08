---
id: I-0111
title: "flai check finds `narrative.state` outside the story at close-out"
class: efficiency
status: closed
count: 3
first_reported: 2026-10-07T08:55:19Z
last_reported: 2026-10-07T09:24:19Z
updated: 2026-10-08T06:17:03Z
---

# I-0111 flai check finds `narrative.state` outside the story at close-out

## Description
flai check finds `narrative.state` outside the story at close-out

## Instances

### 2026-10-07T08:55:19Z
Story: S-0265.
flai check found outside the story:
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0215 --current, or the MCP tool stream_state
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0215 --next, or the MCP tool stream_state

### 2026-10-07T09:08:04Z
Story: S-0246.
flai check found outside the story:
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0215 --current, or the MCP tool stream_state
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0215 --next, or the MCP tool stream_state

### 2026-10-07T09:24:19Z
Story: S-0246.
flai check found outside the story:
`wip/agents/S-0293.md`: S-0293 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0293 --current, or the MCP tool stream_state
`wip/agents/S-0293.md`: S-0293 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0293 --next, or the MCP tool stream_state

## Remediation

Story S-0325 remediates this issue, created from it at 2026-10-07T18:59:57Z.
Closed 2026-10-08T06:17:03Z: S-0323 (ADR-0125), which closed I-0109, the same issue opened twice: a check scoped to a story now leaves out a narrative.state finding on the narrative of another story neither done nor cancelled, so a close-out no longer records it.
