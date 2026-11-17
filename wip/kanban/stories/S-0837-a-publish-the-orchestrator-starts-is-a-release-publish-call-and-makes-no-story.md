---
id: S-0837
type: story
nature: improvement
title: A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired
status: backlog
owner: alex
created: 2026-11-17T08:21:17Z
updated: 2026-11-17T12:57:33Z
transitions: []
tags: [flai, template]
topics: [release, orchestration, upgrade]
touches: [design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md, design/adrs/0240-under-publish-the-orchestrator-starts-the-baseline-publish-formula-whose-main.md, design/adrs/0273-under-judgement-the-orchestrator-s-judgement-reaches-the-publish-formula-as-an.md, design/adrs/README.md, flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_test.go, flai/cmd/formula.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/release.go, flai/cmd/release_push_action_test.go, docs/users/flai-reference.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/harness/harness_test.go, flai/internal/roles/roles_test.go, flai/internal/guard/roles_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, flai/cmd/upgrade.go, flai/cmd/upgrade_test.go, template/root/design/workflows/formulas/publish.md, template/root/design/workflows/formulas/README.md, design/workflows/formulas/README.md, scripts/fetch.sh, template/root/scripts/fetch.sh, scripts/README.md, template/root/scripts/README.md, flai/cmd/publish_judgement_test.go, design/conventions/work-management.md, design/conventions/strategic-agents.md, template/root/design/conventions/work-management.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/workflow.md, design/system/pushing-from-the-board.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, docs/contributors/index.md, docs/contributors/template.md]
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
      seconds: 119
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 21
          output: 3811
          cache_read: 2018847
          cache_write: 27783
          cost: 0.7023
    - kind: orchestrator
      seconds: 105
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 33
          output: 316
          cache_read: 1371614
          cache_write: 203475
          cost: 0.4931
cost_of_delay:
  inputs:
    penalty_per_week: 30
    by: orchestrator
    at: 2026-11-17T12:28:37Z
  value: 30
  by: planner-S-0837
  at: 2026-11-17T12:39:39Z
forecast:
  duration: 1h38m
  delivery: 2026-11-17T22:26:00Z
  basis: "Its own forecast of 1h38m; 16th in the pull order with an in-progress limit of 3, behind S-0806, S-0829, S-0834, S-0811, S-0840, S-0809, S-0842, S-0812, S-0808, S-0807, S-0810, S-0241, S-0813, S-0814, S-0815, S-0822, S-0824 and S-0835."
  by: flai
  at: 2026-11-17T12:57:33Z
finalized:
  by: orchestrator
  at: 2026-11-17T12:41:09Z
---
# S-0837 A publish the orchestrator starts is a release_publish call and makes no story, and the baseline publish formula is retired

## Goal

[ADR-0316](../../../design/adrs/0316-a-publish-the-orchestrator-starts-is-a-release-publish-call-and-makes-no-story.md): a publish the orchestrator starts makes no story. Today each one runs the baseline `publish` formula (ADR-0240), which made 51 stories in this repository between 2026-11-07 and 2026-11-17, 131 USD and 2.5 hours of agent time, each taking an in-progress slot, a close-out, a review slot, and an acceptance though its branch changes nothing. The operator asked on 2026-11-17 that this stop. Accept the ADR, or record what the operator changes in it, as this story's first task.

## Acceptance criteria

- [ ] The ADR is accepted, ADR-0240 and ADR-0273 name it as what supersedes them, and no accepted ADR is otherwise edited.
- [ ] With `publish` on and the policy met, the orchestrator's publish is one `release_publish` call: no story, task, branch, or worktree is made, shown by a test of the orchestrator's publish that reads the board before and after.
- [ ] `release_publish` fetches the remote's branch and tags, evaluates the policy again, publishes only when it is met and something is pending, and says which release covers each component when another publish got there first, changing nothing.
- [ ] `release_publish` no longer refuses for a formula named `publish`, flai refuses the orchestrator `flai formula run` for every formula, and the `push` host action still gates the pending publish for an agent flai serve started.
- [ ] The template and this repository no longer hold `design/workflows/formulas/publish.md`; `flai upgrade` removes a project's unchanged copy and reports, without removing, one with project additions; `scripts/fetch.sh` goes with it where nothing else uses it.
- [ ] The orchestrator's role file, in the baseline, the template, and this repository, tells it to publish with `release_publish` only, and the harness definition is generated again from it.
- [ ] The conventions (`work-management.md`, `strategic-agents.md`), the design (`strategic-agents.md`, `flai-cli.md`, `project-manifest.md`), the user and operator guides, and the formulas' README no longer describe a publish as a story, in the template's copies too.
- [ ] S-0811's criterion that names the publish formula is rewritten to speak of `release_publish`, by a thread to the operator or the orchestrator if this story's agent may not edit it.

## Tasks
- T-3600 Accept ADR-0316 as superseding ADR-0240 and ADR-0273, and ask on S-0811 for its publish formula criterion to name release_publish
- T-3601 release_publish fetches, evaluates the policy again, publishes only what is pending, names the release that covers each component, and no longer refuses for a publish formula
- T-3602 flai formula run refuses the orchestrator every formula, and the push host action still gates a served agent's pending publish
- T-3603 The orchestrator's role file, in the baseline, the template, and this repository, publishes with release_publish only, and its definition is generated again
- T-3604 flai upgrade removes a retired template file a project has not changed, and reports without removing one with project additions
- T-3605 The template no longer ships the publish formula or scripts/fetch.sh, and the formulas' and scripts' READMEs no longer list them
- T-3606 The work-management and strategic-agents conventions, here and in the template, describe a publish as a release_publish call, and the template's changelog records the change
- T-3607 The strategic agents, flai CLI, manifest, workflow, and publishing designs, and the user, operator, and contributor guides, describe a publish as a release_publish call

## Notes

- The stopgap until this lands, if the operator wants it sooner: remove `design/workflows/formulas/publish.md` from this repository's main checkout. The orchestrator's role file already says to call `release_publish` where `flai formula list` shows no `publish`, and `release_publish` already works with no formula. Check first that `flai check` and `flai upgrade` do not put the file back or report it missing.
- The archived publish stories stay as they are.

### Planning

The stopgap is applied. Commit 4052fc603 removed `design/workflows/formulas/publish.md` from this repository on 2026-11-17, so only the template's copy is left to delete.

After: none. No other story writes what this one reads. S-0811 (ready) changes the publish formula and `flai/cmd/release.go`, which this story deletes or edits. That is an overlap for the trial merge to judge, not an `after`.

Touches, all files, no folder. `flai touches suggest` listed nothing, since the story declared no touches. Each touch comes from the criteria, the design, and the code layout the explorer read:

| Touch | From | Why |
|-------|------|-----|
| `design/adrs/0316-…md`, `0240-…md`, `0273-…md`, `design/adrs/README.md` | criterion 1 | Accept ADR-0316, set `superseded_by` on the two it supersedes, and regenerate the index. |
| `flai/internal/mcpserver/orchestrate.go`, `orchestrate_test.go` | criteria 2-4, layout | `releasePublish`, `noPublishFormula`, and the `publishProject` fixture the new before-and-after test uses. |
| `flai/internal/hostapi/writes_outcome.go`, `writes_test.go` | criterion 3, layout | `publishArgs` and `publishExits` map only exit 3. The `covered_by` answer needs exit 2. |
| `flai/cmd/formula.go`, `orchestrate_permissions_test.go` | criterion 4, layout | `orchestratorRunsFormula` allows `publish` alone. Its test asserts that. |
| `flai/cmd/release.go`, `release_push_action_test.go`, `docs/users/flai-reference.md` | criterion 4, layout | The `push` gate stays. Its comments and help name the formula. The reference is generated from the help. |
| The three `orchestrator.md` role files, `flai/internal/harness/harness_test.go`, `flai/internal/roles/roles_test.go`, `flai/internal/guard/roles_test.go` | criterion 6, layout | The role's publish entry, paragraph, and step, and the tests that pin its text. |
| `flai/internal/upgrade/upgrade.go`, `upgrade_test.go`, `flai/cmd/upgrade.go`, `upgrade_test.go` | criterion 5, layout | `Compute` never classifies a path the template stopped shipping. Retiring one is new code. |
| `template/root/design/workflows/formulas/publish.md`, both formulas READMEs, both `scripts/fetch.sh`, both scripts READMEs, `flai/cmd/publish_judgement_test.go` | criterion 5, layout | The files that go, the READMEs that list them, and the test that reads the template's formula. |
| Both copies of `work-management.md` and `strategic-agents.md` under the conventions, `template/CHANGELOG.md` | criterion 7, design | The conventions name the formula, and the template's change needs its changelog entry. |
| `design/system/strategic-agents.md`, `flai-cli.md`, `project-manifest.md`, `workflow.md`, `pushing-from-the-board.md` | criterion 7, design | The current design of a publish as a formula's story. |
| `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `docs/contributors/index.md`, `docs/contributors/template.md` | criterion 7, layout | The guides that name the formula or the orchestrator's publish as one. |

Left out: `design/system/dashboard-host-channel.md` and `server-performance.md` name the formula only as history. `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, and a flaiover test also match "publish", but the explorer did not confirm they describe the formula. A task's close widens the touches if one of them changes.

Forecast 1h38m, delivery 2026-11-17T22:40:00Z, as `flai forecast` gives it with the 46 touches. Its first figure was 17m, with no touches counted. Kept as computed. About half the touches are short doc rewrites or deletions, which argues for less. The new retirement path in `flai upgrade` and the fetch and `covered_by` in `release_publish` are new code with tests, which argues for more.

Cost of delay 30.00 USD a week, as `flai cod` gives it from `penalty_per_week: 30`. The orchestrator set that input on TH-1292, as the planner's lower alternative, because the stopgap already stopped this repository's publish stories. What remains is the template and the projects made from it. Kept.
