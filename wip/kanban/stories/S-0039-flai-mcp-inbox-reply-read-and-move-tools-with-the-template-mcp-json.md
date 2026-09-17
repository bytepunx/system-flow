---
id: S-0039
type: story
nature: feature
title: "flai mcp: inbox, reply, read, and move tools with the template .mcp.json"
status: backlog
parent: E-0006
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-17T19:46:28Z
transitions: []
tags: [dashboard, cli]
---

# S-0039 flai mcp: inbox, reply, read, and move tools with the template .mcp.json

## Goal
`flai mcp` serves the repository to agents over the Model Context Protocol on stdio: inbox of open threads and questions, reply and resolve, read items and documents, transition items, and a long-poll for new events, so a session sees the designer's input within seconds without polling files.

## Acceptance criteria
- [ ] `flai mcp` speaks MCP over stdio with tools `inbox`, `thread_reply`, `thread_resolve`, `item_get`, `item_move`, `doc_get`, `who_touches`, and `wait_for_events(timeout)`; resources expose design and docs files; every write goes through the same code as the CLI
- [ ] `wait_for_events` blocks up to the timeout and returns as soon as a watched file changes (threads, items, narratives), so an idle agent reacts within a second
- [ ] The template ships `.mcp.json` registering `flai mcp`, and the template `CLAUDE.md` tells the agent to check the inbox at session start, at every task transition, and before review
- [ ] Behavior tests drive the server through a client; an integration test round-trips a designer reply into an agent inbox; docs/users and flai-cli.md describe the server; design/tech records the MCP Go library and version

## Tasks

## Notes
Stdio first because it needs no network and no dashboard running; S-0043 adds the HTTP transport through flaiover for remote agents and the hub.
