---
id: S-0497
type: story
nature: research
title: Custom Dockerized Agents
status: backlog
owner: alex
created: 2026-10-25T21:25:23Z
updated: 2026-10-25T21:25:23Z
transitions: []
tags: []
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: medium
---
# S-0497 Custom Dockerized Agents

## Goal

Docker Agents seem to have quite a few limitations and prerequisites that don't make it a clear win over allowing flai to work with Docker images built on a custom harness that allows flai to control different harness/agent combinations such that flai controls exactly what gets mounted to the container, what tools are available, the flai MCP and flai guard, and controls over how the cost is captured or calculated for specific harness/model/agent combinations.

Provide a draft ADR that considers this approach to allowing flai to distribute work across agents that utilize any combination of hosted models and local models using a custom baseline Docker image.

Note: it is also worth keeping in mind that flai could eventually integrate with kubernetes as the means of instantiating agents through this more customized approach and the harness we write in the baseline image could be used to access secrets and configuration for a given agent from signet.

## Acceptance criteria
- [ ] A written ADR that examines taking a more custom approach to delegating work to agents running in Docker containers built on one or more custom image baselines
- [ ] A careful analysis of the trade-offs in this approach vs. the current way flai works vs. using Docker agents

## Tasks

## Notes
