---
id: S-0220
type: story
nature: feature
title: The orchestrator answers threads, or recommends an answer, as its permission allows
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T00:25:39Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", flai/internal/threads, design/system/strategic-agents.md, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/metrics, design/system/metrics.md, design/adrs, flaiover/src/routes/inbox, flaiover/src/routes/threads, flaiover/src/routes/api/inbox, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 97.56
  by: planner-E-0016
  at: 2026-10-04T04:48:20Z
forecast:
  duration: 2h
  delivery: 2026-10-05T13:06:00Z
  basis: "Its own forecast of 2h; 18th in the pull order with an in-progress limit of 3, behind S-0252, S-0259, S-0249, S-0266, S-0268, S-0253, S-0257, S-0244, S-0262, S-0258, S-0260, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218 and S-0219."
  by: flai
  at: 2026-10-05T00:25:39Z
---
# S-0220 The orchestrator answers threads, or recommends an answer, as its permission allows

## Goal

Agents wait on threads for the operator. With `answer_threads: recommend`, the orchestrator replies to each thread awaiting the operator with its recommended answer, marked as a recommendation, for the operator to confirm; with `autonomous`, it answers and the agent goes on.

## Acceptance criteria
- [ ] In `recommend`, the orchestrator posts a reply tagged as a recommendation that does not set the thread to `answered`, citing the design, ADR, or convention it drew on; the inbox shows the operator a thread with a recommendation to confirm with one action (confirm makes it the answer)
- [ ] In `autonomous`, it answers the thread (setting `answered`) when it can cite a source, and escalates to the operator with a recommendation when it cannot or when the question names the operator's judgement (a decision, a scope change, money)
- [ ] It never resolves a thread it did not open, and never answers a thread opened by itself
- [ ] Each answer is logged with the source cited; the metrics count thread waits ended by the orchestrator separately
- [ ] Tests cover a recommendation, an autonomous answer, and an escalation

## Tasks

## Notes
