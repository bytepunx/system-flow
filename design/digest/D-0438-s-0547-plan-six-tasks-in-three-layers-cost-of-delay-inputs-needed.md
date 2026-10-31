---
id: D-0438
title: "S-0547 plan: six tasks in three layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 3h` on S-0547, as the planner recommends"
kind: answer
sources: [TH-0774]
topics: [testing, cli, git, orchestration, planning, flai, go, code, all]
paths: [flai/internal/check, flai/internal/verify, flai/cmd, system-flow.yaml, design/system/flai-cli.md, design/conventions/testing.md, docs/users/flai.md, design/adrs, design/adrs/README.md, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/operators/settings.md, docs/users/flai-reference.md, design/system/project-manifest.md, design/conventions/code-quality.md, design/conventions/tooling.md, design/issues/I-0221-the-close-out-s-go-test-tier-runs-only-the-changed-packages-so-a-change-that-breaks-a-caller-package-s-tests-passes-before-review.md, design/issues/I-0203-a-story-the-operator-accepts-outside-the-merge-queue-lands-without-the-queue-tiers-so-a-stale-generated-file-on-main-fails-every-later-batch-at-integration.md, design/issues/summary.md]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0438 S-0547 plan: six tasks in three layers; cost of delay inputs needed: Recommendation: set `time_lost_per_cycle: 3h` on S-0547, as the planner recommends

The operator, alex, answered TH-0774 at 2026-10-31T09:08:57Z.

What was asked:

```text
**Question: cost of delay inputs.** Neither S-0547 nor E-0031 has cost of delay inputs, so `flai cod S-0547` refuses and I wrote no value.
```

The answer:

````text
**Recommendation:** set `time_lost_per_cycle: 3h` on S-0547, as the planner recommends:

```bash
flai edit S-0547 --time-lost-per-cycle 3h
```

## The figure

- **The goal is time lost.** About 7 send-backs a week, each about 25 minutes of a restarted agent's fix, another close-out, a re-enqueue, and another verification. That is agent time, not revenue or a penalty.
- **On the story, not E-0031.** The epic's stories each measure a different loss (TH-0772, TH-0773). Inputs on the epic would also stop flai taking inputs on any of its stories.

## Why this needs you

I tried to set it, and flai refused: "S-0547 is ready, not in the backlog: the orchestrator gives cost of delay inputs to a story in the backlog alone". You moved S-0547 to ready at 2026-10-31T07:24:13Z. A confirmation of this entry does not write the input. Running the command does.

## The plan: approved

- **The tasks fit the criteria.** T-2377 covers criterion 1, T-2376 and T-2378 criterion 2, T-2379 criterion 3, T-2375 and T-2380 criterion 4. The three layers follow the real dependencies.
- **Criterion 1 through `{importers}`: right.** Changing what `{packages}` means would change every project's tiers. The criterion's outcome, changed packages and their importers tested, is met.
- **`when: ["flai/**", "template/**"]` for integration and smoke: right.** Criterion 3 asks for exactly the flai and template paths.
- **Touches and forecast fit.** The 4h forecast over flai's 18m states its reason. No `after` is needed: the overlap hold on S-0548 and S-0552 already orders it.

The thread stays open for the cost of delay input only.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
````

Source: TH-0774.
