---
id: I-0109
title: "flai check finds `narrative.state` outside the story at close-out"
class: efficiency
status: open
count: 3
first_reported: 2026-10-07T08:48:03Z
last_reported: 2026-10-07T14:49:15Z
updated: 2026-10-07T14:49:15Z
---

# I-0109 flai check finds `narrative.state` outside the story at close-out

## Description
flai check finds `narrative.state` outside the story at close-out

## Instances

### 2026-10-07T08:48:03Z
Story: S-0214.
flai check found outside the story:
`wip/agents/S-0265.md`: S-0265 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0265 --current, or the MCP tool stream_state
`wip/agents/S-0265.md`: S-0265 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0265 --next, or the MCP tool stream_state

### 2026-10-07T08:48:29Z
Story: S-0214.
flai check found outside the story:
`wip/agents/S-0308.md`: S-0308 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0308 --current, or the MCP tool stream_state
`wip/agents/S-0308.md`: S-0308 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0308 --next, or the MCP tool stream_state

### 2026-10-07T14:49:15Z
Story: S-0311.
flai check found outside the story:
`wip/agents/S-0298.md`: S-0298 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0298 --current, or the MCP tool stream_state
`wip/agents/S-0298.md`: S-0298 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0298 --next, or the MCP tool stream_state

## Remediation
