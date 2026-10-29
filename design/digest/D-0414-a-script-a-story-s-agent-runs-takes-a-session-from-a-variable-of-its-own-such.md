---
id: D-0414
title: A script a story's agent runs takes a session from a variable of its own, such as FLAI_AGENT_SESSION, never FLAI_SESSION, since flai guard refuses a story's agent a command that sets FLAI_SESSION
kind: decision
sources: [S-0507]
topics: [harness, docker, guard, scripts]
paths: [scripts/agent-container.sh, scripts/README.md]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0414 A script a story's agent runs takes a session from a variable of its own, such as FLAI_AGENT_SESSION, never FLAI_SESSION, since flai guard refuses a story's agent a command that sets FLAI_SESSION

S-0507's `scripts/agent-container.sh` starts a story's agent in a container and needs a session to start or resume. It reads that session from `FLAI_AGENT_SESSION`, not `FLAI_SESSION`. S-0507 recorded why: "flai guard refuses a story's agent a command that sets `FLAI_SESSION` (ADR-0143), and the caller's own session is not the container's."

A later script follows the same rule when a story's agent may run it and it takes a session or an author for something other than the caller: a story's agent, a strategic agent, or a container. Give that input a variable of its own. Do not reuse `FLAI_SESSION` or `FLAI_AGENT`, which name the caller and which the guard holds to the caller's own values. `scripts/README.md` lists the script's variables.

The same story passes the MCP token to the container with `docker run -e NAME`, copied from the script's environment, never as an argument or in a file.

Source: S-0507 (wip/archive/agents/S-0507.md, Decisions).
