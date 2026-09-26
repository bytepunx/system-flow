---
id: T-0430
type: task
nature: improvement
title: A retry on a ready story with a full in-progress limit is queued, not refused
status: done
parent: S-0118
owner: arobson
created: 2026-09-26T03:02:39Z
updated: 2026-09-26T03:13:33Z
transitions:
  - to: ready
    at: 2026-09-26T03:02:43Z
    by: agent-S-0118
  - to: in-progress
    at: 2026-09-26T03:13:33Z
    by: agent-S-0118
  - to: done
    at: 2026-09-26T03:13:33Z
    by: agent-S-0118
stream: S-0118
tags: []
touches: [flai/cmd, flai/internal/serve]
---
# T-0430 A retry on a ready story with a full in-progress limit is queued, not refused

## Work

This task follows the designer's answer to TH-0011. In `flai/internal/serve`, `Restart` for a story in ready with no room under the in-progress limit marks the failed run as queued instead of refusing. The launcher's look then starts it as soon as there is room, the same way it starts a ready story. `Activity` and `agent.status` say the run is queued. The pane shows that and offers no Retry. `flai serve agent restart` prints what happened.

## Done when

- The Go tests show that a queued retry is recorded, is not started while there is no room, and is started by the launcher once there is.
- The Go tests show that a retry with room, or on a story in progress, still starts at once.
- The pane test shows the queued state with no Retry button.
- `make flai-test` and the flaiover tests pass.

## Notes
