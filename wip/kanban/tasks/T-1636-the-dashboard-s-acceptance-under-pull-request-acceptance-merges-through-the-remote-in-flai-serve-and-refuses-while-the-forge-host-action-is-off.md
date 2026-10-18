---
id: T-1636
type: task
nature: feature
title: The dashboard's acceptance under pull-request acceptance merges through the remote in flai serve, and refuses while the forge host action is off
status: done
parent: S-0397
owner: alex
created: 2026-10-16T04:40:11Z
updated: 2026-10-18T03:06:58Z
transitions:
  - to: ready
    at: 2026-10-18T02:42:39Z
    by: agent-S-0397
  - to: in-progress
    at: 2026-10-18T02:42:40Z
    by: agent-S-0397
  - to: done
    at: 2026-10-18T03:06:58Z
    by: agent-S-0397
stream: S-0397
tags: [cli, dashboard]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-1634]
usage:
  source: log
  seconds: 1441
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 35
      output: 189
      cache_read: 1180229
      cache_write: 68313
      cost: 0.5427
---
# T-1636 The dashboard's acceptance under pull-request acceptance merges through the remote in flai serve, and refuses while the forge host action is off

## Work

- `accept.run` in `flai/internal/hostapi/writes.go` runs `flai accept` as a child in flai serve's environment, where the token's variable is at hand. So the child's merge is T-1635's direct merge, in flai serve's process tree, and never a hand-off.
- For a project under pull-request acceptance, refuse `accept.run` while the `forge` host action is off, before the child runs, naming `flai serve enable forge` as the other host actions' refusals do (`notEnabled`). Leave a project under `local` acceptance as it is.
- Keep `orchestratorAccepts` as S-0395 left it: the gate first.
- Test: a pull-request project with the host action off is refused and runs no child; with it on, the child runs with serve's environment; a local project is unaffected.
- **Waits for T-1634 only.** It runs beside T-1635: they share no path, and the refusal needs nothing T-1635 adds.

## Done when

A dashboard acceptance under pull-request acceptance is refused with the host action off and runs `flai accept` with it on. The hostapi tests pass.

## Notes
