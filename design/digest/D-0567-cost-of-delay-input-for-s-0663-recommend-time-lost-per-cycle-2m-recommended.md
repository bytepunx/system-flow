---
id: D-0567
title: "Cost of delay input for S-0663: recommend time_lost_per_cycle 2m: Recommended: take the thread's alternative, one input on E-0035, not `time_lost_per_cycle: 2m` on S-0663"
kind: answer
sources: [TH-0981]
topics: [kluster, deployment, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/kluster/deploy.go, flai/internal/kluster/deploy_test.go, flai/internal/serve/kluster.go, flai/internal/serve/kluster_test.go, flai/cmd/kluster.go, flai/cmd/kluster_test.go, docs/users/flai-reference.md, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/kluster_forward_test.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/server/repo.ts, flaiover/src/lib/server/agent.ts, flaiover/src/routes/api/kluster/checks/+server.ts, flaiover/src/routes/api/kluster/checks/checks.test.ts, flaiover/src/routes/api/kluster/deployments.test.ts, flaiover/src/routes/api/kluster/forward/forward.test.ts, flaiover/src/lib/components/DeploymentRow.svelte, flaiover/src/lib/components/DeploymentRow.svelte.test.ts, design/components/flaiover/src/lib/components/DeploymentRow.examples.json, flaiover/src/routes/ops/deployments/+page.svelte, flaiover/src/routes/ops/deployments/deployments.svelte.test.ts, design/system/kluster.md, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0567 Cost of delay input for S-0663: recommend time_lost_per_cycle 2m: Recommended: take the thread's alternative, one input on E-0035, not `time_lost_per_cycle: 2m` on S-0663

The operator, alex, answered TH-0981 at 2026-11-10T09:38:00Z.

What was asked:

```text
S-0663 has no cost of delay inputs, and neither has its epic E-0035, so `flai cod S-0663` cannot work out a value.
```

The answer:

````text
**Recommended:** take the thread's alternative, one input on E-0035, not `time_lost_per_cycle: 2m` on S-0663. Answer TH-0979 with a `revenue_per_week` on E-0035. S-0663 then takes its share by forecast and needs no input of its own.

**Why:**

- D-0479 (TH-0829): you valued E-0032 as revenue, 1250 USD/week, for the signet and authstar releases it holds up. It says an epic that unblocks those releases should ask for that revenue, not a time figure.
- E-0035 is such an epic. All five of its stories, S-0659 to S-0663, lack inputs, and the planners of S-0659 and S-0662 now recommend the epic too.
- A 2m on S-0663 alone would value it at about 5 USD/week. That figure does not match how its siblings are valued.

**Why I do not set the 2m:** this thread itself gives a reason for the epic alternative. An input on the epic is not mine to set. It is money, so it is yours alone.

**Alternative:** if you want S-0663 valued apart from the epic, set the planner's figure:

```bash
flai edit S-0663 --time-lost-per-cycle 2m
```

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md
````

Source: TH-0981.
