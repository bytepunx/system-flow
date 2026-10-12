---
id: T-1218
type: task
nature: feature
title: flai stats reports the coordination values under coordination, in JSON and in print
status: done
parent: S-0337
owner: alex
created: 2026-10-07T20:17:26Z
updated: 2026-10-12T16:08:10Z
transitions:
  - to: ready
    at: 2026-10-12T15:35:42Z
    by: agent-S-0337
  - to: in-progress
    at: 2026-10-12T15:48:17Z
    by: agent-S-0337
  - to: done
    at: 2026-10-12T16:08:10Z
    by: agent-S-0337
stream: S-0337
tags: [flai]
touches: [flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1216, T-1217]
usage:
  source: log
  seconds: 1193
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 117
      output: 39684
      cache_read: 7384878
      cache_write: 240616
      cost: 3.8739
---
# T-1218 flai stats reports the coordination values under coordination, in JSON and in print

## Work

Compute and report what T-1216 defines. It waits for T-1216, whose definitions it computes, and T-1217, whose log it reads.

- `statsread.Read` loads the message files and the conflict log into `metrics.Options`, beside the threads and the activities it loads today; one it cannot read is logged and its values left out, as the commits are.
- `flai/internal/metrics/coordination.go` works out each value per week of the window from them and the transitions.
- `metrics.go` adds `coordination` to the stats; `flai stats` prints a Coordination section, and `--json` and the host API's `stats.get` return it.

## Done when

- Tests over fixtures cover each value, an empty window, and a missing log.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
