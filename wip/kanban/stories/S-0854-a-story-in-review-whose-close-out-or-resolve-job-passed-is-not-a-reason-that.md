---
id: S-0854
type: story
nature: improvement
title: A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event
status: ready
owner: alex
created: 2026-11-27T22:55:26Z
updated: 2026-11-28T00:03:04Z
transitions:
  - to: ready
    at: 2026-11-28T00:02:05Z
    by: claude-concierge
tags: [flai, template, serve]
topics: [orchestration, review]
touches: [design/adrs/README.md, design/adrs/0322-a-close-out-or-resolve-job-that-passed-on-a-story-in-review-starts-the.md, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go, docs/users/flai.md, template/CHANGELOG.md, design/issues/I-0413-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md, design/issues/summary.md]
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
      seconds: 40
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 25
          output: 5909
          cache_read: 1476986
          cache_write: 57063
          cost: 0.7659
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29514
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
cost_of_delay:
  inputs:
    time_lost_per_cycle: 48h
    by: flai
    at: 2026-11-27T22:55:26Z
  value: 7200
  by: planner-S-0854
  at: 2026-11-27T23:59:03Z
forecast:
  duration: 2h
  delivery: 2026-11-28T06:30:00Z
  basis: "Its own forecast of 2h; 8th in the pull order with an in-progress limit of 3, behind S-0837, S-0844, S-0846, S-0835, S-0814, S-0852, S-0813, S-0822, S-0848 and S-0853."
  by: flai
  at: 2026-11-28T00:03:04Z
finalized:
  by: alex
  at: 2026-11-28T00:01:54Z
---
# S-0854 A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event

## Goal

This story remediates [I-0413](../../../design/issues/I-0413-a-story-in-review-whose-close-out-or-resolve-job-passed-is-not-a-reason-that.md), "A story in review whose close-out or resolve job passed is not a reason that starts the orchestrator, so an eligible story waits for the operator or for an unrelated event", which the analysis in [2026-11-27-all.md](../../../design/analysis/2026-11-27-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0413 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0413 is closed with `flai issue close I-0413 --reason` saying what fixed it

## Tasks
- T-3641 Record in an ADR that a passed close-out or resolve job of a story in review starts the orchestrator, and update the design and conventions
- T-3642 Add the verified event and the passed-in-review condition to the role trigger vocabulary
- T-3645 flai serve starts the orchestrator with the reason review:S-nnnn on a passed verified event of a story in review
- T-3646 Add the passed-in-review verified trigger to the orchestrator's role file, its baseline, and its template copy
- T-3648 Document the verified trigger for users and the template's changelog, and close I-0413 with what fixed it

## Notes

The analyzer found I-0413: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0413. time_lost_per_cycle 48h carried over from I-0413's Impact section. Its Impact time_lost_per_cycle was taken rather than the 48h derived from its cost and count.

Strategic usage carried over from I-0413: analyzer 1m25s, 1.4522 USD.

### Planning

Proposed solution, from I-0413's instance. The close-out job and the rebase-and-resolve job already append a `verified` event with verdict `passed` when they end (`announce` in `flai/internal/serve/verifyjob.go`). The orchestrator has no story of its own, so its inbox gets the event. Nothing matches it, for three reasons. `verified` is not in the trigger vocabulary (`Events` in `flai/internal/roles/roles.go`). `eventNames` in `flai/internal/serve/orchestrate_reasons.go` has no case for it. Its `prime` maps no verified reason to a form. The fix makes the event a trigger and gives it the review reason:

- `verified` joins the trigger events. A new condition, `passed-in-review`, holds when three things are true: the verdict is `passed`, the job ran something (its `recorded` field is not `true`), and the story is in review when flai serve looks.
- The orchestrator's role file, in its three copies, gains `{event: verified, when: passed-in-review}`. So do the built-in triggers flai serve uses when the file is missing.
- A matched event becomes the reason `review:S-nnnn`, primed as a review's story and its last verify result. The orchestrator's run that left the story "for a later run" then gets that run.
- A verdict answered from the record is left out. Otherwise the orchestrator's own `verify` call on a story it cannot accept would start it again, in a loop.

Touches, and where each came from:

| Touch | Source | Why |
| --- | --- | --- |
| `flai/internal/roles/roles.go`, `roles_test.go`, `conditions.go`, `conditions_test.go` | layout | The trigger events, the conditions table, and `Occurrence`, which gains the verdict, the recorded flag, and the story's state. |
| `flai/internal/serve/orchestrate_reasons.go`, `orchestrate_reasons_test.go`, `orchestrate_test.go` | layout | `builtinTriggers`, `eventNames`, `occurrenceOf`, and `prime`, and the test that reproduces I-0413. |
| `template/root/design/workflows/roles/orchestrator.md`, `flai/internal/roles/baseline/orchestrator.md`, `design/workflows/roles/orchestrator.md` | co-change | `flai touches suggest` puts them at 52% and 37%, and the project's own copy is the seed. All three hold the same trigger list. |
| `flai/internal/harness/harness_test.go` | co-change | 22%. It reads the generated definitions from the baseline role files. |
| `design/workflows/roles/README.md` and its template copy | design | They list the trigger events and conditions. |
| `design/system/strategic-agents.md` | design, co-change | § The orchestrator › Starting and stopping it holds the trigger table. § Primed for its reasons (S-0537) maps reasons to forms. |
| `design/conventions/strategic-agents.md` and its template copy | design, co-change | § Accepting a story says "leave the story for a later run". It now names what starts that run. |
| `docs/users/flai.md` | design | Its role trigger table and its `review:S-nnnn` row. |
| `template/CHANGELOG.md` | co-change | The template's role file and convention change. |
| `design/adrs/README.md`, `design/adrs/0322-…` | design | One ADR, a change to the trigger vocabulary as ADR-0211 was. Its number is a guess: `flai adr new` takes one more than the highest on any branch. `flai task done` corrects the touch, so no folder touch is kept. |
| `design/issues/I-0413-…`, `design/issues/summary.md` | criteria | Criterion 2 closes I-0413 with `flai issue close`, which writes both. |

Some co-change files are left out because no tool, permission, setting, or other role changes: `flai/internal/harness/harness.go`, `flai/internal/serve/orchestrate.go`, `flai/internal/guard/*`, `flai/internal/mcpserver/server.go`, `design/workflows/roles/planner.md`, and `docs/operators/settings.md`.

After: none. No open story writes a file, command, schema, or decision that this one reads. S-0845 and S-0824 also change the resolve job and the orchestrator's role file. A trial merge judges that.

Figures:

- **Forecast.** With the touches written, `flai forecast` gives 44m: 109 s per unit of size over 187 done large improvement stories on claude-opus-5-5, times size 24, delivered 2026-11-28T07:12Z. It is raised to 2h. The work is five tasks across two Go packages, an ADR, three role-file copies, and the docs, and a close-out takes about 10m. S-0847, a serve change with an ADR and three tasks, took 1h03m in progress. The delivery moves later by the same 1h16m, to 08:28Z.
- **Cost of delay.** `flai cod` gives 7,200 USD a week: 48h lost per 168h cycle at 150 USD an hour. It comes from the operator's input `time_lost_per_cycle` 48h, carried over from I-0413. It stands as given. In I-0413's instance, two stories with passing verdicts waited six days in review.
