---
id: D-0430
title: An ADR that widens the orchestrator's powers is written proposed while the operator's answer on its plan thread is pending, and accepted on that answer, a confirm_after confirmation included
kind: decision
sources: [S-0524, TH-0695, S-0529, TH-0701]
topics: [orchestration, planning, conventions]
paths: [design/adrs, design/adrs/0222-with-plan-backlog-epics-or-plan-backlog-stories-on-the-orchestrator-may-ask-the.md, design/adrs/0221-while-the-merge-queue-can-batch-a-hand-acceptance-is-refused-until-the-queue.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0430 An ADR that widens the orchestrator's powers is written proposed while the operator's answer on its plan thread is pending, and accepted on that answer, a confirm_after confirmation included

S-0524 built a change that lets the orchestrator start the planner on items planned before. Its story's agent wrote the ADR before the operator had confirmed the plan, and kept it `proposed` until then:

> The ADR was written `proposed` while the operator's confirmation on TH-0695 was pending, since the change widens the orchestrator's powers.

The answer came from `flai serve`, not the operator. It confirmed the orchestrator's recommendation once `orchestration.permissions.confirm_after` had passed, as ADR-0206 allows. The agent took that confirmation as the answer and accepted the ADR:

> TH-0695 settled: flai serve confirmed the orchestrator's recommendation to approve the remedy at 2026-10-30T14:44:33Z, after 48h with no answer from the operator (orchestration.permissions.confirm_after, ADR-0206), and the orchestrator resolved the thread. ADR-0222 is accepted on that answer. The operator can still send S-0524 back at review if they would have chosen otherwise.

S-0529 did the same for a fix that limits the operator's own acceptance rather than widening the orchestrator's powers. Its planner proposed the fix on TH-0701, and the orchestrator recommended it there because it reverses a consequence of ADR-0174 and binds the operator. The story's agent built it before the answer came:

> Build the fix the planner proposed and the orchestrator recommends on TH-0701, whose confirmation by the operator is still pending. The ADR is written `proposed` until the operator confirms. If the operator picks another fix there, T-2263 and T-2264 are redone and the ADR is rejected.

`flai serve` confirmed that recommendation after `confirm_after` (48h) at 2026-10-30T18:19:21Z, and T-2346 accepted ADR-0221 on it.

So an ADR whose decision waits on the operator's answer on its plan thread stays `proposed` until that thread is answered, whether it widens an agent's powers or limits the operator's. A confirmation after `confirm_after` counts as that answer, and review is the operator's last chance to decline it. Should the operator pick another fix, the ADR is set `rejected` (ADR-0223) and the work redone.

Source: the Decisions of wip/archive/agents/S-0524.md and TH-0695; the Decisions of wip/archive/agents/S-0529.md and TH-0701.
