---
id: T-2599
type: task
nature: feature
title: The host read kluster.log answers a step's log from an offset, and the dashboard requires it
status: done
parent: S-0595
owner: alex
created: 2026-11-07T22:06:35Z
updated: 2026-11-07T22:24:33Z
transitions:
  - to: ready
    at: 2026-11-07T22:09:57Z
    by: agent-S-0595
  - to: in-progress
    at: 2026-11-07T22:10:29Z
    by: agent-S-0595
  - to: done
    at: 2026-11-07T22:24:33Z
    by: agent-S-0595
stream: S-0595
tags: [cli, dashboard]
touches: [flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flaiover/src/lib/server/agent.ts]
usage:
  source: log
  seconds: 844
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 43
      output: 14941
      cache_read: 7303189
      cache_write: 57624
      cost: 2.1352
---
# T-2599 The host read kluster.log answers a step's log from an offset, and the dashboard requires it

## Work

ADR-0263 §4. In `flai/internal/hostapi/kluster.go` add the read `kluster.log`, whatever the host action says, with `{namespace, step, offset}`: it validates the namespace (a DNS label) and the step (a step name: `build`, `build-<name>`, `load`, `deploy`, `teardown`; no path separators), reads `.flai-cache/kluster/logs/<namespace>/<step>.log` from `offset`, at most 64 KB, and answers `namespace`, `step`, `log` (from the root), `offset`, `next`, `size`, `text`, `reset` (true when the file is shorter than `offset`, and the text then starts at 0), `outcome` (the step's outcome in the namespace's record, empty when none), and `exists`. A missing log answers `exists` false and empty text, not an error. Register it in `klusterSpecs()` with `reads: true` and a `describe`. Add `kluster.log` to `REQUIRED_METHODS` in `flaiover/src/lib/server/agent.ts` so `TestTheDashboardAsksForExactlyWhatFlaiOffers` holds. Layer 1: waits for nothing; it reads the record and the logs as they are named today.

## Done when

- Tests in `kluster_test.go` cover a read from 0, a read from the middle, the 64 KB cap with `next`, `reset` after truncation, a missing log, and refusals of a bad namespace or step (`../x`), and fail without the change.
- `flai test` on the changed paths passes, the contract test included.

## Notes
