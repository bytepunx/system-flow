---
id: S-0696
type: story
nature: remediation
title: "The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say"
status: done
owner: alex
created: 2026-11-10T08:54:31Z
updated: 2026-11-11T08:48:36Z
transitions:
  - to: ready
    at: 2026-11-11T02:30:59Z
    by: alex
  - to: in-progress
    at: 2026-11-11T08:00:55Z
    by: agent-S-0696
  - to: review
    at: 2026-11-11T08:21:16Z
    by: agent-S-0696
  - to: done
    at: 2026-11-11T08:48:36Z
    by: alex
tags: [flai, testing, manifest]
topics: [testing]
touches: [system-flow.yaml, flai/internal/verify/select_test.go, design/system/project-manifest.md, design/system/devex.md, design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1257
  priced_by: harness
  turns:
    - day: 2026-11-11
      test_runs: 1
      work: 38
  models:
    - model: claude-opus-5-5
      input: 84
      output: 12060
      cache_read: 9565281
      cache_write: 280075
      cost: 4.3952
  strategic:
    - kind: orchestrator
      seconds: 49
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 121
          cache_read: 677068
          cache_write: 108427
          cost: 0.2183
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-10T08:54:31Z
  value: 12.5
  by: planner-S-0696
  at: 2026-11-10T10:29:19Z
forecast:
  duration: 30m
  delivery: 2026-11-11T08:00:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 3, behind S-0686, S-0691 and S-0693."
  by: flai
  at: 2026-11-11T07:23:13Z
finalized:
  by: orchestrator
  at: 2026-11-10T11:04:45Z
---
# S-0696 The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say

## Goal

This story remediates [I-0270](../../../design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md), "The integration and smoke tiers' when names template/**, but their paths do not select a template-only change, so a story that changes only template/ runs neither at close-out, against what code-quality.md and tooling.md say". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0270 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0270 is closed with `flai issue close I-0270 --reason` saying what fixed it

## Tasks
- T-2987 Add template/** to the integration and smoke tiers' paths, with a manifest test that a template-only change runs both
- T-2990 Say in project-manifest.md and devex.md that this repository's integration and smoke paths select template/**
- T-2991 Close I-0270 naming the widened tier paths and the test that holds them

## Notes

Cost of delay inputs set by flai from I-0270. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-07T04:41:00Z, 3.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix: add `template/**` to the `paths` of the `integration` and `smoke` tiers in `system-flow.yaml`. `verify.SelectStory` (`flai/internal/verify/select.go`) picks an `all_only` tier only when its `paths` select a changed path, and `SelectStoryLeavingQueue` applies `when` only to the tiers already picked. So `when: template/**` does nothing today. The merge queue's batch also picks tiers by `paths`, so a template-only change runs neither tier anywhere before main. With the paths widened, `code-quality.md`, `tooling.md`, and `docs/users/flai.md` are true as written. Many Go tests read `template/`, for example `flai/internal/roles/baseline_test.go` and `flai/internal/check/roles_test.go`. The other fix, dropping `template/**` from `when` and from those documents, would leave those tests to run only on the next change under `flai/`. The planner rejected it for that reason.

Touches, and where each came from:

- `system-flow.yaml`: design. The issue names the two tiers' `paths` and `when`. S-0634 and S-0644 changed tiers there by hand in the same way.
- `flai/internal/verify/select_test.go`: layout. It holds this repository's manifest tests (`TestRepositoryManifest...`), where the regression test belongs.
- `design/system/project-manifest.md`: design. Its `tests` section says "this repository gives `integration` and `smoke` `flai/**`".
- `design/system/devex.md`: design. Its test-tiers row says "integration, and smoke (both `flai/**`)".
- `design/issues/I-0270-the-integration-and-smoke-tiers-when-names-template-but-their-paths-do-not.md` and `design/issues/summary.md`: declared by the second criterion. `flai issue close` writes both.
- Co-change (`flai touches suggest`) also listed `design/system/flai-cli.md` (43%) and `docs/users/flai.md` (41%). The planner left both out. Read in full, `docs/users/flai.md` already says a story that changes `template/` runs both tiers, and the planner found no line in `flai-cli.md` naming this repository's tier paths. The other co-change files cover settings and the template, which this story does not change.

Figures:

- Forecast 30m, delivery 2026-11-10T14:56:00Z. `flai forecast` gave 18m: 134 s per unit over 54 done remediation stories in the medium band, at size 8 (2 criteria, 6 touches). The planner added about 12m, the time `integration` and `smoke` take together (ADR-0233). This story's close-out is sure to run them, since it changes `flai/internal/verify`. The delivery is flai's 14:44 plus those 12 minutes.
- Cost of delay value 12.50 USD a week, as `flai cod` worked it out from the operator's input: 5m lost per 168h cycle at 150 USD an hour. Kept as given. The risk the issue names is a template change that breaks a repository test and reaches main untested. That risk is larger than one 5m occurrence, but the inputs are the operator's.
