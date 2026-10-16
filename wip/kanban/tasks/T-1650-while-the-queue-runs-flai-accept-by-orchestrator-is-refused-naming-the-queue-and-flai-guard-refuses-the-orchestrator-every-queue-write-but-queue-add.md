---
id: T-1650
type: task
nature: feature
title: While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add
status: backlog
parent: S-0403
owner: alex
created: 2026-10-16T05:01:25Z
updated: 2026-10-16T05:01:25Z
transitions: []
stream: S-0403
tags: [cli]
touches: [flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/guard/writes.go, flai/internal/guard/orchestrate_test.go]
---
# T-1650 While the queue runs, flai accept --by orchestrator is refused, naming the queue, and flai guard refuses the orchestrator every queue write but queue_add

## Work

- In `flai accept`, refuse `--by orchestrator` while the `merge-queue` host action is on for the project, unless the queue itself is landing. Name `queue_add` in the refusal.
- In `flai guard`, under `FLAI_ROLE=orchestrate`, refuse `flai queue remove`, `hold`, and `resume`, and their host writes. Allow `flai queue add` and `queue_add`.
- Test both refusals and what stays allowed.
- First in its layer: it waits for no task.

## Done when

The orchestrator can only enqueue while the queue runs, and the tests pass.

## Notes
