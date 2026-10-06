---
id: T-0995
type: task
nature: remediation
title: The host API action item.criteria and the dashboard's HTTP route tick a story's criteria
status: done
parent: S-0282
owner: alex
created: 2026-10-06T03:47:58Z
updated: 2026-10-06T04:17:15Z
transitions:
  - to: ready
    at: 2026-10-06T04:06:49Z
    by: agent-S-0282
  - to: in-progress
    at: 2026-10-06T04:06:49Z
    by: agent-S-0282
  - to: done
    at: 2026-10-06T04:17:15Z
    by: agent-S-0282
stream: S-0282
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, "flaiover/src/routes/api/items/[id]/criteria", flaiover/src/lib/server/agent.ts]
after: [T-0993]
usage:
  source: log
  seconds: 626
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 40
      output: 11908
      cache_read: 1907363
      cache_write: 56139
      cost: 0.9578
---
# T-0995 The host API action item.criteria and the dashboard's HTTP route tick a story's criteria

## Work

Add the host API action `item.criteria` in `flai/internal/hostapi/writes.go`, which runs `flai criteria tick|untick` with the owner as `--by`, `--autocommit`, and the trailer, as `item.finalize` does; and the dashboard's HTTP route `POST /api/items/[id]/criteria` (body `{tick, untick, hash}`) in `flaiover/src/routes/api/items/[id]/criteria/+server.ts`, with `item.criteria` added to `REQUIRED_METHODS` in `flaiover/src/lib/server/agent.ts`.

Waits for the CLI task: the action runs the command it adds.

## Done when

- [ ] `item.criteria` is tested in `writes_test.go`, and the contract test passes.
- [ ] The route is tested in `criteria.test.ts` and `npm test` passes for it.

## Notes
