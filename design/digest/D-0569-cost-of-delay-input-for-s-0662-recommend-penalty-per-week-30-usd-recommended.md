---
id: D-0569
title: "Cost of delay input for S-0662: recommend penalty_per_week 30 USD: Recommended: take the planner's revised recommendation"
kind: answer
sources: [TH-0983]
topics: [kluster, deployment, review, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/kluster/evidence.go, flai/internal/kluster/evidence_test.go, flai/cmd/move.go, flai/cmd/move_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/review_test.go, flai/internal/verify/story.go, flai/internal/verify/kluster.go, flai/internal/verify/kluster_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/hostapi/writes_host.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/review.ts, flaiover/src/lib/review.test.ts, flaiover/src/lib/components/Review.svelte, flaiover/src/lib/components/Review.svelte.test.ts, flai/internal/preview/orchestrator.go, flai/internal/preview/evidence_test.go, flai/cmd/accept.go, flai/cmd/accept_orchestrator_test.go, flai/internal/context/story.go, flai/internal/context/context.go, flai/internal/context/context_test.go, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/workflow.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, flai/internal/roles/baseline/orchestrator.md, template/CHANGELOG.md, design/system/kluster.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0569 Cost of delay input for S-0662: recommend penalty_per_week 30 USD: Recommended: take the planner's revised recommendation

The operator, alex, answered TH-0983 at 2026-11-10T09:39:48Z.

What was asked:

```text
**Recommended:** `penalty_per_week: 30` USD on S-0662.
```

The answer:

````text
**Recommended:** take the planner's revised recommendation. Answer TH-0979 with a `revenue_per_week` on E-0035, not `penalty_per_week: 30` on S-0662. S-0662 then takes its share by forecast and needs no input of its own.

**Why:**

- D-0479 (TH-0829): you valued E-0032 as revenue, 1250 USD/week, for the signet and authstar releases it holds up. It says an epic that unblocks those releases should ask for that revenue.
- E-0035 is such an epic. One input on it serves S-0659 to S-0663, and TH-0979 to TH-0982 ask the same.

**Why I do not set the 30:** the planner withdrew it in favour of the epic. An input on the epic is not mine to set. It is money, so it is yours alone.

**Alternative:** if you want S-0662 valued apart from the epic, set the planner's fallback:

```bash
flai edit S-0662 --penalty-per-week 30
```

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md
````

Source: TH-0983.
