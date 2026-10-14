---
id: T-1575
type: task
nature: remediation
title: The host write stream.restore runs flai stream restore, and accept.run takes keep_removals, each a path and a reason, for the operator's acceptance only
status: done
parent: S-0371
owner: alex
created: 2026-10-14T20:59:28Z
updated: 2026-10-14T22:30:43Z
transitions:
  - to: ready
    at: 2026-10-14T22:14:49Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T22:14:50Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T22:30:43Z
    by: agent-S-0371
stream: S-0371
tags: [cli, go]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1574, T-1479]
usage:
  source: log
  seconds: 953
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 57
      output: 21111
      cache_read: 3330735
      cache_write: 156117
      cost: 1.9458
---
# T-1575 The host write stream.restore runs flai stream restore, and accept.run takes keep_removals, each a path and a reason, for the operator's acceptance only

## Work

It waits for T-1574 and T-1479, whose command and flags the two writes call, so that the dashboard story after S-0371 has both writes to build on.

- In `flai/internal/hostapi/writes.go`, a new write `stream.restore` `{id, paths, request_id}` in `itemSpecs`, beside `stream.state`, that runs `flai stream restore <id> --path=<p>… --json`, checked as the other writes check their input (`needID`, each path relative, inside the checkout, and never a flag). flai's refusals pass through as `Rule`, as `item.move`'s do.
- In the same file, `accept.run` takes `keep_removals`, a list of `{path, reason}`, each reason not blank, and passes them as `--keep-removal` and `--reason`. It refuses `keep_removals` with `by: orchestrator`, as flai does, so the orchestrator never keeps a removal. Its `exits: {4: Refused}` already passes T-1479's refusal data through.
- The host read `accept.preview` needs no change: it calls `preview.Accept` in process, which gives T-1479's `undoes_main`.
- In `flai/internal/hostapi/writes_test.go`, `TestEveryWriteIsCovered` needs a `good` and a `refused` entry for `stream.restore`; add `keep_removals` to `accept.run`'s. Refused input: a path that is a flag, an absolute path, a `..` path, a blank reason, and `keep_removals` with `by: orchestrator`.
- In `flaiover/src/lib/server/agent.ts`, add `stream.restore` to `REQUIRED_METHODS`: `TestTheDashboardAsksForExactlyWhatFlaiOffers` in `flai/internal/hostapi/contract_test.go` fails while flai offers a method the dashboard does not list. Update `agent.test.ts` too if it lists the methods.

## Done when

- `stream.restore` and `accept.run` with `keep_removals` build the arguments above, and refuse the bad input above, in `writes_test.go`.
- `TestTheDashboardAsksForExactlyWhatFlaiOffers` passes.
- `flai test` on the three files passes.

## Notes

The docs task adds `stream.restore` and `keep_removals` to the `flai hostapi` row of `design/system/flai-cli.md` and to `docs/users/flai.md`'s `flai hostapi` section. The dashboard's route and the review page's cards are the dashboard story's, which also updates `design/system/flaiover-dashboard.md` § Writes.
