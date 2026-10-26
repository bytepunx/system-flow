---
id: T-2118
type: task
nature: research
title: The agent adapters design records what a container harness of flai's own, a baseline image flai builds and starts, gives against flai's concerns, and weighs it against flai today and against Docker Agent
status: done
parent: S-0497
owner: alex
created: 2026-10-25T21:26:59Z
updated: 2026-10-25T21:31:34Z
transitions:
  - to: ready
    at: 2026-10-25T21:27:12Z
    by: system-flow
  - to: in-progress
    at: 2026-10-25T21:27:13Z
    by: system-flow
  - to: done
    at: 2026-10-25T21:31:34Z
    by: system-flow
stream: S-0497
tags: []
touches: [design/system/agent-adapters.md]
usage:
  source: log
  seconds: 261
  estimated: true
  priced_by: harness
  models:
    - model: claude-fable-5-1
      input: 245
      output: 11188
      cache_read: 1410007
      cache_write: 64448
      cost: 2.2033
---
# T-2118 The agent adapters design records what a container harness of flai's own, a baseline image flai builds and starts, gives against flai's concerns, and weighs it against flai today and against Docker Agent

## Work

Add a subsection to `design/system/agent-adapters.md`, beside the Docker Agent survey, that describes a container harness of flai's own: a baseline image flai builds, holding a thin harness flai writes that runs the agent loop, mounts only what flai chooses, carries `flai mcp` and `flai guard` inside, reads its secrets and configuration from the environment or a service such as signet, and reports usage per harness, model, and agent. State what it gives against each of flai's concerns (start and resume, MCP, roles, the guard, the permission handler, usage and cost, instruction files), add its row to the harnesses table, and weigh it in one table against flai today (Claude Code on the host) and against Docker Agent (ADR-0202). Waits for nothing: it is the first layer.

## Done when

- The subsection exists, every claim about this repository names its file or ADR, and every claim about Docker Agent cites the survey S-0496 wrote.
- The harnesses table has a row for the container harness.
- A three-way comparison table weighs flai today, Docker Agent, and the container harness on isolation, mounts, tools, the guard and the permission prompt, the subscription, local models, cost capture, operational burden, and the path to Kubernetes.

## Notes
