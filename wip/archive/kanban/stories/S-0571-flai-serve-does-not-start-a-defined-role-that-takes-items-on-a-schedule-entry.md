---
id: S-0571
type: story
nature: remediation
title: flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it
status: done
owner: alex
created: 2026-11-02T21:51:45Z
updated: 2026-11-07T04:17:13Z
transitions:
  - to: ready
    at: 2026-11-07T03:48:43Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-07T03:48:53Z
    by: agent-S-0571
  - to: review
    at: 2026-11-07T04:14:58Z
    by: agent-S-0571
  - to: done
    at: 2026-11-07T04:17:13Z
    by: orchestrator
tags: [flai, serve, roles]
topics: [planning]
touches: [flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/settings.md, template/CHANGELOG.md, design/issues/I-0262-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1575
  priced_by: harness
  turns:
    - day: 2026-11-07
      hand_edits: 1
      work: 50
  models:
    - model: claude-opus-5-5
      input: 146
      output: 31445
      cache_read: 18094469
      cache_write: 336325
      cost: 6.8542
  strategic:
    - kind: orchestrator
      seconds: 80
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 16
          output: 115
          cache_read: 693115
          cache_write: 74406
          cost: 0.1864
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: flai
    at: 2026-11-02T21:51:45Z
  value: 25
  by: planner-S-0571
  at: 2026-11-07T03:45:53Z
forecast:
  duration: 1h
  delivery: 2026-11-07T10:15:00Z
  basis: "flai forecast's 24m (101 s per unit over 53 large-band remediation stories, size 14) raised to 1h: S-0554, which changed role_schedule.go and the same docs, took 77m, and this story repeats about two thirds of that scope plus a serve integration close-out"
  by: planner-S-0571
  at: 2026-11-07T03:45:53Z
finalized:
  by: orchestrator
  at: 2026-11-07T03:48:13Z
---
# S-0571 flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it

## Goal

This story remediates [I-0262](../../../design/issues/I-0262-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md), "flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0262 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0262 is closed with `flai issue close I-0262 --reason` saying what fixed it

## Tasks
- T-2511 The role scheduler starts a defined role that takes stories once per ready story under each-ready-story, with a held start kept per story
- T-2512 flai check reports a schedule entry under each-ready-story on a role whose item does not list story
- T-2515 The design, the roles README and its template copy, the user and operator guides, and the template changelog say a defined role runs once per ready story under each-ready-story
- T-2517 I-0262 is closed with a reason naming the per-story schedule starts and the check

## Notes

Cost of delay inputs set by flai from I-0262. time_lost_per_cycle 10m: 10m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T21:38:10Z, 0 days before this story; under one cycle counts as one).

### Planning

The planner's proposed fix, from I-0262's one instance, is the issue's first option. When a `{schedule: <setting>, when: each-ready-story}` entry comes round for a defined role that takes items, `flai serve` starts the role once for each ready story, in pull order (`roles.EachReadyStory`, as the planner's replanner does in `replan.go`). It keeps a held start per role, setting, and story, so that a start the budget state holds is tried again at each look. A start the checks refuse is logged for that story alone. This is what ADR-0250 §5 already says the condition means, so no new ADR is planned. The issue's other option, `flai check` reporting the entry on a defined role, would only hide the gap. A role whose `item` does not list `story` gets nothing from `each-ready-story`, so `flai check` reports that entry as `roles.file` too.

Touches, and where each came from (`flai touches suggest` listed nothing, since the story declared no touches):

| Touch | Source | Why |
|-------|--------|-----|
| `flai/internal/serve/role_schedule.go` | design, layout | I-0262 names `lookAt` here as where a role that takes items is skipped |
| `flai/internal/serve/role_schedule_test.go` | layout | The test that reproduces the cause, beside the scheduler's other tests |
| `flai/internal/roles/roles.go` | layout | The `roles.file` check that passes the entry; it gains the story-kind condition |
| `flai/internal/roles/roles_test.go` | layout | The check's test |
| `design/system/strategic-agents.md` | design | § On an event or a schedule (S-0418) says the entry starts nothing |
| `design/workflows/roles/README.md` | design | The trigger and condition tables, which say what `each-ready-story` does |
| `template/root/design/workflows/roles/README.md` | layout | The template's copy of that README, identical to it today |
| `docs/users/flai.md` | design | § What starts a role your project defines says a schedule starts a role that takes no item only |
| `docs/operators/settings.md` | design | The `roles.<name>.schedule` row says a role that takes items is not started by one |
| `template/CHANGELOG.md` | layout | Where flai's user-facing changes are recorded |
| `design/issues/I-0262-…md` and `design/issues/summary.md` | goal | Criterion 2 closes the issue, which rewrites the generated summary |

No folder touch is kept, and the story has no `after`: it reads no other open story's output. S-0554's launchers and `roles.EachReadyStory` are on main already.

Figures:

- Forecast 1h, raised from `flai forecast`'s 24m. S-0554 changed `role_schedule.go` and the same docs, and took 77m. This story repeats about two thirds of that scope. Its close-out also runs the integration tier, since it touches `flai/`.
- Cost of delay value 25.00 USD a week, as `flai cod` worked it out from the inputs flai set from I-0262: 10m lost per 168h cycle at 150 USD an hour. It stands, with no reason to adjust it.

### Accepted by the orchestrator

- Verified: 18f6f866dc80e3421ff3a427ce00cf2dadb5aa52
- At: 2026-11-07T04:17:13Z

Verdict: pass
- 1: flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/users/flai.md, docs/operators/settings.md, template/CHANGELOG.md
- 2: design/issues/I-0262-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md, design/issues/summary.md
