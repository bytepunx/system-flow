---
id: S-0634
type: story
nature: improvement
title: flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier
status: done
owner: alex
created: 2026-11-07T23:41:12Z
updated: 2026-11-10T05:44:41Z
transitions:
  - to: ready
    at: 2026-11-09T22:01:06Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T22:01:11Z
    by: agent-S-0634
  - to: review
    at: 2026-11-09T22:49:32Z
    by: agent-S-0634
  - to: in-progress
    at: 2026-11-10T05:13:22Z
    by: flai
  - to: review
    at: 2026-11-10T05:44:26Z
    by: agent-S-0634
  - to: done
    at: 2026-11-10T05:44:41Z
    by: alex
tags: []
topics: [testing, flai, docs]
touches: [flai/cmd/reference_test.go, system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md, design/issues/I-0308-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md, design/issues/summary.md, design/issues/I-0321-testtcpanswersalistenerandnotaclosedone-fails-on-every-run-on-the-operator-s.md, design/issues/I-0309-flai-cmd-s-host-and-mcp-over-http-timing-tests-fail-under-a-host-load-of-about.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, design/issues/I-0329-teststartdevserversayswhatavitethatdoesnotanswer200answered-and.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4764
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 2
      hand_edits: 1
      work: 90
    - day: 2026-11-10
      ceremony: 3
      hand_edits: 1
      work: 28
  models:
    - model: claude-opus-5-5
      input: 260
      output: 28425
      cache_read: 36602254
      cache_write: 638273
      cost: 12.9962
  strategic:
    - kind: planner
      seconds: 292
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 226
          output: 763
          cache_read: 11373243
          cache_write: 642091
          cost: 0.799
        - model: claude-opus-5-5
          input: 471
          output: 31867
          cache_read: 25913297
          cache_write: 847628
          cost: 8.7057
    - kind: orchestrator
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 11
          cache_read: 67262
          cache_write: 10491
          cost: 0.0208
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-07T23:41:12Z
  value: 25
  by: planner-S-0634
  at: 2026-11-09T21:57:57Z
forecast:
  duration: 30m
  delivery: 2026-11-09T23:02:00Z
  basis: "flai forecast's 16m (104 s per unit of size over 27 done medium improvement stories, size 9) raised to 30m for three sequential tasks and a close-out whose integration tier alone took 6m26s in S-0595; 8th in the pull order with an in-progress limit of 5."
  by: planner-S-0634
  at: 2026-11-09T21:57:57Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:01:02Z
---
# S-0634 flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier

## Goal

This story remediates [I-0308](../../../design/issues/I-0308-flai-test-on-a-change-to-the-manifest-s-settings-catalog-does-not-run.md), "flai test on a change to the manifest's settings catalog does not run TestReferenceIsCurrent, so a stale flai-reference.md first fails in the close-out's integration tier". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0308 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0308 is closed with `flai issue close I-0308 --reason` saying what fixed it

## Tasks
- T-2773 Run TestReferenceIsCurrent and TestFlagIndexIsCurrent under -short
- T-2776 Have the settings-index tier run the reference tests and select docs/users/flai-reference.md
- T-2780 Describe the reference tests' tier in devex.md and tooling.md and close I-0308

## Notes

Cost of delay inputs set by flai from I-0308. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T22:56:50Z, 0 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from I-0308's instance and the code. `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` (`flai/cmd/reference_test.go`) skip under `-short` because they read the monorepo's docs. Together they take 0.04 s. The `go-test` tier already runs the tests of `flai/cmd`, an importer of `flai/internal/manifest`, under `-short` whenever `flai/internal/manifest/settings.go` or any command's flags change. So the fix is to skip them only when the docs are absent, as S-0592 did for `TestSettingsIndexHasEveryKey`. The `settings-index` tier then also runs them, selected by `docs/users/flai-reference.md` too, so a hand edit of either generated page is checked. No new tier is needed.

Touches, file by file:

| Touch | Source | Why |
|-------|--------|-----|
| `flai/cmd/reference_test.go` | design (I-0308) | The two tests that skip under `-short` |
| `system-flow.yaml` | design (I-0308) | The `settings-index` tier's `paths` and `-run` |
| `flai/internal/verify/select_test.go` | co-change | Pins the manifest's `settings-index` tier, whose argv changes |
| `design/system/devex.md` | design (D-0515's absorbing section) | § Used by this monorepo describes `settings-index` |
| `design/conventions/tooling.md` | co-change | Its project additions describe `settings-index` |
| `design/issues/I-0308-…md` | declared by criterion 2 | `flai issue close` writes it |
| `design/issues/summary.md` | layout | `flai issue close` regenerates it |

`flai touches suggest` also listed `flai/internal/verify/select.go`, `design/system/project-manifest.md`, and `docs/operators/settings.md`. They are left out: the fix changes no tier selection code and no manifest schema, and the settings page changes only if `make flai-reference` finds it stale. No folder touch. No `after`: the story reads nothing an open story writes.

Figures:

- Forecast 30m, raised from flai's 16m (size 9 from 2 criteria and 7 touches). The three tasks run one after another, each with its `flai test` run, and the close-out runs the integration tier, which took 6m26s in S-0595, and smoke, since the story touches `flai/` and `system-flow.yaml`.
- Cost of delay 25 USD a week, as `flai cod` gives it: 10m lost per 168h cycle at 150 USD an hour. It stands, since I-0308 has one instance.
- 2026-11-10T05:13:22Z: moved to in-progress: the rebase of story/S-0634 onto main stopped on conflicts in design/conventions/tooling.md, design/system/devex.md, which the acceptance of S-0634 by alex cannot settle, so it aborted the rebase; its agent settles it in the worktree and syncs it with flai stream sync S-0634
