---
id: S-0617
type: story
nature: remediation
title: "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai"
status: ready
owner: alex
created: 2026-11-07T23:39:25Z
updated: 2026-11-09T10:38:30Z
transitions:
  - to: ready
    at: 2026-11-09T10:38:24Z
    by: orchestrator
tags: [flai]
touches: [flai/internal/release/outdated.go, flai/internal/release/outdated_test.go, flai/internal/serve/outdated_test.go, design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md, design/issues/summary.md]
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
    - kind: planner
      seconds: 205
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 46
          output: 13937
          cache_read: 2923178
          cache_write: 144249
          cost: 2.0176
    - kind: orchestrator
      seconds: 70
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 20
          output: 108
          cache_read: 736488
          cache_write: 114542
          cost: 0.2229
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-11-07T23:39:25Z
  value: 12.5
  by: planner-S-0617
  at: 2026-11-09T10:35:14Z
forecast:
  duration: 30m
  delivery: 2026-11-09T11:14:00Z
  basis: "Its own forecast of 30m; 1st in the pull order with an in-progress limit of 5, behind S-0610, S-0615 and S-0616."
  by: flai
  at: 2026-11-09T10:38:30Z
finalized:
  by: orchestrator
  at: 2026-11-09T10:37:07Z
---
# S-0617 TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai

## Goal

This story remediates [I-0252](../../../design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md), "TestServeSaysTheFlaiIsOlderThanTheProject fails under load: project.info is asked before flai serve has worked out the outdated flai". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0252 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0252 is closed with `flai issue close I-0252 --reason` saying what fixed it

## Tasks
- T-2723 FlaiOutdated keeps no answer from a git that failed, so one failed look is not served as not outdated for a minute
- T-2724 Show the outdated serve tests hold under load since S-0569, and fix what still fails
- T-2725 Close I-0252 naming S-0569's fix and T-2723's

## Notes

Cost of delay inputs set by flai from I-0252. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T06:04:39Z, 5.7 days before this story; under one cycle counts as one).

### Planning

The proposed solution, from the instance. I-0252's one instance, S-0506 on 2026-11-02, failed at the old `flai/internal/serve/outdated_test.go:175`, `t.Errorf("project.info: %v %s", rerr, res)`, printing `project.info: <nil>`: no error and no result. That is I-0247's symptom, a watcher notification read as the reply, which S-0569 fixed on 2026-11-06 (`ea5280843`), after the instance. `project.info` works out `flai_outdated` itself at each call (`flai/internal/hostapi/hostapi.go`), so it does not wait on serve's own look. The issue's reading does not hold. One way is left for it to answer `flai_outdated` nil for the right project: `release.FlaiOutdated` keeps a failed git look's nil for a minute. T-2723 removes that, with a test that reproduces it. T-2724 shows the serve tests hold under load. T-2725 closes the issue.

Touches, none declared before this plan:

- `flai/internal/release/outdated.go`, `flai/internal/release/outdated_test.go`: design and layout. The answer cache T-2723 changes, and its tests.
- `flai/internal/serve/outdated_test.go`: design. The failing test, which T-2724 changes only if its stress run fails.
- `design/issues/I-0252-testservesaystheflaiisolderthantheproject-fails-under-load-project-info-is.md`, `design/issues/summary.md`: layout. `flai issue close` writes both (T-2725).
- Co-change: `flai touches suggest` lists only issue files and docs that change with the summary. None bears on this story. No user-facing behaviour changes, so no docs.

No folder touch. No `after`: the story reads no other open story's output. S-0569 is done.

Figures:

- Forecast: flai gives 14m, 119 s per unit of size over 43 medium remediation stories, size 7. Raised to 30m for T-2724's stress run under load, which past remediations did not run. Delivery is flai's 10:55Z moved by the 16m added, 11:11Z.
- Cost of delay: 12.50 USD a week, as `flai cod` works it out from the input flai set from I-0252: 5m lost per 168h cycle at 150 USD an hour. Kept: the issue has one instance, and nothing gives a reason to change it.
