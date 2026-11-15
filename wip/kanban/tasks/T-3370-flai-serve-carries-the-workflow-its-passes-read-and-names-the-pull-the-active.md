---
id: T-3370
type: task
nature: feature
title: flai serve carries the workflow its passes read, and names the pull, the active-work room, and the kinds they read it by
status: done
parent: S-0742
owner: alex
created: 2026-11-15T01:18:57Z
updated: 2026-11-15T01:26:38Z
transitions:
  - to: ready
    at: 2026-11-15T01:26:25Z
    by: agent-S-0742
  - to: in-progress
    at: 2026-11-15T01:26:25Z
    by: agent-S-0742
  - to: done
    at: 2026-11-15T01:26:38Z
    by: agent-S-0742
stream: S-0742
tags: []
touches: [flai/internal/serve/serve.go, flai/internal/serve/flow.go, flai/internal/serve/kinds.go, flai/internal/serve/kinds_test.go, flai/internal/serve/agents_launch.go]
usage:
  source: log
  seconds: 13
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 40
      cache_read: 574416
      cache_write: 1365
      cost: 0.217
---
# T-3370 flai serve carries the workflow its passes read, and names the pull, the active-work room, and the kinds they read it by

## Work

- Take S-0743's seam as it is (MS-0567): `Options.Workflow` in `serve.go` and `flow.go` (`Options.flow`, `openRepo`), byte for byte as `story/S-0743` has them, so that the two branches merge cleanly.
- Give the launcher the options' workflow (`launcher.flow`, set in `newLauncher`) and `launcher.open`, which opens a project running it.
- `kinds.go`: the pull transition (`pullTransition`, `isPullable`, `isPulledInto`, `pullQueue`), the kinds (`isActive`, `isWaiting`, `isUnderway`, `isQueued`, `isAccepted`, `isClosed`), and `activeRoom`, the room the WIP limits of the active-work states leave.
- `kinds_test.go`: `kindsFlow`, S-0743's renamed board under a name of its own, and tests of the helpers and the seam against it and the baseline.
- Waits for nothing. The four other tasks wait for it: each reads the workflow through these helpers.

## Done when

- The helpers answer for the baseline and for a board with in-progress renamed doing and testing added after it, as the kinds give.
- `flai test` passes on the paths changed.

## Notes

- S-0743's `flow_test.go` is not copied: it calls `AwaitsAcceptance`, which S-0743 adds to `flai/internal/workflow`. The fixture is copied into `kinds_test.go` as `kindsFlow` instead, so that no names collide when both stories land.
