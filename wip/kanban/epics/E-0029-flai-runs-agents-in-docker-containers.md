---
id: E-0029
type: epic
nature: feature
title: flai runs agents in Docker containers
status: in-progress
owner: alex
created: 2026-10-26T04:31:29Z
updated: 2026-10-28T12:14:47Z
transitions:
  - to: ready
    at: 2026-10-26T04:31:36Z
    by: alex
  - to: in-progress
    at: 2026-10-28T12:04:42Z
    by: agent-S-0505
tags: []
usage:
  source: sum
  seconds: 1995
  priced_by: harness
  turns:
    - day: 2026-10-28
      ceremony: 1
      hand_edits: 2
      work: 92
  models:
    - model: claude-opus-5-5
      input: 372
      output: 109259
      cache_read: 20629781
      cache_write: 530116
      cost: 10.0897
  strategic:
    - kind: planner
      seconds: 190
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 546
          output: 14923
          cache_read: 3320733
          cache_write: 204659
          cost: 8.1773
    - kind: orchestrator
      seconds: 26487
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 338
          output: 4530
          cache_read: 67392913
          cache_write: 267303
          cost: 14.8408
        - model: claude-sonnet-5-5
          input: 20
          output: 108
          cache_read: 240792
          cache_write: 93399
          cost: 0.2285
cost_of_delay:
  inputs:
    penalty_per_week: 150
    by: alex
    at: 2026-10-26T12:24:54Z
---
# E-0029 flai runs agents in Docker containers

## Outcome

ADR 0203 was drafted through a research story to determine the feasibility and approach to moving agents into docker container instances.

The one refinement to make to the ADR as proposed is to look at creating multiple, layered docker images so the bedrock image has things like the flai CLI, git, etc. that would be in common across all images while each specialized image would contain the particular harness or, in the case of models that run locally, the tooling necessary to invoke and work with those.

## Stories
- S-0505 Claude Code works one story of this repository inside a hand-built container from a bedrock image and a claude-code layer, and the results settle what ADR-0203 leaves open
- S-0506 flai's HTTP MCP server serves an agent in a container as the stdio server serves one on the host: every tool, permission_prompt under auto-approve for an agent flai serve started, and a listener a container reaches
- S-0507 pilot, a Go binary in the bedrock image, is the container's entrypoint: it reads a start request, starts the inner harness through flai's adapters, relays its log, and resumes, as an ADR records beside the layered images
- S-0508 flai serve runs a story's agent in a container through a container harness adapter, with the worktree, the story's cache folder, and the login mounted, the keys and the MCP token by environment, and the image and network the host sets
- S-0509 The planner, the orchestrator, the analyzer, and the curator run in a container when their harness is container, with the main checkout mounted as an ADR settles, pilot passing their role, and their activity logged as on the host
- S-0510 A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container
- S-0513 docs/operators/runbooks gains a runbook for building the agent images and running a story's agent in a container by hand

## Notes

### Considerations for Planning Stories

- flai is capable of instantiating the correct docker image as a container for specific agents (strategic agents, story agents, task agents, etc)
- the bedrock image contains tooling that handles spinning up the correct harness and provides necessary interface(s) so that flai is able to communicate with the agent in the container in the same way it does with agents run through a harness today without limiting what harness is used (it should likely have adapters for each harness - claude-code, etc.)
- if common functionality is required across target harness/model/agent approaches, this should be built as a separate Go binary that ships in the container - we can call that binary "pilot"
- each of the docker images should be built locally for the scope of this epic and skip making changes to github CI/CD
- eagerly look at research or experiment stories where feasible to address anticipated limitations or issues in the ADR
- 2026-10-28T12:04:42Z: moved to in-progress: follows S-0505, which moved to in-progress
