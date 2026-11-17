---
id: T-3568
type: task
nature: remediation
title: The move route judges a cancellation's preview and an acceptance's timeout by the target lane's kind
status: done
parent: S-0838
owner: alex
created: 2026-11-17T12:29:28Z
updated: 2026-11-17T12:42:04Z
transitions:
  - to: ready
    at: 2026-11-17T12:36:44Z
    by: agent-S-0838
  - to: in-progress
    at: 2026-11-17T12:36:44Z
    by: agent-S-0838
  - to: done
    at: 2026-11-17T12:42:04Z
    by: agent-S-0838
stream: S-0838
tags: [flaiover, workflow]
touches: ["flaiover/src/routes/api/items/[id]/move/+server.ts", "flaiover/src/routes/api/items/[id]/move/kinds.test.ts", flaiover/src/lib/lanes.fixture.ts]
usage:
  source: log
  seconds: 320
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 11572
      cache_read: 2508202
      cache_write: 94171
      cost: 1.3637
---
# T-3568 The move route judges a cancellation's preview and an acceptance's timeout by the target lane's kind

## Work

In `flaiover/src/routes/api/items/[id]/move/+server.ts`, read the board's flow, its `lanes` and `transitions`, from the repo's board view (`repo().boardView()`, kept until flai says a file changed, as `board()` in `flaiover/src/lib/server/board.ts` reads it). Replace `body.to === 'cancelled'`, which sends a dry run to `item.move.preview`, with `isCancelled(flow, body.to)`, and `body.to === 'done'`, which gives an acceptance its 600000 ms timeout, with `isAccepted(flow, body.to)`, both from `$lib/workflow`. A dry run to a lane of any other kind stays a plain move, as it is today.

Add to `flaiover/src/lib/lanes.fixture.ts` a flow whose lanes of kind `accepted` and `cancelled` have other names (such as `shipped` and `dropped`) and whose pullable queue lane is not `ready` (such as `queued`). This task's test and the two tasks after it use it.

Write `flaiover/src/routes/api/items/[id]/move/kinds.test.ts`: with a repo whose board view answers that flow, a dry-run move to the cancelled lane answers `item.move.preview`, a move to the accepted lane is written with the 600000 ms timeout, and a move to `done` or `cancelled`, which name no lane of that kind there, is a plain move. This waits for no task: it is the first layer.

## Done when

- The route compares no state name: `grep -nE "'(done|cancelled)'" "flaiover/src/routes/api/items/[id]/move/+server.ts"` finds nothing
- `kinds.test.ts` passes, and `designer.test.ts` in the same folder still passes on the baseline
- `flai test` on the task's paths passes

## Notes
