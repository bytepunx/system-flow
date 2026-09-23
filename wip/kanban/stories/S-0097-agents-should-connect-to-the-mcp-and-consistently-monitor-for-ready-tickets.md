---
id: S-0097
type: story
nature: improvement
title: Agents should connect to the MCP and consistently monitor for ready tickets
status: review
owner: alex
created: 2026-09-23T02:18:31Z
updated: 2026-09-23T02:53:39Z
transitions:
  - to: ready
    at: 2026-09-23T02:18:42Z
    by: alex
  - to: in-progress
    at: 2026-09-23T02:45:57Z
    by: system-flow
  - to: review
    at: 2026-09-23T02:53:39Z
    by: system-flow
tags: [cli]
touches: [flai/cmd]
---
# S-0097 Agents should connect to the MCP and consistently monitor for ready tickets

## Goal

Agents should connect to the MCP and monitor for tickets moved to ready and pick them up as soon as there is room in In Process and Ready.

## Acceptance criteria
- [x] when there is capacity in "in process" and tickets in ready, a ticket should be pulled to work on
- [x] when there is not capacity in "in process" and tickets in ready, the agent should wait for available capacity and then pull from ready
- [x] if no tickets are in ready, the agent should monitor via the MCP for tickets to land in ready, then pull the next ticket

## Tasks
- T-0347 wait_for_work returns the story to pull once there is room and a story is ready, and waits otherwise
- T-0348 The MCP instructions and the template's conventions tell an idle agent to hold wait_for_work and pull what it returns
- T-0349 Tried against a real MCP server: no room, then room; nothing ready, then ready

## Notes
