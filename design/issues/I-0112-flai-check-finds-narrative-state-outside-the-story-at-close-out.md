---
id: I-0112
title: "flai check finds `narrative.state` outside the story at close-out"
class: efficiency
status: open
count: 2
first_reported: 2026-10-07T08:55:42Z
last_reported: 2026-10-07T09:01:18Z
updated: 2026-10-07T09:01:18Z
---

# I-0112 flai check finds `narrative.state` outside the story at close-out

## Description
flai check finds `narrative.state` outside the story at close-out

## Instances

### 2026-10-07T08:55:42Z
Story: S-0275.
flai check found outside the story:
`wip/agents/S-0265.md`: S-0265 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0265 --current, or the MCP tool stream_state
`wip/agents/S-0265.md`: S-0265 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0265 --next, or the MCP tool stream_state

### 2026-10-07T09:01:18Z
Story: S-0275.
flai check found outside the story:
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Current state is empty or the template's placeholder; write what is true now with flai stream state S-0215 --current, or the MCP tool stream_state
`wip/agents/S-0215.md`: S-0215 is in-progress and its narrative's ## Next steps is empty or the template's placeholder; write them as a list, the very next action first, with flai stream state S-0215 --next, or the MCP tool stream_state

## Remediation
