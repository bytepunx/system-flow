---
id: D-0508
title: "S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show: Recommended: keep S-0562 a draft until you answer TH-0811"
kind: answer
sources: [TH-0839]
topics: [orchestration, cli, planning, template, conventions, all]
paths: [wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md, design/issues, design/issues/summary.md, design/issues/I-0283-a-publish-instance-with-nothing-to-release-stops-at-its-evaluate-step-down-the.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0508 S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show: Recommended: keep S-0562 a draft until you answer TH-0811

The operator, alex, answered TH-0839 at 2026-11-07T09:40:50Z.

What was asked:

```text
**Questions for you, each with my recommendation first:**
```

The answer:

```text
**Recommended:** keep S-0562 a draft until you answer TH-0811. Then word criterion 1 to match that answer before you finalize it. I will not finalize S-0562 while TH-0811 is open, even though it will be complete once the planner writes its value.

**Why it is your call:** the release policy belongs to you. On TH-0256 you chose `judgement` over `threshold` (D-0216). The planner recommends switching to `policy: threshold`, which would reverse that choice. Only you can do that. Your two options:

- **Switch to `threshold`** on TH-0811. The instance then releases for real, and criterion 1 stays as written.
- **Keep `judgement`.** Then criterion 1 has to be reworded so that the judgement stop at step 2 counts, because under judgement `--require-met` always stops. In that case the run shows only the stop.

T-2504 runs nothing until TH-0811 is answered, so nothing is lost while you decide.

Source: design/digest/D-0216-orchestration-release-sets-value-1000-but-no-policy-so-flai-evaluates-it-as.md
```

Source: TH-0839.
