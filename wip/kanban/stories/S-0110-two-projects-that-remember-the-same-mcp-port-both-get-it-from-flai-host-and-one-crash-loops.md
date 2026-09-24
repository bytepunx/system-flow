---
id: S-0110
type: story
nature: remediation
title: Two projects that remember the same MCP port both get it from flai host, and one crash-loops
status: backlog
owner: alex
created: 2026-09-24T04:49:46Z
updated: 2026-09-24T04:49:46Z
transitions: []
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0110 Two projects that remember the same MCP port both get it from flai host, and one crash-loops

## Goal

flai host starts every project's HTTP MCP server with an address of its own, even when two projects remember the same one.

## Acceptance criteria

- [ ] when a project's remembered MCP address is already given to another project the host runs, the project gets a free port no other project has, remembers it, and the host log says so
- [ ] the operator's host on 1.15.1 crash-looped authstar's MCP server (bind: address already in use on 127.0.0.1:4243, which system-flow also remembers); with the fix, every project's server runs

## Tasks

## Notes

Seen on the operator's host on 2026-09-24 at 04:44Z, flai 1.15.1. system-flow has remembered 127.0.0.1:4243 since 2026-09-23, when 4243 was the default for `flai mcp start`, and authstar remembered it too. `mcpAddrFor` returns a remembered address as it is, and leaves out the addresses given to other projects only when it picks a new port. Under flai serve before S-0106 the loser failed with one warning; flai host starts it again with a growing wait, so it shows as a crash loop.
