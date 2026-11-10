---
id: S-0625
type: story
nature: improvement
title: "Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand"
status: in-progress
owner: alex
created: 2026-11-07T23:40:15Z
updated: 2026-11-10T05:19:09Z
transitions:
  - to: ready
    at: 2026-11-09T22:04:44Z
    by: alex
  - to: in-progress
    at: 2026-11-10T05:12:17Z
    by: agent-S-0625
tags: [flai, template]
topics: [release, orchestration, conventions]
touches: [design/adrs/README.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, flai/cmd/release.go, flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, flai/cmd/release_evaluate_test.go, flai/cmd/publish_judgement_test.go, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, template/CHANGELOG.md, design/issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md, design/issues/summary.md, design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1905
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-11-10
      ceremony: 3
      hand_edits: 1
      work: 20
  models:
    - model: claude-opus-5-5
      input: 144
      output: 986
      cache_read: 16546411
      cache_write: 1197818
      cost: 6.5209
  strategic:
    - kind: orchestrator
      seconds: 69
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 11
          output: 104
          cache_read: 613670
          cache_write: 15815
          cost: 0.1674
    - kind: analyzer
      seconds: 42
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 15
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
cost_of_delay:
  inputs:
    time_lost_per_cycle: 10m
    by: orchestrator
    at: 2026-11-09T21:57:26Z
  value: 25
  by: planner-S-0625
  at: 2026-11-09T22:01:20Z
forecast:
  duration: 50m
  delivery: 2026-11-10T00:32:00Z
  basis: "Its own forecast of 50m; 2nd in the pull order with an in-progress limit of 5, behind S-0633, S-0636, S-0638 and S-0640."
  by: flai
  at: 2026-11-09T23:22:43Z
finalized:
  by: alex
  at: 2026-11-09T21:56:42Z
---
# S-0625 Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand

## Goal

This story remediates [I-0291](../../../design/issues/I-0291-under-the-judgement-release-policy-the-publish-permission-publishes-nothing.md), "Under the judgement release policy the publish permission publishes nothing: every release since 2026-11-02 was the operator's by hand", which the analysis in [2026-11-07-all.md](../../../design/analysis/2026-11-07-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0291 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0291 is closed with `flai issue close I-0291 --reason` saying what fixed it

## Tasks
- T-2788 Record in an ADR how the orchestrator's judgement reaches the publish formula's evaluate step
- T-2790 flai release --evaluate --judged meets the judgement policy with the orchestrator's reason
- T-2791 The orchestrator passes its judgement to the publish formula as the judgement input
- T-2792 The publish formula takes a judgement input and its evaluate step passes it as --judged
- T-2793 Document publishing under judgement through the formula and close I-0291

## Notes

The analyzer found I-0291: the analysis the goal links gives its evidence.

Strategic usage carried over from I-0291: analyzer 42s, 0.8156 USD.

### Planning

Proposed solution. Under `judgement`, `flai release --evaluate` never sets `Met` (`flai/internal/release/evaluate.go`, the judgement case of `Evaluate`). The publish formula's step 2 is `release --evaluate --require-met`, so it exits 2. Since S-0621 that exit stops the instance cleanly, releasing nothing. `release_publish` refuses while the formula exists (`noPublishFormula` in `flai/internal/mcpserver/orchestrate.go`). So the orchestrator's judgement has no way to reach the step that checks it. The fix carries it there:

- `flai release --evaluate` gains `--judged "<reason>"`. Under `judgement`, a non-empty reason makes the evaluation met when something is accepted and not yet released and `whole_epics` holds nothing back. Its reason names the judgement. Under `threshold` and `theme` the flag changes nothing, so one formula serves every policy.
- The publish formula declares an input `judgement`, with `default: ""`, and step 2 runs `release --evaluate --require-met --judged "{{judgement}}"`. With no input, step 2 behaves as today.
- Under `judgement`, the orchestrator runs `flai formula run publish --ready --autocommit --input judgement="<its reasoning>"`. The guard already passes `--input` on a formula entry (`flai/internal/guard/roles.go`, `formulaRunFlags`), so the guard does not change.

This project's own symptom ended on 2026-11-07, when the operator set `orchestration.release.policy: threshold` (TH-0811, D-0504). The cause remains for any project on `judgement`, which is the default.

Touches, and where each came from:

| Touch | Source | Why |
|-------|--------|-----|
| `design/adrs/0272-…md`, `design/adrs/README.md` | design | A new ADR refines ADR-0240 and ADR-0094. The ADR's number and slug are predicted. |
| `design/system/strategic-agents.md` | design, co-change | § Publishing a release (S-0222) and its `judgement` row |
| `design/system/flai-cli.md` | co-change (22%) | The `release --evaluate [--require-met]` row |
| `design/system/project-manifest.md` | design | Says judgement "is never met by itself" |
| `flai/cmd/release.go`, `flai/internal/release/evaluate.go` | layout | The flag and the judgement case |
| `flai/internal/release/evaluate_test.go`, `flai/cmd/release_evaluate_test.go` | layout | Judgement tests, `TestReleaseEvaluateRequireMet` |
| `flai/cmd/publish_judgement_test.go` | layout | New: runs the real publish formula's evaluate step under judgement, the test that reproduces I-0291 |
| `design/workflows/formulas/publish.md`, `template/root/…/publish.md` | declared by the issue, layout | The input and step 2 |
| `design/workflows/roles/orchestrator.md`, its template copy, `flai/internal/roles/baseline/orchestrator.md` | layout | The publish paragraph |
| `flai/internal/harness/harness_test.go` | layout | `TestTheOrchestratorsPromptSaysHowItPublishes` |
| `design/conventions/strategic-agents.md`, its template copy | design | § As the orchestrator, the publish rule |
| `docs/users/flai.md`, `docs/users/flai-reference.md`, `docs/operators/settings.md` | co-change (30%, 37%, 22%) | § Whether a release is due and § When the orchestrator publishes; the generated flag tables |
| `template/CHANGELOG.md` | layout | The formula, role, and convention changes |
| `design/issues/I-0291-…md`, `design/issues/summary.md` | design | Closed with `flai issue close` |

`flai touches suggest` also listed `flai/internal/release/release.go`, `flai/cmd/accept.go`, `flai/cmd/push.go`, `flai/internal/preview/push.go`, `design/conventions/git.md`, `design/system/workflow.md`, and `docs/operators/index.md`. They are left out because publishing itself, acceptance, and pushing do not change. `flai/internal/mcpserver/orchestrate.go` is left out too, since `release_publish` already takes a judgement reason. There are no folder touches.

No `after`. The template's publish formula gains a `{{judgement}}` placeholder. S-0616, in progress, changes how `TestRenderPrototypeTemplate` treats `{{`. Today's test already passes declared formula inputs (`withoutFormulaInputs`), and S-0616's change passes copied files. Either way the placeholder passes, so this story reads nothing S-0616 writes.

Figures:

- Forecast: `flai forecast` gave 3m: 74 s per size unit at size 2, with no touches declared then. I raised it to 50m. S-0621 is the same kind of change (a flai flag or command, the publish formula in both copies, the orchestrator's instructions, docs, close an issue), and its agent took 49m (2954 s). Delivery is 22:50Z, after the stories ahead in the pull order.
- Cost of delay: 25 USD a week, as `flai cod` worked it out from `time_lost_per_cycle: 10m` at 150 USD an hour over a 168h cycle. The orchestrator set that input on TH-0943 as recommended. It stands: this project no longer loses time to the cause, and what remains is a project on `judgement`.
