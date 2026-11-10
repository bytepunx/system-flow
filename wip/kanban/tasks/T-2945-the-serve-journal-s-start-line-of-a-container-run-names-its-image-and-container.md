---
id: T-2945
type: task
nature: improvement
title: The serve journal's start line of a container run names its image and container
status: backlog
parent: S-0665
owner: alex
created: 2026-11-10T08:50:39Z
updated: 2026-11-10T08:50:39Z
transitions: []
stream: S-0665
tags: [serve, containers]
touches: [flai/internal/serve/agents_start.go, flai/internal/serve/plan.go, flai/internal/serve/orchestrate.go, flai/internal/serve/analyze.go, flai/internal/serve/curate.go, flai/internal/serve/container_test.go]
after: [T-2943]
---
# T-2945 The serve journal's start line of a container run names its image and container

## Work

The journal's start entry reads `started docker (container) ...` for a container run today. Make it name the image and the container, such as `started docker (container, image flai-agent-sf-claude-code:local, container flai-agent-plan-<session>)`, from the fields the run now records. Write one helper for the `started %s (%s)` part and use it in every start line: the story's variants in `flai/internal/serve/agents_start.go`, and the planner's, orchestrator's, analyzer's, and curator's in `plan.go`, `orchestrate.go`, `analyze.go`, and `curate.go`. A run on any other harness keeps its line as it is. The mounts stay in `agents.json` alone, since they would make the line too long to read.

It waits for T-2943, whose run fields it reads, and which changes `agents_start.go` and `container_test.go` too.

## Done when

- A container run's start entry in the journal names its image and its container's name, for a story's agent and for a strategic agent's.
- A run on another harness logs the line it logged before.
- A test in `flai/internal/serve/container_test.go` checks the journal entry of a container story run and a container planner run.
- `flai test` passes on the changed paths.

## Notes
