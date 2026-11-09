---
id: S-0643
type: story
nature: improvement
title: A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make
status: backlog
owner: alex
created: 2026-11-09T21:22:33Z
updated: 2026-11-09T22:33:17Z
transitions: []
tags: [flai, release, formula]
topics: [release, formula]
touches: [flai/cmd/release.go, flai/cmd/release_pending_test.go, flai/internal/release/release.go, flai/internal/release/release_test.go, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, template/CHANGELOG.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/strategic-agents.md, design/issues/I-0320-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 142
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 252
          cache_read: 943474
          cache_write: 89878
          cost: 0.2755
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-09T21:22:33Z
  value: 12.5
  by: planner-S-0643
  at: 2026-11-09T22:14:12Z
forecast:
  duration: 1h
  delivery: 2026-11-10T01:39:00Z
  basis: "Its own forecast of 1h; 13th in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633, S-0634, S-0645, S-0644, S-0638, S-0636, S-0640, S-0625, S-0630, S-0628, S-0631, S-0646, S-0241 and S-0637."
  by: flai
  at: 2026-11-09T22:33:17Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:30:00Z
---
# S-0643 A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make

## Goal

This story remediates [I-0320](../../../design/issues/I-0320-a-publish-instance-s-evaluate-and-publish-steps-race-a-publish-made-meanwhile.md), "A publish instance's evaluate and publish steps race a publish made meanwhile, so its last step finds nothing pending and closes on a release it did not make". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0320 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0320 is closed with `flai issue close I-0320 --reason` saying what fixed it

## Tasks
- T-2818 flai release --pending --require-pending exits 2 when nothing is pending, naming the release that covered it
- T-2819 The publish formula's publish step runs with --require-pending and stops cleanly on exit 2
- T-2820 Document --require-pending and the publish instance that another publish outran
- T-2821 Close I-0320 with what fixed it

## Notes

Cost of delay inputs set by flai from I-0320. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T16:58:32Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed solution, from I-0320's one instance (S-0641). The instance does not hold off a second publish. Instead, its last step tells a publish that another release covered apart from one it made.

- `flai release --pending` gains `--require-pending`. With nothing pending, it exits 2 and names each component's newest release tag, commit, tagger, and time: the release that covered what was pending.
- The publish formula's `publish` step runs with that flag and `stop: [2]` (ADR-0271). So an instance that a publish outran stops cleanly, and its record names the release that covered it.
- Plain `--pending` keeps exit 0 with `nothing pending`, so the dashboard's Publish, `release_publish`, and the operator's scripts are unchanged.
- Not proposed: holding the operator's hand publish back while an instance runs. The race does no harm once the step can tell it apart, and the operator's own shell is never held today (S-0423).

Touches, all files and no folders:

| Touch | Source |
|-------|--------|
| `flai/cmd/release.go`, `flai/cmd/release_pending_test.go` | layout: `publishPending` and its `nothing pending` result; the test file holds `TestReleasePendingNothingPending` |
| `flai/internal/release/release.go`, `flai/internal/release/release_test.go` | layout: `CurrentVersion` reads the release tags; no lookup gives a tag's commit, tagger, and time yet |
| `design/workflows/formulas/publish.md`, `template/root/design/workflows/formulas/publish.md` | design: the formula I-0320 names, and its template copy |
| `design/workflows/formulas/README.md`, `template/root/design/workflows/formulas/README.md` | design: they name the evaluate step's exit 2 |
| `template/CHANGELOG.md` | layout: a template formula change adds an Unreleased line, as S-0621 did |
| `docs/users/flai.md`, `design/system/flai-cli.md`, `design/system/strategic-agents.md` | design: they describe `--pending` and the publish instance's steps |
| `docs/users/flai-reference.md`, `docs/operators/settings.md` | layout: regenerated from the help text by `make flai-reference` |
| `design/issues/I-0320-…md`, `design/issues/summary.md` | criteria: the issue is closed |

`flai touches suggest` had nothing to start from, since the story declared no touches. The explorer predicted the touches from the code layout. The dashboard's publish (`flai/internal/hostapi/writes.go`), `release_publish`, and `flai/internal/preview/push.go` keep plain `--pending`, so the story does not touch them.

No `after`. S-0625, which is ready, changes the same publish formula and `flai/cmd/release.go`. This story reads none of its output, and a trial merge judges the overlap (ADR-0230).

Figures:

- Forecast 1h, delivery 2026-11-10T02:00:00Z. `flai forecast` gave 30m from 98 s per unit of size over 142 done improvement stories, times size 18. I doubled it for the new git read of a tag's tagger and time, a test that reproduces the race, two copies of the formula, and a reference regenerated in the close-out's integration tier.
- Cost of delay 12.50 USD a week, as `flai cod` worked it from the 5m time-lost input. It stands as computed.

Tasks, in three layers:

1. T-2818, the flag, the lookup, and the reproducing test.
2. T-2819, the formula and its changelog line, and T-2820, the documents and the regenerated reference. Both wait for T-2818 and share no path.
3. T-2821 closes I-0320, after T-2819 and T-2820.
