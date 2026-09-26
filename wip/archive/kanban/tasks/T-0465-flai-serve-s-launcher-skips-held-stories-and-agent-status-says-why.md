---
id: T-0465
type: task
nature: feature
title: flai serve's launcher skips held stories and agent.status says why
status: done
parent: S-0128
owner: alex
created: 2026-09-26T08:05:56Z
updated: 2026-09-26T08:12:56Z
transitions:
  - to: ready
    at: 2026-09-26T08:06:04Z
    by: agent-S-0128
  - to: in-progress
    at: 2026-09-26T08:09:41Z
    by: agent-S-0128
  - to: done
    at: 2026-09-26T08:12:56Z
    by: agent-S-0128
stream: S-0128
tags: []
touches: [flai/internal/serve, flai/cmd]
---
# T-0465 flai serve's launcher skips held stories and agent.status says why

## Work

- The launcher does not start a held story's agent and starts the next ready story in pull order that is not held, within the limit; a story it starts in a look claims its paths for the rest of that look.
- `Activity` reports a held ready story as waiting with its reason, whether or not it ever had an agent.
- `flai serve agent start` still starts a held story's agent.

## Done when

- Behaviour tests in `internal/serve` cover skip-ahead, the held story starting first once clear, and the held activity; `go test ./internal/serve` passes.

## Notes
