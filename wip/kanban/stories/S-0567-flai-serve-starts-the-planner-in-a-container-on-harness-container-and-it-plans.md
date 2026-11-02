---
id: S-0567
type: story
nature: feature
title: flai serve starts the planner in a container on harness container, and it plans one backlog story of this repository with its activity logged
status: backlog
parent: E-0029
owner: alex
created: 2026-11-02T19:02:15Z
updated: 2026-11-02T19:02:15Z
transitions: []
tags: []
topics: [harness, docker, planning]
touches: [wip/agents/planner.md, design/system/strategic-agents.md, design/issues]
after: [S-0509]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0567 flai serve starts the planner in a container on harness container, and it plans one backlog story of this repository with its activity logged

## Goal

S-0509 lets the planner, the orchestrator, the analyzer, and the curator run in a container (ADR-0247). It proved this with docker faked, since the host's flai serve could not run the story's code before its release (TH-0815). This story makes the live run: `flai serve` starts a planner on `harness: container`, and it plans one backlog story of this repository in the container. It shows that the mounts, the guard, and the activity entry hold on the real path.

## Acceptance criteria

- [ ] With the host's container harness set (`flai serve agent harness container --login <path>`) and the planner's agent on `harness: container`, `flai serve` starts a planner in `flai-agent-claude-code:local` (or the project's layer). The planner plans one backlog story of this repository, and its plan's thread and drafted tasks are written through flai.
- [ ] `wip/agents/planner.md` gains the run's activity entry, with duration and cost measured from the container's log as from a host run.
- [ ] The guard refuses the planner inside the container what it refuses on the host, shown by at least one refusal in the run's log or by a call made by hand in the container.
- [ ] Each limitation met is an issue under `design/issues`, and `design/system/strategic-agents.md` § In a container (S-0509) says what the live run showed.

## Tasks

## Notes

Waits for S-0509 (`after`), whose start paths and adapter it runs, and for a flai release that holds them. Before the run the operator sets the host's container harness login and image with `flai serve agent harness container`. An agent cannot do that, since the guard refuses it a host-configuration write (ADR-0220). Moved here from S-0509's criterion 4 on TH-0815.
