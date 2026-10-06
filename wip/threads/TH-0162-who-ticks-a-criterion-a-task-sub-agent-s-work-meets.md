---
id: TH-0162
title: Who ticks a criterion a task sub-agent's work meets?
anchor:
  path: wip/kanban/stories/S-0282-story-s-should-correctly-tick-acceptance-criteria-when-they-have-verified-them.md
  item: S-0282
status: resolved
participants: [agent-S-0282, alex]
created: 2026-10-06T03:46:58Z
updated: 2026-10-06T04:39:28Z
---

# TH-0162 Who ticks a criterion a task sub-agent's work meets?

On wip/kanban/stories/S-0282-story-s-should-correctly-tick-acceptance-criteria-when-they-have-verified-them.md.

## Entries

### 2026-10-06T03:46:58Z agent-S-0282
The goal says a task sub-agent completing a criterion should tick it. Today the guard refuses every flai write by a sub-agent (ADR-0059, ADR-0060), and work-management.md says never to tick a criterion you did not verify.

Recommendation: the story's agent ticks. A task sub-agent names in its final message the criteria its task meets. The story's agent reviews the work, then ticks each criterion it verified with the new command (`flai criteria tick S-nnnn <n>`, MCP `criteria_tick`). The guard keeps refusing sub-agents' writes; `flai criteria list` stays a read they may run. The conventions and flai serve's start prompt say this.

Alternative: the guard lets a task sub-agent run `flai criteria tick` on its own story only. That needs a new ADR beside ADR-0060, and the guard must tell a task sub-agent from an explorer or verifier.

I am building the command, the MCP tool, and the host API action now. Only the guard and convention wording wait on this answer.

### 2026-10-06T04:39:19Z alex
take recommendation

### 2026-10-06T04:39:28Z agent-S-0282
Resolved.
