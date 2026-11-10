---
id: D-0566
title: "S-0661 needs cost-of-delay inputs: recommended one input on E-0035, as TH-0979 asks: Recommended: take this thread's alternative, a `revenue_per_week` on E-0035, with the figure yours"
kind: answer
sources: [TH-0980]
topics: [kluster, deployment, flai, cli, go, code, all]
paths: [flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/internal/serve/checks.go, flai/internal/serve/checks_test.go, flai/internal/kluster/request.go, flai/internal/kluster/request_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, flai/cmd/checks.go, flai/cmd/checks_test.go, flai/internal/mcpserver/kluster.go, flai/internal/mcpserver/kluster_test.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/server_test.go, flai/internal/guard/guard_tables.go, flai/internal/guard/guard_rules.go, flai/internal/guard/kluster_test.go, design/system/kluster.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0566 S-0661 needs cost-of-delay inputs: recommended one input on E-0035, as TH-0979 asks: Recommended: take this thread's alternative, a `revenue_per_week` on E-0035, with the figure yours

The operator, alex, answered TH-0980 at 2026-11-10T09:36:54Z.

What was asked:

```text
S-0661 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0661` cannot work out a value. Those inputs are yours.
```

The answer:

````text
**Recommended:** take this thread's alternative, a `revenue_per_week` on E-0035, with the figure yours. I recommend the same on TH-0979, so one answer serves S-0660 and S-0661:

```bash
flai edit E-0035 --revenue-per-week <USD>
```

**Why:**

- Your answer on TH-0829 (D-0479) valued E-0032, kluster for signet and authstar, as `revenue_per_week: 1250`, not as time lost. It says an epic that unblocks signet's or authstar's releases should ask for that revenue.
- E-0035 is such an epic: it lets a story in authstar or signet reach review already deployed and checked.

**Why I do not set it:** the input goes on the epic, which `plan_backlog_stories` does not let me set. The figure is yours alone: only you know how much revenue is still held up beyond E-0032's 1250 USD a week.

**Alternative:** `penalty_per_week: 40` on E-0035, as TH-0979's planner recommends, if you judge E-0035 holds up no new revenue.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md
````

Source: TH-0980.
