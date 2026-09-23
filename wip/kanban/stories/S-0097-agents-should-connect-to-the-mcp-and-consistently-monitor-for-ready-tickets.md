---
id: S-0097
type: story
nature: improvement
title: Agents should connect to the MCP and consistently monitor for ready tickets
status: ready
owner: alex
created: 2026-09-23T02:18:31Z
updated: 2026-09-23T02:18:42Z
transitions:
  - to: ready
    at: 2026-09-23T02:18:42Z
    by: alex
tags: [cli]
touches: [flai/cmd]
---
# S-0097 Agents should connect to the MCP and consistently monitor for ready tickets

## Goal

Agents should connect to the MCP and monitor for tickets moved to ready and pick them up as soon as there is room in In Process and Ready.

## Acceptance criteria
- [ ] when there is capacity in "in process" and tickets in ready, a ticket should be pulled to work on
- [ ] when there is not capacity in "in process" and tickets in ready, the agent should wait for available capacity and then pull from ready
- [ ] if no tickets are in ready, the agent should monitor via the MCP for tickets to land in ready, then pull the next ticket

## Tasks

## Notes
