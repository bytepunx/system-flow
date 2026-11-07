---
id: S-0567
type: story
nature: feature
title: flai serve starts the planner in a container on harness container, and it plans one backlog story of this repository with its activity logged
status: in-progress
parent: E-0029
owner: alex
created: 2026-11-02T19:02:15Z
updated: 2026-11-07T04:07:38Z
transitions:
  - to: ready
    at: 2026-11-07T03:43:06Z
    by: alex
  - to: in-progress
    at: 2026-11-07T03:43:19Z
    by: agent-S-0567
blocked:
  - from: 2026-11-07T03:50:10Z
    until: 2026-11-07T03:54:00Z
    reason: awaiting the operator on TH-0838
  - from: 2026-11-07T03:54:43Z
    until: 2026-11-07T04:06:50Z
    reason: awaiting the operator on TH-0843
  - from: 2026-11-07T04:07:38Z
    reason: awaiting the operator on TH-0843
tags: []
topics: [harness, docker, planning]
touches: [wip/agents/planner.md, design/system/strategic-agents.md, design/issues, design/issues/summary.md]
after: [S-0509]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 488
  priced_by: harness
  turns:
    - day: 2026-11-07
      ceremony: 2
      hand_edits: 2
      work: 38
  models:
    - model: claude-opus-5-5
      input: 98
      output: 20515
      cache_read: 8592542
      cache_write: 241035
      cost: 4.0575
  strategic:
    - kind: orchestrator
      seconds: 87
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 19
          output: 139
          cache_read: 1528902
          cache_write: 16614
          cost: 0.375
forecast:
  duration: 1h
  basis: "The planner's figure on TH-0838: flai's 19m size estimate plus the operator's setup and watching one planner session, as S-0565 is forecast at 1h30m"
  by: agent-S-0567
  at: 2026-11-07T03:49:29Z
finalized:
  by: alex
  at: 2026-11-07T03:41:47Z
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
- T-2498 The live planner run's preconditions are in place, and the operator has set the planner on harness container and asked for a plan of one backlog story
- T-2499 The planner's container run is watched to its end: it plans the chosen story through flai, and wip/agents/planner.md gains its entry with duration and cost from the container's log
- T-2500 A guard refusal of the planner inside its container is shown from the run's log, or from a call made by hand in the container with the planner's role
- T-2501 Each limitation the live run met is an issue, and strategic-agents.md § In a container (S-0509) says what the live run showed
- T-2502 The live planner run is checked: its container, its plan, its activity entry's duration and cost, and a guard refusal inside the container
- T-2503 Each limitation the live planner run met is an issue, and strategic-agents.md § In a container says what the run showed

## Notes

Waits for S-0509 (`after`), whose start paths and adapter it runs, and for a flai release that holds them. Before the run the operator sets the host's container harness login and image with `flai serve agent harness container`. An agent cannot do that, since the guard refuses it a host-configuration write (ADR-0220). Moved here from S-0509's criterion 4 on TH-0815.

### Live run

Checked by agent-S-0567 on T-2498, 2026-11-07:

| Check | Found |
|-------|-------|
| Installed flai | 1.83.0 at `~/.flai/bin/flai`, built 2026-11-03; S-0509 is in every release from `flai/v1.81.0`. `flai serve` 1.83.0 runs since 2026-11-03T18:12:07Z |
| Image | `flai-agent-sf-claude-code:local`, the layer `agent.image.layer` names, over `flai-agent-claude-code:local`. The images on the host were built 2026-11-02T18:51Z, before S-0509's pilot commit (18:56Z), so `make agent-images` rebuilt all three from main at `e41d6124e`: `sha256:521d859c7bca`, Claude Code 2.1.312, flai 1.84.0 inside |
| Host container harness | program `docker`, default bridge network; login not set |
| Planner's agent | `planning.agent` names a model and no harness |
| Story to plan | S-0559, the first story `flai plan --candidates` lists with no plan |
