---
id: T-0783
type: task
nature: feature
title: flai stats reads threads, the limit, and committed files, and prints the new aggregates
status: done
parent: S-0205
owner: alex
created: 2026-10-03T20:40:06Z
updated: 2026-10-03T21:21:09Z
transitions:
  - to: ready
    at: 2026-10-03T21:10:03Z
    by: agent-S-0205
  - to: in-progress
    at: 2026-10-03T21:10:03Z
    by: agent-S-0205
  - to: done
    at: 2026-10-03T21:21:09Z
    by: agent-S-0205
stream: S-0205
tags: []
touches: [flai/cmd/stats.go, flai/cmd/check_stats_test.go, flai/cmd/hostapi_reads_test.go, flai/internal/hostapi/reads.go, flai/internal/statsread, flai/internal/storygit/committed.go, flai/internal/storygit/committed_test.go, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0778, T-0779, T-0780, T-0781, T-0782]
usage:
  source: log
  seconds: 666
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 137
      output: 57292
      cache_read: 8272516
      cache_write: 223939
      cost: 4.1424
---
# T-0783 flai stats reads threads, the limit, and committed files, and prints the new aggregates

## Work

`flai stats` reads what the new metrics need, the threads, the board's in-progress limit, the manifest's projects, and each story's committed files from git (a reader in `flai/internal/storygit`), passes them to `metrics.Compute`, and prints the new aggregates; `docs/users/flai.md` and the generated reference say so. Waits for every metric task: it wires and prints what they compute.

## Done when

- [ ] `flai stats --json` carries every new section on a project, and `flai stats` prints the aggregates
- [ ] The committed-files reader is tested against a real git repository, and the command against a fixture
- [ ] `docs/users/flai.md` and `docs/users/flai-reference.md` describe the change

## Notes
