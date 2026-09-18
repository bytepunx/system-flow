---
id: S-0039
type: story
nature: feature
title: "flai mcp: inbox, reply, read, and move tools with the template .mcp.json"
status: done
parent: E-0006
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-18T17:07:30Z
transitions:
  - to: ready
    at: 2026-09-18T16:58:16Z
    by: alex
  - to: in-progress
    at: 2026-09-18T16:58:16Z
    by: alex
  - to: review
    at: 2026-09-18T17:03:20Z
    by: alex
  - to: done
    at: 2026-09-18T17:07:30Z
    by: alex
tags: [dashboard, cli]
touches: [flai/internal/mcpserver, flai/cmd/mcp.go, flai/go.mod, template/root/.mcp.json, template/root/CLAUDE.md.tmpl, design/tech/go-libraries.md]
---

# S-0039 flai mcp: inbox, reply, read, and move tools with the template .mcp.json

## Goal
`flai mcp` serves the repository to agents over the Model Context Protocol on stdio: inbox of open threads and questions, reply and resolve, read items and documents, transition items, and a long-poll for new events, so a session sees the designer's input within seconds without polling files.

## Acceptance criteria
- [x] `flai mcp` speaks MCP over stdio with tools `inbox`, `thread_reply`, `thread_resolve`, `item_get`, `item_move`, `doc_get`, `who_touches`, and `wait_for_events(timeout)`; resources expose design and docs files; every write goes through the same code as the CLI
- [x] `wait_for_events` blocks up to the timeout and returns as soon as a watched file changes (threads, items, narratives), so an idle agent reacts within a second
- [x] The template ships `.mcp.json` registering `flai mcp`, and the template `CLAUDE.md` tells the agent to check the inbox at session start, at every task transition, and before review
- [x] Behavior tests drive the server through a client; an integration test round-trips a designer reply into an agent inbox; docs/users and flai-cli.md describe the server; design/tech records the MCP Go library and version

## Tasks
- T-0141 MCP server on stdio with the Go SDK: inbox, thread_reply, thread_resolve, item_get, item_move, doc_get, who_touches; writes through the CLI's code; agents cannot accept
- T-0142 wait_for_events: block up to a timeout and return on the first change under threads, kanban, or agents; resources for design and docs
- T-0143 Template .mcp.json and CLAUDE.md inbox checkpoints; this repository's .mcp.json
- T-0144 Behavior tests through an in-memory client; integration round trip of a designer reply into the agent inbox; docs, CLI design, tech pins

## Notes
Stdio first because it needs no network and no dashboard running; S-0043 adds the HTTP transport through flaiover for remote agents and the hub.
- Added beyond the listed tools: `thread_get` and `thread_open`, so an agent can read a whole thread and ask the designer a question, which is what `## Open questions` was for. `item_move` refuses story and epic acceptance, in line with S-0046.
