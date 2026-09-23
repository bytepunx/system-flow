---
id: S-0096
type: story
nature: feature
title: Flai serve should start a flai mcp in the background as a child process
status: done
owner: alex
created: 2026-09-23T02:09:38Z
updated: 2026-09-23T03:21:30Z
transitions:
  - to: ready
    at: 2026-09-23T02:09:42Z
    by: alex
  - to: in-progress
    at: 2026-09-23T02:35:36Z
    by: system-flow
  - to: review
    at: 2026-09-23T02:45:42Z
    by: system-flow
  - to: done
    at: 2026-09-23T03:21:30Z
    by: alex
tags: [cli]
touches: [flai/cmd]
---
# S-0096 Flai serve should start a flai mcp in the background as a child process

## Goal

The MCP process should be running any time a flai serve process is started so that the agent(s) can communicate.

## Acceptance criteria
- [x] flai MCP process starts with flai serve
- [x] flai MCP process is stopped when flai serve process is stopped

## Tasks
- T-0344 flai serve keeps each served project's HTTP MCP server running and stops the ones it started
- T-0345 Recorded as an ADR refining ADR-0030; serve status and the docs say where each project's MCP listens
- T-0346 Tried with a real flai serve: the MCP answers, stops with serve, and exits by itself when serve is killed

## Notes
