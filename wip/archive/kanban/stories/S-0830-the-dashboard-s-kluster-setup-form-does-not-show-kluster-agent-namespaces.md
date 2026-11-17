---
id: S-0830
type: story
nature: remediation
title: The dashboard's kluster setup form does not show kluster.agent_namespaces
status: done
owner: alex
created: 2026-11-17T08:17:22Z
updated: 2026-11-17T10:14:13Z
transitions:
  - to: ready
    at: 2026-11-17T10:04:07Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T10:04:54Z
    by: agent-S-0830
  - to: review
    at: 2026-11-17T10:12:48Z
    by: agent-S-0830
  - to: done
    at: 2026-11-17T10:14:13Z
    by: orchestrator
tags: [flaiover, kluster]
topics: [kluster]
touches: [flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/components/DeploySettings.svelte.test.ts, flaiover/src/lib/server/repo.ts, design/components/flaiover/src/lib/components/DeploySettings.examples.json, design/system/flaiover-dashboard.md, design/system/kluster.md, docs/operators/settings.md, docs/users/flaiover.md, design/issues/I-0356-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 497
  priced_by: harness
  turns:
    - day: 2026-11-17
      test_runs: 1
      work: 37
  models:
    - model: claude-opus-5-5
      input: 84
      output: 14984
      cache_read: 8047025
      cache_write: 247468
      cost: 3.8892
  strategic:
    - kind: planner
      seconds: 16
      estimated: true
      models:
        - model: claude-haiku-5-5
          input: 8
          output: 24
          cache_read: 645501
          cache_write: 13179
          cost: 0.045
        - model: claude-opus-5-5
          input: 37
          output: 3075
          cache_read: 2593309
          cache_write: 45690
          cost: 0.9539
    - kind: orchestrator
      seconds: 69
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 17
          output: 167
          cache_read: 644027
          cache_write: 129428
          cost: 0.2409
        - model: claude-sonnet-5-5
          input: 10
          output: 349
          cache_read: 95098
          cache_write: 51298
          cost: 0.1179
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-17T08:17:22Z
  value: 25
  by: planner-S-0830
  at: 2026-11-17T08:28:38Z
forecast:
  duration: 25m
  delivery: 2026-11-17T10:39:00Z
  basis: "Its own forecast of 25m; 1st in the pull order with an in-progress limit of 3, behind S-0754 and S-0806."
  by: flai
  at: 2026-11-17T10:04:46Z
finalized:
  by: orchestrator
  at: 2026-11-17T08:33:52Z
---
# S-0830 The dashboard's kluster setup form does not show kluster.agent_namespaces

## Goal

This story remediates [I-0356](../../../design/issues/I-0356-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md), "The dashboard's kluster setup form does not show kluster.agent_namespaces". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0356 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0356 is closed with `flai issue close I-0356 --reason` saying what fixed it

## Tasks
- T-3526 The kluster setup form shows and saves kluster.agent_namespaces
- T-3528 The design and the guides say the setup form sets kluster.agent_namespaces, and I-0356 is closed

## Notes

Cost of delay inputs set by flai from I-0356. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-13T08:42:09Z, 4 days before this story; under one cycle counts as one).

### Planning

Proposed fix, from the instance. flai already gives the setting to the dashboard: `kluster.status`'s `settings` lists `kluster.agent_namespaces` with kind `number`, whole, default `2`, and a JSON number value (`flai/internal/manifest/settings.go`, `klusterSettings` in `flai/internal/hostapi/kluster.go`). `flai manifest set` and `settings.manifest` already write and check it. So the fix is in flaiover alone, in three places. `PLAIN_KEYS` in `DeploySettings.svelte` gains the key. `typed()` shows a number value as its digits. `KlusterSetting` in `repo.ts` gains the `number` kind. No Go change is planned.

Touches:

| Path | From |
|------|------|
| `flaiover/src/lib/components/DeploySettings.svelte` | the issue's instance names it; layout |
| `flaiover/src/lib/components/DeploySettings.svelte.test.ts` | co-change (`flai touches suggest --min 1`) |
| `design/components/flaiover/src/lib/components/DeploySettings.examples.json` | co-change; `components.examples` keeps it in step with the props |
| `flaiover/src/lib/server/repo.ts` | layout: `KlusterSetting` has no `number` kind |
| `design/system/kluster.md` | design: says the setup form does not show it yet |
| `docs/operators/settings.md` | design: its row says the same |
| `design/system/flaiover-dashboard.md` | design: the `/ops/klusters` row lists the form's fields |
| `docs/users/flaiover.md` | design: step 4 lists the form's fields |

Co-change also listed `flaiover/src/lib/settings.ts` and its test, at 1 of 4158 commits. They are left out: the form's existing helpers cover a number. `flai task done` widens the touches if the agent needs them. The story has no folder touch and no `after`.

Figures:

- Forecast 25m, delivery 2026-11-17T15:09:00Z, as `flai forecast` gives them: the median of 68 done remediation stories in the medium band, times size 10 (2 criteria, 8 touches). It stands. Half the touches are one-sentence doc edits, and the basis already holds flaiover close-outs.
- Cost of delay 25.00 USD a week, as `flai cod` works it from the inputs flai set from I-0356: 10m lost per 168h cycle at 150 USD an hour. It stands: one occurrence, and `flai manifest set` is a workaround.

### Accepted by the orchestrator

- Verified: 30d18f601884be6438a2f42c48c63546bc1bf192
- At: 2026-11-17T10:14:13Z

Verdict: pass
- 1: flaiover/src/lib/components/DeploySettings.svelte, flaiover/src/lib/server/repo.ts, flaiover/src/lib/components/DeploySettings.svelte.test.ts, design/components/flaiover/src/lib/components/DeploySettings.examples.json
- 2: design/issues/I-0356-the-dashboard-s-kluster-setup-form-does-not-show-kluster-agent-namespaces.md, design/issues/summary.md
