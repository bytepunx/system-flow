---
id: D-0508
title: "S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show: Recommended: keep S-0562 a draft until you answer TH-0811"
kind: answer
sources: [TH-0839, TH-0811]
topics: [orchestration, planning, releases, all]
paths: [system-flow.yaml, wip/kanban/stories/S-0562-the-publish-formula-is-run-once-in-this-repository-its-instance-worked-to.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0508 S-0562 plan: three tasks in three layers; cost of delay inputs needed, and TH-0811 decides what criterion 1 can show: Recommended: keep S-0562 a draft until you answer TH-0811

Changing this project's release policy is the operator's alone: the orchestrator does not reverse the operator's choice (D-0216), and a draft whose criteria depend on the policy stays a draft until the operator answers.

planner-S-0562 asked on TH-0839 for S-0562's cost of delay and for a switch to `policy: threshold`, since under `judgement` the publish formula's step 2, `flai release --evaluate --require-met`, always stops. The orchestrator set `penalty_per_week: 25` on S-0562, the figure S-0545 and S-0582 carry for a proof run that blocks no other work, and recommended: "keep S-0562 a draft until you answer TH-0811 ... the release policy belongs to you. On TH-0256 you chose `judgement` over `threshold` (D-0216) ... Only you can do that."

The operator, alex, confirmed that recommendation at 2026-11-07T09:40:50Z, and had answered TH-0811 with `threshold` (D-0504), so S-0562's criterion 1 stood as written and the story went ahead.

Sources: TH-0839, TH-0811.
