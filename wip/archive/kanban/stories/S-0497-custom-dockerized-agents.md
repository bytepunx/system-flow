---
id: S-0497
type: story
nature: research
title: Custom Dockerized Agents
status: done
owner: alex
created: 2026-10-25T21:25:23Z
updated: 2026-10-26T04:01:15Z
transitions:
  - to: ready
    at: 2026-10-25T21:25:28Z
    by: alex
  - to: in-progress
    at: 2026-10-25T21:25:38Z
    by: agent-S-0497
  - to: review
    at: 2026-10-25T21:36:42Z
    by: system-flow
  - to: done
    at: 2026-10-26T04:01:15Z
    by: alex
tags: [cli]
topics: [harness, docker]
touches: [design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/adrs/0203-flai-may-run-a-story-s-agent-in-a-container-from-a-baseline-image-it-builds.md]
agent:
  harness: claude-code
  model: claude-fable-5-1
  config:
    effort: medium
usage:
  source: log
  seconds: 684
  priced_by: harness
  turns:
    - day: 2026-10-25
      hand_edits: 1
      work: 26
  models:
    - model: claude-fable-5-1
      input: 866
      output: 39482
      cache_read: 4975852
      cache_write: 227434
      cost: 7.7754
  strategic:
    - kind: orchestrator
      seconds: 23815
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 232
          output: 4018
          cache_read: 30182043
          cache_write: 57122
          cost: 6.6301
        - model: claude-sonnet-5-5
          input: 6
          output: 37
          cache_read: 78784
          cache_write: 22786
          cost: 0.0695
forecast:
  duration: 25m
  delivery: 2026-10-25T22:10:00Z
  basis: "flai forecast gave 41m (median 482 s per unit over 4 done research stories in the small band, size 5: 2 criteria, 3 touches) and delivery 2026-10-26T01:01Z by a cycle factor of 5.25; lowered to 25m because S-0496, the same shape (a section of agent-adapters.md and a draft ADR) on the same agent and model, was in progress 10m15s and done 17m35s after it started, and this story adds a three-way comparison but surveys no new product; delivery is the start at 2026-10-25T21:26Z plus 25m times S-0496's observed cycle factor of 1.7."
  by: planner-S-0497
  at: 2026-10-25T21:28:15Z
---
# S-0497 Custom Dockerized Agents

## Goal

Docker Agents seem to have quite a few limitations and prerequisites that don't make it a clear win over allowing flai to work with Docker images built on a custom harness that allows flai to control different harness/agent combinations such that flai controls exactly what gets mounted to the container, what tools are available, the flai MCP and flai guard, and controls over how the cost is captured or calculated for specific harness/model/agent combinations.

Provide a draft ADR that considers this approach to allowing flai to distribute work across agents that utilize any combination of hosted models and local models using a custom baseline Docker image.

Note: it is also worth keeping in mind that flai could eventually integrate with kubernetes as the means of instantiating agents through this more customized approach and the harness we write in the baseline image could be used to access secrets and configuration for a given agent from signet.

## Acceptance criteria
- [x] A written ADR that examines taking a more custom approach to delegating work to agents running in Docker containers built on one or more custom image baselines
- [x] A careful analysis of the trade-offs in this approach vs. the current way flai works vs. using Docker agents

## Tasks
- T-2118 The agent adapters design records what a container harness of flai's own, a baseline image flai builds and starts, gives against flai's concerns, and weighs it against flai today and against Docker Agent
- T-2119 A draft ADR proposes custom dockerized agents on a flai baseline image, with the trade-offs against flai today and Docker Agent, its limitations, its security concerns, and the Kubernetes and signet directions

## Notes

### Planning

Planned by planner-S-0497 on 2026-10-25, after the story's agent had started it and written T-2118 and T-2119.

Touches. `flai touches suggest S-0497` refused: the story declared no touches, so there were no seeds for co-change. Every touch comes from the layout of the research story this one continues, S-0496, and from the design:

- `design/system/agent-adapters.md` (design, layout): S-0496 recorded its Docker Agent survey as a section of this document, beside the harnesses it compares, and ADR-0202 cites it. The analysis of a custom baseline image against today's `claude-code` on the host and against Docker Agent belongs beside it, in a section of its own, so that the ADR cites the comparison rather than carrying it. T-2118 names it.
- `design/adrs/README.md` (layout): `flai adr new` writes the draft ADR's index row. T-2119 names it.
- `design/adrs` (layout, folder kept): the draft ADR's number and slug are given by `flai adr new` from the title the story's agent chooses when it is written, so no task can name the file. The next number is 0203 and no story in review carries an ADR, but the slug is the title's. `flai task done` widens the touches with the committed path (ADR-0144). The folder is the only one kept, and only for that file.

Forecast. See the forecast's basis.

Cost of delay. The story has no inputs and no epic, so `flai cod` refuses a value; asked on TH-0645, recommending `penalty_per_week: 50` USD as S-0496 had. The value is written once the operator sets an input.
