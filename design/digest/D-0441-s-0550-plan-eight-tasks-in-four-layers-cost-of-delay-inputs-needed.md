---
id: D-0441
title: "S-0550 plan: eight tasks in four layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 10m` on S-0550, as the planner recommends"
kind: answer
sources: [TH-0777, S-0550]
topics: [git, cli, conventions, release, testing, orchestration, planning, flai, go, template, code, all]
paths: [flai/internal/serve, ".gitattributes", design/conventions/git.md, design/system/flai-cli.md, design/adrs/README.md, design/issues/summary.md, template, design/adrs, design/adrs/0157-flai-stream-sync-and-flai-accept-merge-design-adrs-readme-md-s-index-row-by-row.md, flai/internal/check/check.go, flai/internal/check/fix.go, flai/internal/check/fix_test.go, flai/internal/check/issues_test.go, flai/cmd/check.go, flai/cmd/check_test.go, flai/cmd/accept_overlap.go, flai/cmd/accept_overlap_test.go, flai/internal/serve/hold_test.go, template/root/.gitattributes, template/.gitattributes, flai/internal/release/release.go, flai/internal/release/release_test.go, flai/cmd/stream_sync_test.go, template/root/design/conventions/git.md, template/CHANGELOG.md, design/system/workflow.md, flai/internal/adrindex/adrindex.go, flai/internal/adrindex/adrindex_test.go, flai/internal/adr/adr.go, flai/internal/adr/adr_test.go, flai/cmd/adr.go, flai/cmd/adr_test.go, template/root/design/adrs/README.md, flai/internal/issues/generated.go, flai/internal/issues/generated_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/continuous-improvement.md, design/conventions/decisions.md, template/root/design/conventions/decisions.md, docs/operators/settings.md]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0833, TH-0884, TH-0944]
---

# D-0441 S-0550 plan: eight tasks in four layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 10m` on S-0550, as the planner recommends

The operator, alex, answered TH-0777 at 2026-10-31T09:13:39Z.

What was asked:

```text
**Question: cost of delay inputs.** S-0550 and E-0031 have no cost of delay inputs, so `flai cod S-0550` refuses and I wrote no value.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 10m` on S-0550, as the planner recommends:

```bash
flai edit S-0550 --time-lost-per-cycle 10m
```

Then drop "and the digest index" from criterion 2. The digest has no index file.

## The figure

- **The goal is time lost.** It costs about 2.6 extra `stream sync` runs per acceptance, plus the hub-file conflicts worked by hand. That is agent time, not revenue or a penalty.
- **The story gives no reason for an alternative.** No spend comes from these files, so no penalty applies.
- **Put it on the story, not on E-0031**, as on TH-0773 to TH-0775. Inputs on the epic would stop flai taking inputs on any of its stories.

## Why this needs you

I tried to set it, and flai refused: "S-0550 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0550 to ready at 2026-10-31T07:24:05Z. A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.**
  - Criterion 1: T-2394, T-2405, and T-2409.
  - Criterion 2: T-2399.
  - Criterion 3: T-2397.
  - Criterion 4: T-2402.
  - T-2390's ADR comes first because ADR-0157 rejected regenerating the index. T-2413 carries the docs.
- **The assumptions hold.** Stories write under `## Unreleased`, because MD024 would refuse two sibling version headings. T-2399 tests that `merge-tree` honours `merge=union`. Criterion 4 stays unticked at close-out, with T-2402's before figures and method.
- **The forecast and touches fit.** The 5h forecast over flai's 20m states its reason. ADR-0096 §3 narrows the declared `flai/internal/serve` and `template` folders to the files the tasks name.

## The proposals

- **Touches `flai/internal/sync` and `flai/internal/claims`:** the story's agent removes them with `flai touches`. They name no package. Today they also overlap the same two folders in S-0551's touches.
- **Criterion 2's digest clause:** drop it, as above. No digest index exists to mark `merge=union`, and the criterion's words are yours. Left as it is, the clause cannot be ticked honestly at close-out.
- **The `asked` trigger defect:** TH-0771 already reports it. Nothing more is needed here.

The thread stays open for the cost of delay input and the criterion edit.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
````

Source: TH-0777.
