---
id: T-0880
type: task
nature: improvement
title: flai stats reports the project strategic total per kind beside the per-item figures
status: done
parent: S-0226
owner: alex
created: 2026-10-05T04:44:53Z
updated: 2026-10-06T18:11:56Z
transitions:
  - to: ready
    at: 2026-10-06T18:02:22Z
    by: agent-S-0226
  - to: in-progress
    at: 2026-10-06T18:02:22Z
    by: agent-S-0226
  - to: done
    at: 2026-10-06T18:11:56Z
    by: agent-S-0226
stream: S-0226
tags: [flai]
touches: [flai/internal/metrics, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/metrics.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0878]
usage:
  source: log
  seconds: 574
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 56
      output: 20381
      cache_read: 2756377
      cache_write: 90500
      cost: 1.5467
---
# T-0880 flai stats reports the project strategic total per kind beside the per-item figures

## Work

Report what strategic agents spent on no item, per kind, as T-0878's ADR decides. It waits for T-0878 for where the project total comes from; it shares no path with T-0879 and runs beside it. Should the ADR store the total instead of deriving it from the activity documents, give it an `after` on T-0879, which writes it.

- `flai/internal/metrics`: per strategic kind, the project total (cost, seconds, estimated) of the activities that charged no item, read from the activity documents (`StrategicAgent` and `strategic.go`), so that what was charged to items plus the project total equals the activity document's totals, kind by kind.
- `flai/cmd/stats.go`: print it beside the per-item strategic figures and carry it in `--json`; update the command's help, then `docs/users/flai-reference.md` from it.
- `design/system/metrics.md`: define the project strategic total, how the orchestrator's split charge is read per item and up the hierarchy, and the equality with the activity document's totals.
- `docs/users/flai.md`: describe the orchestrator's charge to the items an activity names, and the project total in `flai stats`.
- Tests in `flai/internal/metrics` and `flai/cmd/check_stats_test.go` pin: an orchestrator activity split over two items and rolled up to their epics; an activity naming none in the project total; per-kind item charges plus the project total equal the activity document's totals.

## Done when

- `flai stats` and `flai stats --json` show each strategic kind's project total beside the per-item figures, and the tests above pass with `scripts/flai-test.sh`.
- `design/system/metrics.md`, `docs/users/flai.md`, and `docs/users/flai-reference.md` say what the total is and where it comes from.
- `flai check --strict` and the markdown lint pass.

## Notes

- Drafted by planner-S-0226.
