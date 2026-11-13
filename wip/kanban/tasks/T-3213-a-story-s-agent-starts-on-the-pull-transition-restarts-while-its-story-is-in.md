---
id: T-3213
type: task
nature: feature
title: A story's agent starts on the pull transition, restarts while its story is in active work, and stops as today when it leaves that kind
status: backlog
parent: S-0742
owner: alex
created: 2026-11-12T16:06:57Z
updated: 2026-11-12T16:06:57Z
transitions: []
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_activity.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart_test.go, flai/internal/serve/start_test.go, flai/internal/serve/stop_test.go]
---
# T-3213 A story's agent starts on the pull transition, restarts while its story is in active work, and stops as today when it leaves that kind

## Work

Replace the references to a state's name in the agents' files with reads from the repository's workflow: an agent starts for a story in the state the pull transition leaves from while the active-work states' WIP limit has room (ADR-0043), restarts while its story is in an active-work state (ADR-0108), and is stopped or left to end when the story leaves it. Test a start, a restart, and a stop on a workflow with a state renamed and one added between `in-progress` and `review`.

## Done when

The agents' files name no state, and the tests pass on the baseline and on the other workflow.

## Notes
