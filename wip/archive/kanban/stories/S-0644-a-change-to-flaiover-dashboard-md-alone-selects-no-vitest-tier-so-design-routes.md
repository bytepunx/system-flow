---
id: S-0644
type: story
nature: improvement
title: A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route
status: done
owner: alex
created: 2026-11-09T21:23:04Z
updated: 2026-11-09T23:22:42Z
transitions:
  - to: ready
    at: 2026-11-09T22:19:59Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T22:48:21Z
    by: agent-S-0644
  - to: review
    at: 2026-11-09T23:21:06Z
    by: agent-S-0644
  - to: done
    at: 2026-11-09T23:22:42Z
    by: orchestrator
tags: [flai, flaiover, testing]
topics: [testing, flaiover, docs]
touches: [system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md, design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md, design/issues/summary.md, design/issues/I-0323-flai-check-finds-item-duplicate-id-outside-the-story-at-close-out.md, design/issues/I-0327-testrepositorymanifestlintschangedflaioverfiles-skips-under-short-so-a-new-tier.md, design/issues/I-0324-testeachrequestistimedbeneaththetransport-reads-the-request-answered-event.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1978
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 4
      hand_edits: 2
      work: 97
  models:
    - model: claude-opus-5-5
      input: 218
      output: 21465
      cache_read: 14474940
      cache_write: 158875
      cost: 4.5962
  strategic:
    - kind: orchestrator
      seconds: 67
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 8
          output: 78
          cache_read: 664480
          cache_write: 7688
          cost: 0.179
        - model: claude-sonnet-5-5
          input: 2
          output: 11
          cache_read: 14233
          cache_write: 8148
          cost: 0.0166
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-09T21:23:04Z
  value: 12.5
  by: planner-S-0644
  at: 2026-11-09T22:08:10Z
forecast:
  duration: 25m
  delivery: 2026-11-09T23:10:00Z
  basis: "Its own forecast of 25m; 2nd in the pull order with an in-progress limit of 5, behind S-0622, S-0623, S-0632, S-0633, S-0634 and S-0645."
  by: flai
  at: 2026-11-09T22:20:09Z
finalized:
  by: orchestrator
  at: 2026-11-09T22:19:32Z
---
# S-0644 A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route

## Goal

This story remediates [I-0322](../../../design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md), "A change to flaiover-dashboard.md alone selects no vitest tier, so design-routes.test.ts does not run when the dashboard design gains a row with no route". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0322 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0322 is closed with `flai issue close I-0322 --reason` saying what fixed it

## Tasks
- T-2809 Add a design-routes tier that runs design-routes.test.ts when the dashboard design or a route changes
- T-2810 Describe the design-routes tier in devex.md and tooling.md and close I-0322

## Notes

Cost of delay inputs set by flai from I-0322. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-09T17:35:44Z, 0.2 days before this story; under one cycle counts as one).

### Planning

Proposed fix. Add a narrow tier, `design-routes`, to `system-flow.yaml`'s `tests`, as S-0592 added `settings-index` and `flaiover-flai` for I-0304. It runs `flaiover/src/lib/server/design-routes.test.ts` alone through `scripts/flaiover-unit.sh`, in `flaiover`, with format `vitest-json`. Its paths are `design/system/flaiover-dashboard.md` and `flaiover/src/routes/**`. The vitest tier cannot carry the fix: it runs `vitest related` on the changed files, which follows the TypeScript import graph. The test reads the design and the routes folder from disk and imports neither, so neither a change to the design nor a deleted route selects it. The `flaiover` tier's `covers` gains `design-routes`, so under `--all` the test runs once. A repository-manifest test in `flai/internal/verify/select_test.go` reproduces the gap, as `TestRepositoryManifestChecksTheSettingsIndexWhenItChanges` does for I-0304.

Touches, and where each came from:

- `system-flow.yaml`: design. The issue names the manifest's `tests`. S-0592 changed that block by hand in 38e021d62.
- `flai/internal/verify/select_test.go`: layout. The repository-manifest tier tests live here (I-0304, S-0592).
- `design/system/devex.md`: design. Its Test tiers row lists every tier and why each narrow one exists.
- `design/conventions/tooling.md`: layout. Its project additions list the tiers `flai test` runs.
- `design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md` and `design/issues/summary.md`: declared by criterion 2, since `flai issue close` writes both.

`flai touches suggest` listed only weak co-changes, each in 6% of commits or fewer, such as `docs/users/flai.md` and `design/system/flai-cli.md`. None applies: the fix adds no command and changes no flag. `design/system/flaiover-dashboard.md` is what the tier selects, not what the story changes. No folder touch.

No `after`. S-0634, in progress, changes the same manifest block, `select_test.go`, `devex.md`, and `tooling.md` for I-0308. This story reads nothing S-0634 writes, so a trial merge judges the overlap (ADR-0230).

Forecast: `flai forecast` gave 14m, from 104 s per unit of size over 27 done medium improvement stories, size 8. Raised to 25m. The story has two tasks in sequence. The worktree has to install flaiover's dependencies before vitest runs. The close-out runs integration and smoke, because `select_test.go` is under `flai/`, and integration alone took 6m26s in S-0595. Delivery is flai's 00:16Z moved by the 11m added.

Cost of delay: `flai cod` gave 12.50 USD a week, from the 5m of time lost per 168h cycle at 150 USD an hour. Kept: the input is the operator's, from I-0322's one instance, and no reason to change it.

### Accepted by the orchestrator

- Verified: 4da34bdfda6e1d8c88e0b1e47d453ebb3ed4cd10
- At: 2026-11-09T23:22:42Z

Verdict: pass. Both criteria are met at 4da34bdfda6e; the close-out passed every step there, and every changed file is within the story's touches.
- 1: system-flow.yaml, flai/internal/verify/select_test.go, design/system/devex.md, design/conventions/tooling.md
- 2: design/issues/I-0322-a-change-to-flaiover-dashboard-md-alone-selects-no-vitest-tier-so-design-routes.md, design/issues/summary.md
