---
id: ADR-0020
title: The designer and agents communicate through files, served to agents by flai mcp
status: accepted
date: 2026-09-17
supersedes: []
superseded_by: []
refines: [ADR-0009, ADR-0016]
---

# ADR-0020 The designer and agents communicate through files, served to agents by flai mcp

## Context

E-0006 makes flaiover the place where the designer reads, comments, answers, and directs. Agent-to-designer already has sub-second latency: agents write files through flai and the dashboard pushes changes to open tabs. Designer-to-agent had none: an agent notices a new answer only when it happens to read the file. Options were file watching alone, an HTTP API on the dashboard, or MCP.

## Decision

Durable state stays in files. Threads, answers, and decisions are markdown under `wip/threads/` with front matter that `flai check` validates, so git carries them and an agent works with nothing but the repository. `flai mcp` serves those files and the work items to agents over the Model Context Protocol: stdio for local sessions, registered by the template's `.mcp.json`; Streamable HTTP through flaiover for remote agents and a future hub, delegating to the binary per ADR-0016. Tools cover the inbox, replies, reads, transitions, and a `wait_for_events` long-poll so an idle agent reacts within a second. Waking is the agent's responsibility: check the inbox at session start, at task transitions, and before review; hold the long-poll when idle.

## Consequences

- One source of truth; the dashboard and the server watch the same files and cannot disagree.
- MCP is pull-based; the long-poll gives near-real-time only while an agent is between steps. That is accepted.
- flai gains a server mode and a Go MCP dependency, recorded in design/tech.
- The hub shape is fixed: bearer token, project identity in every response, flaiover dialing out.

## Alternatives considered

- File watching alone: no way to reach an idle agent; kept as the fallback that always works.
- HTTP API only: needs the dashboard running for every agent and has no native client in agent sessions.
- A message queue or database: a second store to keep consistent with the files; rejected.
