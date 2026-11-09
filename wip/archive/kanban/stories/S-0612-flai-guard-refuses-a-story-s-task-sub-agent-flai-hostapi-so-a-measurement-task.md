---
id: S-0612
type: story
nature: improvement
title: flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing
status: done
owner: alex
created: 2026-11-07T23:38:52Z
updated: 2026-11-09T09:32:19Z
transitions:
  - to: ready
    at: 2026-11-09T08:03:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T08:03:52Z
    by: agent-S-0612
  - to: review
    at: 2026-11-09T09:29:54Z
    by: agent-S-0612
  - to: done
    at: 2026-11-09T09:32:19Z
    by: orchestrator
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/hostapi_test.go, flai/cmd/guard_test.go, flai/internal/hostapi/guardreads_test.go, docs/users/flai.md, design/system/flai-cli.md, design/issues/I-0307-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 5176
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 1
      hand_edits: 2
      work: 62
  models:
    - model: claude-opus-5-5
      input: 274
      output: 52629
      cache_read: 35441106
      cache_write: 981908
      cost: 13.9786
  strategic:
    - kind: planner
      seconds: 384
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 64
          output: 18694
          cache_read: 4425041
          cache_write: 175928
          cost: 2.6666
    - kind: orchestrator
      seconds: 296
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 48
          output: 348
          cache_read: 1833117
          cache_write: 239553
          cost: 0.5372
        - model: claude-sonnet-5-5
          input: 7
          output: 377
          cache_read: 94420
          cache_write: 22347
          cost: 0.0864
cost_of_delay:
  inputs:
    time_lost_per_cycle: 20m
    by: flai
    at: 2026-11-07T23:38:52Z
  value: 50
  by: planner-S-0612
  at: 2026-11-09T07:59:19Z
forecast:
  duration: 45m
  delivery: 2026-11-09T08:49:00Z
  basis: "flai forecast's 18m (103 s per unit of size over 25 done medium improvement stories, size 10) raised to 45m for S-0568's like guard-read fix, which took 1927 s, and for classing every host API method in T-2703; delivery is flai's 08:22 moved by the 27m added."
  by: planner-S-0612
  at: 2026-11-09T07:59:19Z
finalized:
  by: orchestrator
  at: 2026-11-09T08:03:12Z
---
# S-0612 flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing

## Goal

This story remediates [I-0307](../../../design/issues/I-0307-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md), "flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0307 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0307 is closed with `flai issue close I-0307 --reason` saying what fixed it

## Tasks
- T-2702 flai guard lets every agent run flai hostapi with no method or with a read method, and refuses the write methods
- T-2703 A test fails when the host API gains a method the guard's hostapi read list does not class
- T-2704 Document that every agent may run flai hostapi's read methods and that its write methods are the story's agent's
- T-2705 Close I-0307 with the guard's hostapi read rule as its fix

## Notes

Cost of delay inputs set by flai from I-0307. time_lost_per_cycle 20m: 20m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T22:43:39Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0307's one instance: `flai guard` treats `flai hostapi` as a read when it names no method or a method on a list of the host API's read methods kept in the guard, and refuses every other method. The guard cannot import `flai/internal/hostapi`, which imports it, so a test in `flai/internal/hostapi` keeps the list in step with the method table. This follows S-0568, which made `flai role list`, `show`, and `definition` reads the same way. No ADR: allowing a read is within ADR-0060. The `git clone` refusal the instance also names is out of scope: the issue is about `flai hostapi`.

Where each touch came from. The story declared none, so `flai touches suggest` ran from `flai/internal/guard/guard.go`, the file I-0307 names, and `flai/cmd/hostapi.go`.

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/guard/guard.go` | design (I-0307) | `cliReads` and `flagReads` gain `hostapi` and its read list |
| `flai/internal/guard/hostapi_test.go` | layout | new test file reproducing I-0307, as `roles_test.go` did for S-0568 |
| `flai/cmd/guard_test.go` | co-change (31%) | S-0568 added its case here |
| `flai/internal/hostapi/guardreads_test.go` | layout | new test that ties the guard's list to the method table |
| `docs/users/flai.md` | co-change (9%), design | § Request timing names `flai hostapi --timing` |
| `design/system/flai-cli.md` | co-change (7%), design | the `flai guard` row lists the reads |
| `design/issues/I-0307-…md` | criteria | criterion 2 closes it |
| `design/issues/summary.md` | criteria | `flai issue close` rewrites it |

Left out: `flai/cmd/hostapi.go`, whose command does not change, and `docs/users/flai-reference.md`, generated from the help, which does not change either. No folder touch, and no `after`: the story reads nothing another open story writes.

Figures:

- Forecast 45m, not flai's 18m. flai's figure is 103 s per unit of size over 25 done medium improvement stories, times size 10. S-0568, the same kind of guard-read fix, took 1927 s of agent time, and T-2703 must class every method of the host API's table. Delivery is flai's 08:22 moved by the 27 minutes added.
- Cost of delay 50 USD a week, as `flai cod` works it out from the operator's input, 20m lost per 168h cycle at 150 USD an hour. It stands as computed.

### Accepted by the orchestrator

- Verified: 8f03c1d62c36c4303e66f541056962c992342f79
- At: 2026-11-09T09:32:19Z

Verdict: pass at 8f03c1d62c36c4303e66f541056962c992342f79: flai guard lets every agent run flai hostapi alone or with a read method, --timing included, and refuses the write methods by name; I-0307 is closed.
- 1: flai/internal/guard/guard.go, flai/internal/guard/hostapi_test.go, flai/cmd/guard_test.go, flai/internal/hostapi/guardreads_test.go, docs/users/flai.md, design/system/flai-cli.md
- 2: design/issues/I-0307-flai-guard-refuses-a-story-s-task-sub-agent-flai-hostapi-so-a-measurement-task.md, design/issues/summary.md
