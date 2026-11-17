---
id: S-0797
type: story
nature: remediation
title: The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed
status: review
owner: alex
created: 2026-11-16T00:41:48Z
updated: 2026-11-16T04:06:58Z
transitions:
  - to: ready
    at: 2026-11-16T01:45:13Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-16T02:59:50Z
    by: agent-S-0797
  - to: review
    at: 2026-11-16T04:06:58Z
    by: agent-S-0797
tags: [merge-queue, flai]
topics: [merge-queue, release]
touches: [design/adrs/README.md, flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/mergequeue/preflight.go, flai/internal/mergequeue/preflight_test.go, flai/internal/mergequeue/state.go, flai/internal/mergequeue/state_test.go, flai/internal/serve/mergequeue.go, flai/internal/serve/mergequeue_test.go, flai/cmd/queue.go, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, design/adrs/0311-the-merge-queue-s-preflight-checks-main-again-with-a-flai-built-from-main-when.md, flai/internal/mergequeue/batch_test.go, flai/internal/mergequeue/land_test.go, flai/internal/check/output.go, flai/internal/check/output_test.go, flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/queue_held_test.go, design/issues/I-0382-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4074
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 6
      work: 47
  models:
    - model: claude-opus-5-5
      input: 330
      output: 120842
      cache_read: 20194787
      cache_write: 775644
      cost: 11.2367
  strategic:
    - kind: planner
      seconds: 36
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 23
          output: 2864
          cache_read: 1810228
          cache_write: 37898
          cost: 0.681
    - kind: orchestrator
      seconds: 176
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 369
          cache_read: 1801001
          cache_write: 249750
          cost: 0.6245
        - model: claude-sonnet-5-5
          input: 4
          output: 49
          cache_read: 13603
          cache_write: 29636
          cost: 0.0345
cost_of_delay:
  inputs:
    time_lost_per_cycle: 3h
    by: flai
    at: 2026-11-16T00:41:48Z
  value: 450
  by: planner-S-0797
  at: 2026-11-16T01:43:27Z
forecast:
  duration: 1h30m
  delivery: 2026-11-16T04:55:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, behind S-0750 and S-0799."
  by: flai
  at: 2026-11-16T02:59:13Z
finalized:
  by: orchestrator
  at: 2026-11-16T01:44:45Z
---
# S-0797 The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed

## Goal

This story remediates [I-0382](../../../design/issues/I-0382-the-merge-queue-s-preflight-holds-the-queue-on-main-s-files-when-the-installed.md), "The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed". The issue recommends this solution:

Proposed, for the story to settle in an ADR: the preflight tells a finding the installed flai cannot judge from one in the files. When main's own flai sub-project is ahead of the running flai's commit, the preflight either runs its `flai check` with a flai built from main, as `scripts/flai.sh` does, or holds with words that name the cause, the installed flai's version against main's unreleased commits, and the remedy, publish and install the release, in place of telling the operator to fix files that are right. Consider also holding an acceptance that changes a format the running flai reads (role files, the workflow file, the manifest) until a flai that reads it is installed.

## Acceptance criteria
- [x] The cause I-0382 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0382 is closed with `flai issue close I-0382 --reason` saying what fixed it

## Tasks
- T-3403 Write ADR-0311 refining ADR-0234 and ADR-0201: the preflight checks main again with a flai built from main when the running flai is behind, else holds naming the release to install
- T-3404 Tell whether the running flai is behind main's flai sub-project, by its build commit or else by the newest flai release tag
- T-3406 The preflight checks main again with the manifest's flai.command when the running flai is behind, and records the cause on a hold it cannot clear that way
- T-3408 Document the preflight's second check and the hold that names the installed flai and the release to install
- T-3411 The hold thread and flai queue name the installed flai, main's unreleased flai commits, and the release to install, and I-0382 is closed

## Notes

Cost of delay inputs set by flai from I-0382. time_lost_per_cycle 3h carried over from I-0382's Impact section. Its Impact time_lost_per_cycle was taken rather than the 3h derived from its cost and count.

### Planning

Touches, by where each came from:

| Touch | Source |
|-------|--------|
| `flai/internal/mergequeue/preflight.go`, `preflight_test.go` | design: `preflightCheck` runs the running flai's `check.Run`, ADR-0234 |
| `flai/internal/mergequeue/state.go`, `state_test.go` | layout: `PreflightHold` gains the cause |
| `flai/internal/release/outdated.go`, `outdated_test.go` | layout: beside `FlaiOutdated`, which compares release tags only |
| `flai/internal/serve/mergequeue.go`, `mergequeue_test.go` | design: `holdText` tells the operator to fix main's files |
| `flai/cmd/queue.go`, `queue_test.go` | layout: `printQueue` prints the preflight's hold |
| `design/adrs/0309-...md`, `design/adrs/README.md` | goal: the issue asks for an ADR; 0308 is S-0796's, so the number may move |
| `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/operators/index.md` | design: the sections that describe the preflight's hold |

- `flai touches suggest` declared nothing at first. From the paths above it listed co-change files: `docs/users/flai-reference.md`, `design/system/metrics.md`, `template/CHANGELOG.md`, and others. None was kept. No flag or command changes. The `preflight_held` event keeps its fields. The template is not changed.
- No folder touch.
- No `after`: the story reads nothing another open story writes.
- `cmd/board.go`'s `preflightSummary` and the dashboard's merge queue panel are left as they are. They show the hold and its findings, and the thread carries the cause.

Figures:

- Forecast 1h30m, delivered by 05:23Z. flai's 33m is the median over large remediation stories, at size 17. It was raised for an ADR, a git comparison, a second check through `flai.command` in a scratch worktree, the hold's words, three documents, and a close-out of about 10m. S-0795 took 54m against a 40m forecast. The delivery moved by the same 57m from flai's 04:26.
- Cost of delay 450 USD a week, as `flai cod` gives it from the inputs flai set from I-0382: 3h lost per cycle at 150 USD an hour, one cycle a week. It stands: the hold stopped every landing for about 3h, and the next format change will cause it again.
