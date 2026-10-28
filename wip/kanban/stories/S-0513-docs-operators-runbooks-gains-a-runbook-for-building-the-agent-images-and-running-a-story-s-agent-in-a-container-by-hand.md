---
id: S-0513
type: story
nature: improvement
title: docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand
status: backlog
parent: E-0029
owner: alex
created: 2026-10-28T12:14:47Z
updated: 2026-10-28T12:14:47Z
transitions: []
tags: [docker]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0513 docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand

## Goal

An operator can build the two agent images and run a story's agent in a container by hand from one runbook, without reading S-0505's scripts or its results. The runbook is the operator-facing record of what S-0505 built: `scripts/agent-image.sh`, `make agent-images`, and `scripts/agent-container.sh`, what each mounts and needs, and how to stop and clean up. This story is also the one S-0505 works inside the container (TH-0676), so its agent runs in `flai-agent-claude-code:local` and reaches flai over the project's HTTP MCP server.

## Acceptance criteria

- [ ] `docs/operators/runbooks/agent-container.md`, with the runbooks' front matter, says what the agent images are (local, unsigned, E-0029), how to build them with `make agent-images`, what `scripts/agent-container.sh S-nnnn` mounts (the worktree, `.flai-cache/agents/S-nnnn`, the Claude login read-only) and never mounts (the Docker socket), what it needs on the host (Docker, `flai host` with the project's HTTP MCP server, `flai mcp token`), and how to remove the images and the story's cache folder afterwards
- [ ] The runbook names the limitations S-0505's results record and links the results document under `design/experiments/` rather than copying them
- [ ] `docs/operators/runbooks/README.md` indexes the runbook, and the markdown lint and `flai check --strict` pass

## Tasks

## Notes

Documentation only: the agent images carry no Go toolchain, so this story changes no code. Created by agent-S-0505 on the operator's confirmation on TH-0676. S-0505 starts it on the host and works it inside the container; it is not for `flai serve` to start an agent for on the host.
