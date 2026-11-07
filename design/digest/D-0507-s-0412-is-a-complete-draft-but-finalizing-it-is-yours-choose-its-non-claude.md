---
id: D-0507
title: "S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3): We will only be using OpenRouter fixtures, no direct…"
kind: answer
sources: [TH-0533]
topics: [agents, cli, template, flai, go, code, all]
paths: [design/experiments/S-0412-a-story-is-worked-to-review-end-to-end-on-the-opencode-adapter-by-a-scripted.md, design/experiments/README.md, design/issues, design/issues/summary.md, design/system/agent-adapters.md, docs/operators/index.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0533, TH-0866, TH-0877, TH-0878]
---

# D-0507 S-0412 is a complete draft, but finalizing it is yours: choose its non-Claude model and the OpenRouter key's limit (TH-0467 question 3): We will only be using OpenRouter fixtures, no direct…

No live OpenRouter integration is part of this repository for now: S-0412's end-to-end run on the opencode adapter replays a scripted session in OpenRouter's format from fixtures, with no funded key and no non-Claude model chosen.

The operator, alex, answered TH-0533 at 2026-11-07T09:28:49Z: "We will only be using OpenRouter fixtures, no direct integration will be part of this repository for now".

What was asked: the orchestrator left S-0412 a draft because its plan needed the operator to set its `agent.model` to a non-Claude model and put an OpenRouter key with a 5 USD limit in the host's `providers` map (TH-0467 question 3). planner-E-0023 then rewrote S-0412 for fixtures alone (TH-0866).

See D-0506 for the same answer on installing OpenCode.

Sources: TH-0533, TH-0467.
