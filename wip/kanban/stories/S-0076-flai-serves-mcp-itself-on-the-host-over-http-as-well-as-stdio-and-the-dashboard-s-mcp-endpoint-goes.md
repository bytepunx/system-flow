---
id: S-0076
type: story
nature: feature
title: flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:50Z
updated: 2026-09-20T07:28:09Z
transitions:
  - to: ready
    at: 2026-09-20T07:28:09Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, design/adrs, template]
---
# S-0076 flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes

## Goal
Agents work with flai on the host, not through the dashboard. `flai mcp` keeps stdio and gains a server of its own over Streamable HTTP for agents that need it; flaiover's `/mcp` endpoint and its bridge are removed.

## Acceptance criteria
- [ ] `flai mcp` serves stdio as now; an HTTP mode starts a server process on the host, loopback by default, with bearer authentication, the agent named as today, and sessions capped; how it is started, stopped, and found (its own command, or part of the host flai) is decided here and recorded
- [ ] The tools, rules, and cursor are the same code on both transports; the MCP revision targeted and what changed in it (sessions, the GET stream) is recorded in `design/tech`
- [ ] flaiover's `/mcp` route, `mcpbridge.ts`, and their settings are removed; a request to `/mcp` says where MCP lives now
- [ ] An ADR supersedes the MCP-over-HTTP part of ADR-0024 and keeps its project identity; `.mcp.json` in the template and here still works unchanged over stdio
- [ ] The conventions, the design documents, and the users' documentation say how an agent connects; convention baselines change in `template/` first

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Decided by the operator on 2026-09-20: "lets move the MCP integration down to flai as a process again (flai MCP can start a server process) since we want flai and agents working together on the host and not having agents interacting with a project through the dashboard." Independent of the read and write stories; must land before the mount is removed, because `/mcp` today needs the mount.
