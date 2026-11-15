---
id: T-3213
type: task
nature: feature
title: A story's agent starts on the pull transition, restarts while its story is in active work, and stops as today when it leaves that kind
status: done
parent: S-0742
owner: alex
created: 2026-11-12T16:06:57Z
updated: 2026-11-15T01:54:57Z
transitions:
  - to: ready
    at: 2026-11-15T01:26:58Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:26:58Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:54:57Z
    by: agent-S-0742
stream: S-0742
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_launch.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_activity.go, flai/internal/serve/restart.go, flai/internal/serve/start.go, flai/internal/serve/stop.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart_test.go]
after: [T-3370]
usage:
  source: log
  seconds: 1679
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 144
      output: 697
      cache_read: 9218142
      cache_write: 363609
      cost: 3.612
---
# T-3213 A story's agent starts on the pull transition, restarts while its story is in active work, and stops as today when it leaves that kind

## Work

Replace the references to a state's name in the agents' files with reads from the repository's workflow: an agent starts for a story in the state the pull transition leaves from while the active-work states' WIP limit has room (ADR-0043), restarts while its story is in an active-work state (ADR-0108), and is stopped or left to end when the story leaves it. Test a start, a restart, and a stop on a workflow with a state renamed and one added between `in-progress` and `review`.

## Done when

The agents' files name no state, and the tests pass on the baseline and on the other workflow.

## Notes
