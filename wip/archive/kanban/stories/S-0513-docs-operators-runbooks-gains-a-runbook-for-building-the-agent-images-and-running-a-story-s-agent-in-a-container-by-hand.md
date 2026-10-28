---
id: S-0513
type: story
nature: improvement
title: docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand
status: done
parent: E-0029
owner: alex
created: 2026-10-28T12:14:47Z
updated: 2026-10-28T12:43:01Z
transitions:
  - to: ready
    at: 2026-10-28T12:22:45Z
    by: system-flow
  - to: in-progress
    at: 2026-10-28T12:22:46Z
    by: system-flow
  - to: review
    at: 2026-10-28T12:30:04Z
    by: agent-S-0513
  - to: done
    at: 2026-10-28T12:43:01Z
    by: alex
tags: [docker]
touches: [docs/operators/runbooks/agent-container.md, docs/operators/runbooks/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 62
  priced_by: harness
  turns:
    - day: 2026-10-28
      work: 18
  models:
    - model: claude-opus-5-5
      input: 56
      output: 5688
      cache_read: 1058407
      cache_write: 89895
      cost: 1.0448
  strategic:
    - kind: orchestrator
      seconds: 1026
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 87
          output: 1266
          cache_read: 11528428
          cache_write: 20238
          cost: 2.5355
        - model: claude-sonnet-5-5
          input: 10
          output: 49
          cache_read: 112797
          cache_write: 42378
          cost: 0.1061
---
# S-0513 docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand

## Goal

An operator can build the two agent images and run a story's agent in a container by hand from one runbook, without reading S-0505's scripts or its results. The runbook is the operator-facing record of what S-0505 built: `scripts/agent-image.sh`, `make agent-images`, and `scripts/agent-container.sh`, what each mounts and needs, and how to stop and clean up. This story is also the one S-0505 works inside the container (TH-0676), so its agent runs in `flai-agent-claude-code:local` and reaches flai over the project's HTTP MCP server.

## Acceptance criteria

- [x] `docs/operators/runbooks/agent-container.md`, with the runbooks' front matter, says what the agent images are (local, unsigned, E-0029), how to build them with `make agent-images`, what `scripts/agent-container.sh S-nnnn` mounts (the worktree, `.flai-cache/agents/S-nnnn`, the Claude login read-only) and never mounts (the Docker socket), what it needs on the host (Docker, `flai host` with the project's HTTP MCP server, `flai mcp token`), and how to remove the images and the story's cache folder afterwards
- [x] The runbook names the limitations S-0505's results record and links the results document under `design/experiments/` rather than copying them
- [x] `docs/operators/runbooks/README.md` indexes the runbook, and the markdown lint and `flai check --strict` pass

## Tasks
- T-2196 docs/operators/runbooks/agent-container.md describes building the agent images and running a story's agent in a container by hand, and the runbooks README indexes it

## Notes

Documentation only: the agent images carry no Go toolchain, so this story changes no code. Created by agent-S-0505 on the operator's confirmation on TH-0676. S-0505 starts it on the host and works it inside the container; it is not for `flai serve` to start an agent for on the host.
