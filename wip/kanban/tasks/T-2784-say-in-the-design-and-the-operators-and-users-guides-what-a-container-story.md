---
id: T-2784
type: task
nature: improvement
title: Say in the design and the operators' and users' guides what a container story runs inside and what still runs on the host
status: done
parent: S-0628
owner: alex
created: 2026-11-09T21:59:18Z
updated: 2026-11-10T06:55:26Z
transitions:
  - to: ready
    at: 2026-11-10T06:50:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T06:50:43Z
    by: system-flow
  - to: done
    at: 2026-11-10T06:55:26Z
    by: agent-S-0628
stream: S-0628
tags: [docs, security]
touches: [design/system/flai-cli.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/tech/docker.md]
after: [T-2772, T-2777, T-2782]
usage:
  source: log
  seconds: 283
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 70
      output: 387
      cache_read: 12868630
      cache_write: 59601
      cost: 4.7549
---
# T-2784 Say in the design and the operators' and users' guides what a container story runs inside and what still runs on the host

## Work

Bring the documents in line with ADR-0272 and the code T-2777 and T-2782 wrote. Today none says that the host runs the worktree's scripts. Some contradict ADR-0253 on where the close-out runs.

- `design/system/flai-cli.md`: change `### The close-out job` and `### The container harness (S-0508)`. The second's "What it does not do" says the close-out is asked of `flai serve` because only the host commits. Describe the job's container run and the host's commit and sync. In `### A strategic agent's container (S-0509)`, say that the orchestrator's `flai verify` and `flai accept` run the tiers on the host.
- `design/system/strategic-agents.md` § In a container (S-0509): one paragraph on what a strategic agent's container does not keep off the host.
- `docs/operators/index.md`: in the `**In a container**` bullet of `### Starting an agent when a story becomes ready`, say what runs inside and what runs on the host. Do the same in `### The checks host action (S-0082)` and in `### Strategic agents in a container (S-0509)`.
- `docs/operators/runbooks/agent-container.md`: the close-out container's mounts beside the agent's.
- `docs/users/flai.md` § In a container: the container harness. Its line saying the close-out cannot run in the container (S-0520) becomes what the job now does. This file is in `claims.shared`.

Waits for T-2772, whose decision it describes, and for T-2777 and T-2782, whose behaviour it documents.

## Done when

- Each section above names what runs in the container and what runs on the host, in the ADR's words, and links ADR-0272.
- No document still says the close-out of a container story runs on the host, or that it cannot run in the container.
- `flai check --strict` and the markdown lint pass on the five files.

## Notes
