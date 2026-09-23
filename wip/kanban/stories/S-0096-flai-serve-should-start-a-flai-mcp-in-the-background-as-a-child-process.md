---
id: S-0096
type: story
nature: feature
title: Flai serve should start a flai mcp in the background as a child process
status: ready
owner: alex
created: 2026-09-23T02:09:38Z
updated: 2026-09-23T02:09:42Z
transitions:
  - to: ready
    at: 2026-09-23T02:09:42Z
    by: alex
tags: [cli]
touches: [flai/cmd]
---
# S-0096 Flai serve should start a flai mcp in the background as a child process

## Goal

The MCP process should be running any time a flai serve process is started so that the agent(s) can communicate.

## Acceptance criteria
- [ ] flai MCP process starts with flai serve
- [ ] flai MCP process is stopped when flai serve process is stopped

## Tasks

## Notes
