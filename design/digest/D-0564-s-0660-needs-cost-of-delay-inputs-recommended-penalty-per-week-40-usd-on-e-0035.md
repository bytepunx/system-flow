---
id: D-0564
title: "S-0660 needs cost-of-delay inputs: recommended penalty_per_week 40 USD on E-0035: Recommended: set a `revenue_per_week` on E-0035, not the planner's `penalty_per_week: 40`"
kind: answer
sources: [TH-0979]
topics: [kluster, deployment, flai, cli, go, code, all]
paths: [flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0564 S-0660 needs cost-of-delay inputs: recommended penalty_per_week 40 USD on E-0035: Recommended: set a `revenue_per_week` on E-0035, not the planner's `penalty_per_week: 40`

The operator, alex, answered TH-0979 at 2026-11-10T09:34:47Z.

What was asked:

```text
S-0660 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0660` cannot work out a value. Those inputs are yours.
```

The answer:

````text
**Recommended:** set a `revenue_per_week` on E-0035, not the planner's `penalty_per_week: 40`. Give the revenue that signet's and authstar's releases still wait on while criteria that need the running system are only checked in review:

```bash
flai edit E-0035 --revenue-per-week <USD>
```

**Why revenue, not a penalty:**

- Your answer on TH-0829 (D-0479) valued E-0032, kluster for signet and authstar, as revenue: `revenue_per_week: 1250`. You said held-up releases of signet and authstar keep you from shipping to waiting subscribers.
- D-0479 also says that an epic which unblocks signet's or authstar's releases should ask for that revenue, not offer a time or penalty figure again.
- E-0035 is that kind of epic. Its outcome is done when a story in a kluster project, such as authstar, reaches review already deployed and checked.

**Why I do not set it:** the input goes on the epic, and under `plan_backlog_stories` I may give only a backlog story its first input. The figure is also yours alone: only you know how much of the 1250 USD a week E-0032 named is still held up, and how much is new.

**Alternative:** if you judge E-0035 holds up no revenue beyond E-0032's, take the planner's `penalty_per_week: 40` on E-0035.

One answer here also serves TH-0980 (S-0661). Once you set the input, I resolve neither thread: the planners write each story's value on their next run.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md
````

Source: TH-0979.
