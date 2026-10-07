---
id: T-1218
type: task
nature: feature
title: flai stats reports the coordination values under coordination, in JSON and in print
status: backlog
parent: S-0337
owner: alex
created: 2026-10-07T20:17:26Z
updated: 2026-10-07T20:21:18Z
transitions: []
stream: S-0337
tags: [flai]
touches: [flai/internal/metrics/coordination.go, flai/internal/metrics/coordination_test.go, flai/internal/metrics/metrics.go, flai/internal/statsread/statsread.go, flai/internal/statsread/statsread_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go]
after: [T-1217]
---
# T-1218 flai stats reports the coordination values under coordination, in JSON and in print

## Work

Compute and report what T-1216 defines. It waits for T-1217, whose log it reads.

- `statsread.Read` loads the message files and the conflict log into `metrics.Options`, beside the threads and the activities it loads today; one it cannot read is logged and its values left out, as the commits are.
- `flai/internal/metrics/coordination.go` works out each value per week of the window from them and the transitions.
- `metrics.go` adds `coordination` to the stats; `flai stats` prints a Coordination section, and `--json` and the host API's `stats.get` return it.

## Done when

- Tests over fixtures cover each value, an empty window, and a missing log.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
