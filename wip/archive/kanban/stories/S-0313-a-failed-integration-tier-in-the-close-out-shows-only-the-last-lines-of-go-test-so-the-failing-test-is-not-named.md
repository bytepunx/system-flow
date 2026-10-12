---
id: S-0313
type: story
nature: improvement
title: A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named
status: done
owner: alex
created: 2026-10-07T14:26:02Z
updated: 2026-10-12T16:25:54Z
transitions:
  - to: ready
    at: 2026-10-12T15:50:17Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-12T15:51:15Z
    by: agent-S-0313
  - to: review
    at: 2026-10-12T16:02:14Z
    by: agent-S-0313
  - to: done
    at: 2026-10-12T16:25:54Z
    by: alex
tags: []
topics: [testing]
touches: [flai/internal/verify/parse.go, flai/internal/verify/parse_test.go, scripts/integration.sh, system-flow.yaml, scripts/README.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md, design/issues/summary.md, design/issues/I-0114-flai-verify-s-integration-tier-keeps-only-the-last-lines-of-go-test-s-output-so-the-failing-test-is-not-named.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 675
  turns:
    - day: 2026-10-12
      hand_edits: 1
      work: 38
  models:
    - model: claude-opus-5-5
      input: 106
      output: 30728
      cache_read: 6731170
      cache_write: 239011
      cost: 3.7384
  strategic:
    - kind: orchestrator
      seconds: 1329
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 129
          output: 1995
          cache_read: 27532504
          cache_write: 60594
          cost: 6.8923
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4m
    by: flai
    at: 2026-10-07T14:26:02Z
  value: 10
  by: planner-S-0313
  at: 2026-10-07T23:38:54Z
forecast:
  duration: 45m
  delivery: 2026-10-12T19:47:00Z
  basis: "Its own forecast of 45m; 12th in the pull order with an in-progress limit of 5, behind S-0337, S-0233, S-0234, S-0235, S-0236, S-0237, S-0238, S-0239, S-0241, S-0304, S-0305 and S-0306."
  by: flai
  at: 2026-10-12T15:48:40Z
finalized:
  by: orchestrator
  at: 2026-10-07T23:40:20Z
---
# S-0313 A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named

## Goal

This story remediates [I-0113](../../../design/issues/I-0113-a-failed-integration-tier-in-the-close-out-shows-only-the-last-lines-of-go-test-so-the-failing-test-is-not-named.md), "A failed integration tier in the close-out shows only the last lines of go test, so the failing test is not named". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0113 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0113 is closed with `flai issue close I-0113 --reason` saying what fixed it
- [x] I-0114 is closed with `flai issue close I-0114 --reason` saying what fixed it

## Tasks
- T-1286 A plain finding names the failures go test printed, not only the last lines
- T-1287 The integration tier reads go test -json, so its findings name each failing test by file and line
- T-1288 The manifest and settings docs say what a plain finding keeps
- T-1289 Close I-0113 and I-0114 with what fixed them

## Notes

Cost of delay inputs set by flai from I-0113. time_lost_per_cycle 4m: 4m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-07T09:13:42Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0113's one instance (S-0215). The `integration` tier in `system-flow.yaml` runs `scripts/integration.sh`, a plain `go test -race -count=1 ./...`, under `format: plain`. `plain()` in `flai/internal/verify/parse.go` keeps only the last 20 lines, so the failing package, early in the run, is cut off. Two fixes, one per layer-1 task:

- T-1287: the tier runs the script with `-json` under `format: go-test-json`, the parser the `go-test` tier already uses, so each finding names its test, file, and line.
- T-1286: a plain finding keeps go test's `--- FAIL`, `FAIL<tab><package>`, and `panic:` lines ahead of the tail. This covers any plain tier that runs go test, the template's `integration` tier among them.

Touches, all files; no folder touch was kept:

| Touch | From |
|-------|------|
| `flai/internal/verify/parse.go`, `parse_test.go` | layout: `plain()` and its tests |
| `scripts/integration.sh`, `system-flow.yaml` | layout: the tier and its script |
| `scripts/README.md` | layout: its `integration.sh` row |
| `design/system/project-manifest.md` | design and co-change: says plain "reads only the exit status" |
| `docs/operators/settings.md` | design: the `tests[].format` row says the same |
| `docs/users/flai.md` | co-change: `flai test` and `flai verify` describe findings |
| `design/issues/I-0113-…md`, `design/issues/summary.md` | criterion 2: `flai issue close`; shared paths |

`flai touches suggest` declared nothing, so it ran from `parse.go`, `integration.sh`, and `system-flow.yaml`. Its template and CI co-changes were left out, since the template's tiers do not change.

Forecast: `flai forecast` gave 17m (83 s per unit of size over 9 done improvement stories, size 12). It is raised to 45m: T-1287's check and the close-out each run the full `go test -race` suite. The delivery moves by the same 28m, to 2026-10-08T06:37:00Z.

Cost of delay: `flai cod` gives 10.00 USD a week, from flai's input of 4m lost per 168h cycle at 150 USD an hour. It stands as computed: one occurrence gives no reason to change it.
