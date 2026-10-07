---
id: S-0308
type: story
nature: improvement
title: golangci-lint fails at once when another story's agent is running it on the same host
status: backlog
owner: alex
created: 2026-10-07T01:07:14Z
updated: 2026-10-07T02:19:57Z
transitions: []
tags: []
touches: [flai/.golangci.yaml, flai/tests/integration/golangci_lock_test.go, design/system/devex.md, design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-07T01:07:14Z
  value: 12.5
  by: planner-S-0308
  at: 2026-10-07T02:19:34Z
forecast:
  duration: 10m
  delivery: 2026-10-07T11:10:00Z
  basis: "flai forecast: median 84 s per unit of size over 8 done improvement stories on claude-opus-5-5 in the medium band, times size 7 (2 criteria, 5 touches), 34th in the pull order with an in-progress limit of 3; kept, a one-key config change, one integration test, a design row, and an issue close fit that size."
  by: planner-S-0308
  at: 2026-10-07T02:19:34Z
finalized:
  by: alex
  at: 2026-10-07T02:17:22Z
---
# S-0308 golangci-lint fails at once when another story's agent is running it on the same host

## Goal

This story remediates [I-0101](../../../design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md), "golangci-lint fails at once when another story's agent is running it on the same host". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0101 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0101 is closed with `flai issue close I-0101 --reason` saying what fixed it

## Tasks
- T-1153 golangci-lint allows parallel runners from flai/.golangci.yaml, with an integration test that holds its lock
- T-1155 devex.md says the lint allows parallel runners, and I-0101 is closed with what fixed it

## Notes

Cost of delay inputs set by flai from I-0101. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T00:34:48Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed remedy, from I-0101's one instance: golangci-lint v2 takes a lock on `golangci-lint.lock` in the system temp folder at start, and gives up with exit 3 after 5 seconds unless `run.allow-parallel-runners` is set (golangci-lint `pkg/commands/run.go`, `acquireFileLock`). Every worktree on the host shares that lock. The `golangci-lint` tier in `system-flow.yaml` passes `--allow-parallel-runners`; `scripts/flai-test.sh`, which the close-out runs, does not. Setting `run.allow-parallel-runners: true` in `flai/.golangci.yaml` covers the script, the tier, CI, and a run by hand from one place. Parallel runners, not `allow-serial-runners`, because the tier already runs that way without trouble and serializing would make parallel close-outs wait on each other.

Touches, all files, no folder touch:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/.golangci.yaml` | layout | The fix: `run.allow-parallel-runners: true` |
| `flai/tests/integration/golangci_lock_test.go` | layout | New test: holds the lock under a test `TMPDIR` and runs golangci-lint with the repository's config, expecting no exit 3 |
| `design/system/devex.md` | design | Its "Go build and lint" row describes the lint config and `flai-test.sh` |
| `design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md` | declared by the goal | Closed by `flai issue close` |
| `design/issues/summary.md` | co-change | Regenerated when I-0101 closes |

`flai touches suggest` listed nothing else that bears on the story: its other paths (`design/system/flai-cli.md`, `docs/users/flai.md`, other issues) co-change with `design/issues/summary.md`, which changes with every issue. `scripts/flai-test.sh` and `system-flow.yaml` are left alone: the script reads the config, and the manifest's tier flag stays harmless.

Forecast 10m, delivery 2026-10-07T11:10:00Z: flai's figure (84 s per unit of size over 8 done improvement stories, size 7), kept, since the work is that small.

Cost of delay 12.50 USD a week: `flai cod` from the operator's input (5m lost per 168h cycle at 150 USD an hour), kept. Each collision fails a whole close-out, so the real loss per occurrence is likely more than 5m, but the input is the operator's and there is one instance only.
