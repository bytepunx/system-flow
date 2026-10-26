---
id: T-2175
type: task
nature: feature
title: A planner run in a container plans one story of this repository with its activity logged, each limitation is an issue, and the docs say how a strategic agent runs in a container
status: backlog
parent: S-0509
owner: alex
created: 2026-10-26T04:51:44Z
updated: 2026-10-26T04:51:44Z
transitions: []
stream: S-0509
tags: [cli, docker]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/issues/summary.md]
after: [T-2170]
---
# T-2175 A planner run in a container plans one story of this repository with its activity logged, each limitation is an issue, and the docs say how a strategic agent runs in a container

## Work

With the planner's harness set to `container` (the project's default agent, or the planner's own setting on the strategic settings page), ask `flai serve` to plan one backlog story of this repository and let the run finish: `wip/agents/planner.md` gains the activity entry with duration and cost measured from the container's log, the story gains its plan, and the guard's refusals inside the container match the host's. Try the orchestrator's start and stop once with its host action, if the operator allows, and record what happened. Record each limitation met as an issue with `flai issue new`. Then write the documentation: `design/system/strategic-agents.md` (each role's "Its agent" section says how it runs in a container and what it mounts, following T-2158's paragraph), `design/system/flai-cli.md` (the strategic start paths through the adapter), `docs/users/flai.md`, and `docs/operators/index.md`. Waits for T-2170, which finishes the start paths.

## Done when

A planner run in a container has planned a story and is logged in `wip/agents/planner.md` with its cost, each limitation met has an issue, the four documents describe the behaviour the tests check, and `flai check --strict` and the markdown lint pass.

## Notes

The proof run costs a planner session; name the story it planned in the story's notes.
