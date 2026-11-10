---
id: D-0568
title: "S-0659 cost of delay: give E-0035 a revenue_per_week input; I recommend 1250 USD/week, as E-0032's: Recommended: set a `revenue_per_week` on E-0035"
kind: answer
sources: [TH-0982]
topics: [kluster, deployment, flai, cli, go, code, all]
paths: [design/system/kluster.md, design/adrs/README.md, design/adrs/0254-flai-runs-one-kluster-per-project-from-the-config-its-manifest-s-kluster-block.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0568 S-0659 cost of delay: give E-0035 a revenue_per_week input; I recommend 1250 USD/week, as E-0032's: Recommended: set a `revenue_per_week` on E-0035

The operator, alex, answered TH-0982 at 2026-11-10T09:38:37Z.

What was asked:

```text
`flai cod S-0659` cannot work out a value: neither S-0659 nor its epic E-0035 has cost of delay inputs, and they are yours to give.
```

The answer:

````text
**Recommended:** set a `revenue_per_week` on E-0035. Take the planner's 1250 unless less of E-0032's held-up revenue is still at stake:

```bash
flai edit E-0035 --revenue-per-week 1250
```

**Why:**

- D-0479 (TH-0829): you valued E-0032 as revenue, 1250 USD/week, because held-up releases of signet and authstar keep you from shipping to waiting subscribers. It says an epic that unblocks those releases should ask for that revenue.
- E-0035 is such an epic. S-0659 is its design story.
- I recommend the same input on TH-0979, TH-0980, TH-0981, and TH-0983. One answer on E-0035 serves all five stories. Each planner then writes its story's share by forecast.

**Why I do not set it:** the input goes on the epic. Under `plan_backlog_stories` I may give only a backlog story its first input. The figure is money and yours alone.

**Alternative:** `penalty_per_week: 40` on E-0035, as TH-0979's planner first proposed, if E-0035 holds up no revenue beyond E-0032's.

S-0659 stays a draft until then: `flai promote --drafts` lists it as lacking only a cost of delay value.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/digest/D-0479-e-0032-cost-of-delay-inputs-so-its-stories-can-be-valued-i-entered-my-own.md
````

Source: TH-0982.
