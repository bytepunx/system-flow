---
id: S-0308
type: story
nature: improvement
title: golangci-lint fails at once when another story's agent is running it on the same host
status: done
owner: alex
created: 2026-10-07T01:07:14Z
updated: 2026-10-07T08:48:14Z
transitions:
  - to: ready
    at: 2026-10-07T08:20:40Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-07T08:42:30Z
    by: system-flow
  - to: review
    at: 2026-10-07T08:47:44Z
    by: agent-S-0308
  - to: done
    at: 2026-10-07T08:48:14Z
    by: orchestrator
tags: []
touches: [flai/.golangci.yaml, flai/tests/integration/golangci_lock_test.go, design/system/devex.md, design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 326
  models:
    - model: claude-opus-5-5
      input: 42
      output: 9449
      cache_read: 1441588
      cache_write: 76560
      cost: 1.0899
  strategic:
    - kind: orchestrator
      seconds: 852
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 64
          output: 1009
          cache_read: 4879700
          cache_write: 24788
          cost: 1.2797
        - model: claude-sonnet-5-5
          input: 6
          output: 34
          cache_read: 44347
          cache_write: 25509
          cost: 0.0469
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
  delivery: 2026-10-07T15:13:00Z
  basis: "Its own forecast of 10m; 28th in the pull order with an in-progress limit of 3, behind S-0275, S-0214, S-0215, S-0216, S-0232, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0246, S-0265, S-0279, S-0280, S-0287, S-0288, S-0289, S-0290, S-0291, S-0293, S-0297, S-0298, S-0304, S-0305 and S-0306."
  by: flai
  at: 2026-10-07T08:17:59Z
finalized:
  by: alex
  at: 2026-10-07T02:17:22Z
---
# S-0308 golangci-lint fails at once when another story's agent is running it on the same host

## Goal

This story remediates [I-0101](../../../design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md), "golangci-lint fails at once when another story's agent is running it on the same host". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0101 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0101 is closed with `flai issue close I-0101 --reason` saying what fixed it

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

### Accepted by the orchestrator

- Verified: 720299aa38b4b8f8f98769de2ec85bbae55a29f8
- At: 2026-10-07T08:48:14Z

Verdict: pass. Criteria 1 and 2 are met, the diff is within the story's touches, and there are no convention breaches, verified at commit 720299aa38b4b8f8f98769de2ec85bbae55a29f8.

- 1: flai/.golangci.yaml, flai/tests/integration/golangci_lock_test.go, design/system/devex.md
- 2: design/issues/I-0101-golangci-lint-fails-at-once-when-another-story-s-agent-is-running-it-on-the-same-host.md, design/issues/summary.md
