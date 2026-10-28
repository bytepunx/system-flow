---
id: S-0512
type: story
nature: remediation
title: search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name
status: backlog
owner: alex
created: 2026-10-26T16:03:03Z
updated: 2026-10-28T13:02:30Z
transitions: []
tags: [cli, template, dashboard]
topics: [cli, template, dashboard]
touches: [flai/internal/adr/adr.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/docedit/docedit.go, flai/internal/docedit/docedit_test.go, design/adrs/0000-template.md, template/root/design/adrs/0000-template.md, template/CHANGELOG.md, template/template.yaml, design/system/documentation-standard.md, design/system/flai-cli.md, docs/users/flai.md, flaiover/src/lib/server/search.test.ts, design/adrs/README.md, design/issues/I-0208-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md, design/issues/summary.md]
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
      seconds: 465
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 62
          output: 863
          cache_read: 6510372
          cache_write: 19372
          cost: 1.4336
draft: true
cost_of_delay:
  inputs:
    time_lost_per_cycle: 5m
    by: flai
    at: 2026-10-26T16:03:03Z
  value: 12.5
  by: planner-S-0512
  at: 2026-10-28T08:27:18Z
forecast:
  duration: 1h
  delivery: 2026-10-29T06:56:00Z
  basis: "Its own forecast of 1h; 27th in the pull order with an in-progress limit of 5, behind S-0502, S-0415, S-0414, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425, S-0506, S-0507, S-0508, S-0509, S-0510 and S-0511."
  by: flai
  at: 2026-10-28T13:02:30Z
---
# S-0512 search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name

## Goal

This story remediates [I-0208](../../../design/issues/I-0208-search-test-ts-lists-adrs-with-the-supersession-chain-fails-on-main-since-adr-0202-s-status-became-rejected-a-status-outside-the-four-the-test-and-the-adr-template-name.md), "search.test.ts 'lists ADRs with the supersession chain' fails on main since ADR-0202's status became rejected, a status outside the four the test and the ADR template name". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0208 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0208 is closed with `flai issue close I-0208 --reason` saying what fixed it

## Tasks
- T-2189 Record the ADR: rejected is the fifth ADR status, immutable once set, and flai check validates an ADR's status
- T-2190 flai keeps one list of the five ADR statuses and flai check reports adr.status on a status outside it
- T-2191 The ADR template's status comment and the documentation standard name rejected as the fifth status, with a template changelog entry and bump
- T-2192 flai doc refuses a change to a rejected ADR as to an accepted one, but its topics
- T-2193 search.test.ts reads the ADR statuses it allows from the template's status comment instead of its own list
- T-2194 flai-cli.md's check row and the user docs' ADR section name rejected and the adr.status finding
- T-2195 Close I-0208 with flai issue close saying what fixed it

## Notes

Cost of delay inputs set by flai from I-0208. time_lost_per_cycle 5m: 5m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-10-26T04:35:14Z, 0.5 days before this story; under one cycle counts as one).

### Planning

Planned by planner-S-0512 on 2026-10-28.

The proposed solution, from I-0208's one instance: S-0501 already added `rejected` to the test's list, so the test passes on main today; the cause, an ADR status landing on main that nothing validates, stays. The remediation makes `rejected` the fifth ADR status (a proposed decision the operator declined, recorded in favour of another, as ADR-0202 was for ADR-0203), immutable once set as an accepted one is, names it where the four are named, gives `flai check` a finding, `adr.status`, on a status outside the five, so that an unknown status stops a close-out before the vitest tier fails, and makes the test read the list from `design/adrs/0000-template.md`'s status comment instead of its own copy. The decision is recorded in an ADR first.

Touches, file by file:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/adr/adr.go` | layout | `Statuses` names what `flai adr new` takes; the list of every status an ADR may carry is kept beside it or in a package `check` and `docedit` both import, since `adr` imports both |
| `flai/internal/check/check.go`, `flai/internal/check/check_test.go` | layout | the `adr.*` findings and their tests; `adr.status` is added there |
| `flai/internal/docedit/docedit.go`, `flai/internal/docedit/docedit_test.go` | layout | `acceptedStatus` lists accepted, superseded, and deprecated; a rejected ADR becomes immutable the same way |
| `design/adrs/0000-template.md`, `template/root/design/adrs/0000-template.md` | issue | the status comment names the four; it gains `rejected` in both copies |
| `template/CHANGELOG.md`, `template/template.yaml` | co-change (21%), layout | a template change gets a changelog entry and a patch bump, as S-0498's did |
| `design/system/documentation-standard.md` | issue | the editing rule names accepted, superseded, and deprecated ADRs as immutable and says nothing of a status check |
| `design/system/flai-cli.md` | co-change (42%) | the `flai check` row lists the findings; `adr.status` is documented there |
| `docs/users/flai.md` | co-change (32%) | the ADR section names the statuses and the `adr.*` findings |
| `flaiover/src/lib/server/search.test.ts` | issue | the test's hard-coded list is replaced by the template's comment |
| `design/adrs/README.md` | co-change (26%) | the new ADR's index row |
| `design/issues/I-0208-….md`, `design/issues/summary.md` | criterion 2 | `flai issue close` writes both |

No folder touch is kept. The new ADR's file cannot be named until `flai adr new` allocates its number; `flai task done` widens the touches with it (ADR-0144).

`touches suggest` was run seeded with the three paths the issue names, since the story declared none; its other co-changed files (`flaiover-dashboard.md`, `docs/users/flaiover.md`, `conventions.md`, `search.ts`, `CLAUDE.md`) are left out: nothing in the dashboard's behaviour, the search index, or the conventions changes.

Forecast: `flai forecast` gave 5m before the touches were written, from size 2 (2 criteria, 0 touches) in the small band. Run again with 16 touches the size is 18, the large band; the figure written is that run's, with the reason recorded beside it below if adjusted.

Cost of delay: 12.50 USD a week stands, from the operator's input of 5m lost per 168h cycle at 150 USD an hour; the immediate failure is fixed on main by S-0501, and the input is the issue's, not the planner's to change.
