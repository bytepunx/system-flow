---
id: E-0029
type: epic
nature: feature
title: flai runs agents in Docker containers
status: backlog
owner: alex
created: 2026-10-26T04:31:29Z
updated: 2026-10-26T04:31:29Z
transitions: []
tags: []
---
# E-0029 flai runs agents in Docker containers

## Outcome

ADR 0203 was drafted through a research story to determine the feasibility and approach to moving agents into docker container instances.

The one refinement to make to the ADR as proposed is to look at creating multiple, layered docker images so the bedrock image has things like the flai CLI, git, etc. that would be in common across all images while each specialized image would contain the particular harness or, in the case of models that run locally, the tooling necessary to invoke and work with those.

## Stories

## Notes

### Considerations for Planning Stories

- flai is capable of instantiating the correct docker image as a container for specific agents (strategic agents, story agents, task agents, etc)
- the bedrock image contains tooling that handles spinning up the correct harness and provides necessary interface(s) so that flai is able to communicate with the agent in the container in the same way it does with agents run through a harness today without limiting what harness is used (it should likely have adapters for each harness - claude-code, etc.)
- if common functionality is required across target harness/model/agent approaches, this should be built as a separate Go binary that ships in the container - we can call that binary "pilot"
- each of the docker images should be built locally for the scope of this epic and skip making changes to github CI/CD
- eagerly look at research or experiment stories where feasible to address anticipated limitations or issues in the ADR
