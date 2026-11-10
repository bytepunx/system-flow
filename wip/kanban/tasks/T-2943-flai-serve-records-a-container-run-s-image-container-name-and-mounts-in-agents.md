---
id: T-2943
type: task
nature: improvement
title: flai serve records a container run's image, container name, and mounts in agents.json, with a test reproducing I-0274
status: backlog
parent: S-0665
owner: alex
created: 2026-11-10T08:50:29Z
updated: 2026-11-10T08:50:29Z
transitions: []
stream: S-0665
tags: [serve, containers]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_start.go, flai/internal/serve/container.go, flai/internal/serve/container_test.go]
after: [T-2942]
---
# T-2943 flai serve records a container run's image, container name, and mounts in agents.json, with a test reproducing I-0274

## Work

Give `AgentRun` (`flai/internal/serve/agents.go`) three fields beside `harness_version`: `image`, `container`, and `mounts`, each empty on a run on any other harness. In `startOf` (`flai/internal/serve/agents_start.go`), which a story's run and every strategic agent's go through, copy them from the container field of the adapter's `Start` onto the run, so `agents.json` keeps them after the container's `--rm` removes it. Where `flai/internal/serve/container.go` derives the container's name again for a stop (`stopContainer`), take the recorded name when the run has one.

Reproduce I-0274 first in `flai/internal/serve/container_test.go`: a story's run and a planner's run on the container harness, with the project's layer image (`agent.image.layer`), whose `agents.json` names no image, container, or mounts. Then make it pass: the record names the layer image, `flai-agent-<subject>-<session>`, and the mounts the argv gave.

It waits for T-2942, whose `Start` field it reads.

## Done when

- A container run's entry in `agents.json` names its image, its container's name, and its mounts, for a story's agent and for a strategic agent's.
- The test fails on main as I-0274 describes and passes with the change.
- A run on another harness records none of the three.
- `flai test` passes on the changed paths.

## Notes
