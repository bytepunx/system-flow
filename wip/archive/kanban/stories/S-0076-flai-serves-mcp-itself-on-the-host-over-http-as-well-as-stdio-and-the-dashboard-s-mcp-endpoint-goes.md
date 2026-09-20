---
id: S-0076
type: story
nature: feature
title: flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T12:57:01Z
transitions:
  - to: ready
    at: 2026-09-20T07:28:09Z
    by: alex
  - to: in-progress
    at: 2026-09-20T12:11:42Z
    by: system-flow
  - to: review
    at: 2026-09-20T12:38:12Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:57:01Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, design/adrs, template]
---
# S-0076 flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes

## Goal
Agents work with flai on the host, not through the dashboard. `flai mcp` keeps stdio and gains a server of its own over Streamable HTTP for agents that need it; flaiover's `/mcp` endpoint and its bridge are removed.

## Acceptance criteria
- [x] `flai mcp` serves stdio as now; an HTTP mode starts a server process on the host, loopback by default, with bearer authentication, the agent named as today, and sessions capped; how it is started, stopped, and found (its own command, or part of the host flai) is decided here and recorded
- [x] The tools, rules, and cursor are the same code on both transports; the MCP revision targeted and what changed in it (sessions, the GET stream) is recorded in `design/tech`
- [x] flaiover's `/mcp` route, `mcpbridge.ts`, and their settings are removed; a request to `/mcp` says where MCP lives now
- [x] An ADR supersedes the MCP-over-HTTP part of ADR-0024 and keeps its project identity; `.mcp.json` in the template and here still works unchanged over stdio
- [x] The conventions, the design documents, and the users' documentation say how an agent connects; convention baselines change in `template/` first

## Tasks
- T-0275 flai push and flai accept send tags three at a time, the branch last (I-0026)
- T-0276 internal/mcphttp: the MCP server over Streamable HTTP, with bearer, naming, caps, and project identity
- T-0277 flai mcp http, start, stop, status, and token: the server as a process of its own
- T-0278 flaiover: /mcp, the bridge, flai.ts, and their settings go; /mcp says where MCP lives now
- T-0279 ADR-0030, the design, the conventions, and the documentation say how an agent connects
- T-0280 End to end in a scratch project: an agent over HTTP beside a dashboard with no /mcp

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Decided by the operator on 2026-09-20: "lets move the MCP integration down to flai as a process again (flai MCP can start a server process) since we want flai and agents working together on the host and not having agents interacting with a project through the dashboard." Independent of the read and write stories; must land before the mount is removed, because `/mcp` today needs the mount.
