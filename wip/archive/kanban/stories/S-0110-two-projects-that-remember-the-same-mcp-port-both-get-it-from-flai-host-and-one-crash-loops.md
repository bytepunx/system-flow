---
id: S-0110
type: story
nature: remediation
title: Two projects that remember the same MCP port both get it from flai host, and one crash-loops
status: done
owner: alex
created: 2026-09-24T04:49:46Z
updated: 2026-09-24T05:21:52Z
transitions:
  - to: ready
    at: 2026-09-24T04:49:46Z
    by: system-flow
  - to: in-progress
    at: 2026-09-24T04:49:46Z
    by: system-flow
  - to: review
    at: 2026-09-24T04:54:50Z
    by: system-flow
  - to: done
    at: 2026-09-24T05:21:52Z
    by: alex
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

- [x] when a project's remembered MCP address is already given to another project the host runs, the project gets a free port no other project has, remembers it, and the host log says so
- [x] the crash loop is reproduced and gone: on a scratch host, two projects that remember the same port both run, and a move never takes a port another registered project remembers

## Tasks
- T-0385 A remembered MCP address another project already has is replaced by a free one, remembered, and logged
- T-0386 Tried on the operator's projects; docs

## Notes

The operator's own host confirms it once it runs this fix: restart `flai host` on the release that carries S-0110, and authstar's and system-flow's servers both run, one of them on a new port that the host log names.

Seen on the operator's host on 2026-09-24 at 04:44Z, flai 1.15.1. system-flow has remembered 127.0.0.1:4243 since 2026-09-23, when 4243 was the default for `flai mcp start`, and authstar remembered it too. `mcpAddrFor` returns a remembered address as it is, and leaves out the addresses given to other projects only when it picks a new port. Under flai serve before S-0106 the loser failed with one warning; flai host starts it again with a growing wait, so it shows as a crash loop.
