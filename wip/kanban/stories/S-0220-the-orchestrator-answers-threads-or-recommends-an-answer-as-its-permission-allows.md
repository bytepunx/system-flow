---
id: S-0220
type: story
nature: feature
title: The orchestrator answers threads, or recommends an answer, as its permission allows
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-02T11:54:43Z
transitions: []
tags: [flai]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", flai/internal/threads, design/system/strategic-agents.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
