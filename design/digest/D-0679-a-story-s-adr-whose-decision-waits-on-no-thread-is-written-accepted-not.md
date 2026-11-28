---
id: D-0679
title: "A story's ADR whose decision waits on no thread is written accepted, not proposed: the operator reviews it when they accept the story, and a proposed ADR would leave the code with no decision to follow"
kind: decision
sources: [S-0846, S-0261]
topics: [conventions, design, flai]
paths: [design/adrs, design/conventions/decisions.md]
date: 2026-11-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0679 A story's ADR whose decision waits on no thread is written accepted, not proposed: the operator reviews it when they accept the story, and a proposed ADR would leave the code with no decision to follow

S-0846 wrote ADR-0320, on how the auto-upgrade drains close-out jobs, with status `accepted` on its own branch, and gave its reason:

> ADR-0320 is written accepted rather than proposed: the operator's acceptance of S-0846 is where it is reviewed, and a proposed ADR would leave the code without a decision to follow.

S-0261 did the same with ADR-0104: "Accepted rather than proposed, since the plan asked for a thread only if the fix overrode TH-0032."

So when no thread holds an operator's answer the decision depends on, the story's agent writes its ADR `accepted` and the code follows it. The operator reviews the ADR with the story's diff at acceptance, and sends the story back if they would decide otherwise. When an answer on a thread is still pending, the other rule applies: the ADR stays `proposed` until that answer comes (D-0430). Until the ADR reaches main it is the story's own draft, amended in place (D-0591).

Sources: the Decisions of wip/archive/agents/S-0846.md and wip/archive/agents/S-0261.md.
