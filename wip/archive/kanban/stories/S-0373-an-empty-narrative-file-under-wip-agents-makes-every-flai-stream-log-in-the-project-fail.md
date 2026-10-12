---
id: S-0373
type: story
nature: remediation
title: An empty narrative file under wip/agents makes every flai stream log in the project fail
status: done
owner: alex
created: 2026-10-12T16:52:52Z
updated: 2026-10-12T18:29:45Z
transitions:
  - to: ready
    at: 2026-10-12T17:30:12Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T18:13:43Z
    by: agent-S-0373
  - to: review
    at: 2026-10-12T18:28:48Z
    by: agent-S-0373
  - to: done
    at: 2026-10-12T18:29:45Z
    by: orchestrator
tags: [cli]
topics: [narratives]
touches: [flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, design/issues/summary.md, flai/internal/workitem/activity_test.go, design/issues/I-0136-an-acceptance-on-main-during-a-close-out-s-tiers-fails-its-final-sync-check-so-the-whole-close-out-runs-again.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 925
  turns:
    - day: 2026-10-12
      ceremony: 3
      hand_edits: 2
      work: 40
  models:
    - model: claude-opus-5-5
      input: 92
      output: 17578
      cache_read: 6805597
      cache_write: 197430
      cost: 3.2925
  strategic:
    - kind: orchestrator
      seconds: 320
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 32
          output: 651
          cache_read: 10445393
          cache_write: 15158
          cost: 2.637
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-12T16:52:52Z
  value: 12.5
  by: planner-S-0373
  at: 2026-10-12T17:29:22Z
forecast:
  duration: 30m
  delivery: 2026-10-12T18:35:00Z
  basis: "Its own forecast of 30m; 3rd in the pull order with an in-progress limit of 5, behind S-0363 and S-0370."
  by: flai
  at: 2026-10-12T18:00:51Z
finalized:
  by: orchestrator
  at: 2026-10-12T17:29:56Z
---
# S-0373 An empty narrative file under wip/agents makes every flai stream log in the project fail

## Goal

This story remediates [I-0129](../../../design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md), "An empty narrative file under wip/agents makes every flai stream log in the project fail". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0129 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0129 is closed with `flai issue close I-0129 --reason` saying what fixed it

## Tasks
- T-1489 WriteIndex lists an unreadable narrative as unreadable and writes the index
- T-1490 flai check counts an unreadable narrative as listed in the index and tests narrative.front-matter
- T-1491 Design and user docs say an unreadable narrative is listed as unreadable instead of failing the write
- T-1492 Close I-0129 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0129. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-08T20:55:33Z, 3.8 days before this story; under one cycle counts as one).

### Planning

Cause, found by the explorer: `Repo.ActiveNarratives` (`flai/internal/workitem/narrative.go`) returns at the first `S-*.md` that does not parse, so `WriteIndex` fails. Every command that rewrites `wip/agents/index.md` fails with it: stream open, log, and state, move, edit, archive, accept, story start, and task done. Each fails after its own change is saved.

Proposed fix: `WriteIndex` lists such a file as an unreadable row, as it already lists an unreadable activity document, and writes the index. `flai check` keeps `narrative.front-matter` as the one finding on it. A story's own unreadable narrative still fails its own log.

Touches, all files, no folder kept:

| Touch | Source |
|-------|--------|
| `flai/internal/workitem/narrative.go`, `narrative_test.go` | layout: the cause and its test (T-1489) |
| `flai/internal/check/check.go`, `check_test.go` | layout and co-change, 38% with the code above (T-1490) |
| `design/system/agent-narrative.md` | design, and co-change at 10% (T-1491) |
| `design/system/flai-cli.md`, `docs/users/flai.md` | co-change at 34% and 29%; both describe the index (T-1491) |
| `design/issues/I-0129-…md`, `design/issues/summary.md` | criterion 2 (T-1492) |

`flai/cmd` is left out on purpose. S-0367 claims that folder, and a test in `flai/internal/workitem` reproduces the failure at `WriteIndex`, the point every command shares.

Figures:

- Forecast 30m: flai's 19m, raised for four tasks in three layers across two packages and three documents, plus the close-out.
- Cost of delay 12.50 USD a week: flai cod's figure from flai's 5m input, unchanged.

### Accepted by the orchestrator

- Verified: ca5a02ce704f27a31e65056e7d0a896859e643bf
- At: 2026-10-12T18:29:45Z

Verdict: pass at ca5a02ce704f27a31e65056e7d0a896859e643bf; flai verify passed every step at that commit, and the verifier matched both criteria to the diff.

- 1: flai/internal/workitem/narrative.go, flai/internal/workitem/narrative_test.go, flai/internal/workitem/activity_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, design/system/agent-narrative.md, design/system/flai-cli.md, docs/users/flai.md
- 2: design/issues/I-0129-an-empty-narrative-file-under-wip-agents-makes-every-flai-stream-log-in-the-project-fail.md, design/issues/summary.md
