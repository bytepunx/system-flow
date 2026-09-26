---
id: T-0464
type: task
nature: feature
title: Compute each ready story's hold from claims, and mark held cards on the board
status: done
parent: S-0128
owner: alex
created: 2026-09-26T08:05:56Z
updated: 2026-09-26T08:09:40Z
transitions:
  - to: ready
    at: 2026-09-26T08:06:04Z
    by: agent-S-0128
  - to: in-progress
    at: 2026-09-26T08:06:04Z
    by: agent-S-0128
  - to: done
    at: 2026-09-26T08:09:40Z
    by: agent-S-0128
stream: S-0128
tags: []
touches: [flai/internal/workitem, flai/internal/hostapi, flai/internal/mcpserver, flai/cmd]
---
# T-0464 Compute each ready story's hold from claims, and mark held cards on the board

## Work

- A story's claim: its `touches` and those of its tasks that are not done or cancelled, a sub-project's name or tag from `system-flow.yaml` read as its path.
- Overlap by the `/`-prefix rule `flai check` uses; an empty claim overlaps every claim.
- A hold reason per ready story (`held (overlap): …`, `held (no-touches): …`) against the stories in progress or in review, with room for a second reason code (`after`).
- `NewBoardView` marks held cards (`held`, `held_reason`); `flai board` prints them; the MCP `board` and `inbox` ready lists carry them.

## Done when

- Behaviour tests cover each claim, component, overlap, and empty-claim case, and the board's held marks; `go test ./internal/workitem ./cmd ./internal/mcpserver` passes.

## Notes
