---
id: D-0409
title: "Plan for S-0502: confirm_after in eight tasks over five layers, with an operator-alone mark for what a timeout must never confirm: Recommended for the operator: `confirm_after: 48h` for this project…"
kind: answer
sources: [TH-0658]
topics: [orchestration, cli]
paths: [system-flow.yaml, design/adrs/0206-flai-serve-confirms-an-orchestrator-recommendation-the-operator-left-unanswered.md, design/system/project-manifest.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0409 Plan for S-0502: confirm_after in eight tasks over five layers, with an operator-alone mark for what a timeout must never confirm: Recommended for the operator: `confirm_after: 48h` for this project…

The operator, alex, confirmed on TH-0658 at 2026-10-28T20:27:20Z the orchestrator's recommendation of 2026-10-26T05:14:17Z: S-0502's plan stands as drafted, and this project sets `orchestration.permissions.confirm_after: 48h`.

The planner asked: "What value should this project set for `confirm_after` (T-2148)? Recommended: `48h`, the figure you used on TH-0598. Alternatives: `24h` ... or leave it unset here and ship the setting for other projects only."

The orchestrator would not answer it itself: "`confirm_after` decides how long before flai serve confirms my own recommendations in your name. That is a decision about your authority, not one a source settles, so it is yours. 24h would have ended TH-0583 and TH-0588 sooner, but it leaves you one day to overrule. 48h matches what you chose."

The setting is in `system-flow.yaml`; what it does is ADR-0206. A later change to the figure is the operator's to make, for the same reason.

Source: TH-0658.
