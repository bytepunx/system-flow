---
id: I-0188
title: flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again
class: efficiency
status: closed
count: 1
first_reported: 2026-10-21T16:01:33Z
last_reported: 2026-10-21T16:01:33Z
updated: 2026-10-30T14:25:58Z
---

# I-0188 flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again

## Description
flai plan --candidates lists no epic or story planned before, so the orchestrator cannot ask the planner to plan one again

## Instances

### 2026-10-21T16:01:33Z
On TH-0579 the operator asked that the orchestrator be allowed to start the planner on epics and stories already planned. The operator confirmed a new E-0024 story after S-0419 that generates the built-in roles' definitions, which only E-0024's planner can draft, but flai plan --candidates lists an epic only in the backlog with no stories, or with every story done or cancelled, and the plan tool refuses the orchestrator anything it does not list. The new story waits for the operator to run flai plan E-0024.

## Remediation

Story S-0524 remediates this issue, created from it at 2026-10-28T12:53:50Z.
Closed 2026-10-30T14:25:58Z: S-0524 fixed it under ADR-0222. flai plan --candidates now lists a third kind of candidate, planned before and asked again: an epic not done or cancelled, or a story in the backlog or ready, whose newest planner run ended and on which, or on an item under which, the operator has written since. It names that thread. The orchestrator asks the planner for it under plan_backlog_epics or plan_backlog_stories when the operator's entry asks for a change to the plan, and the run records orchestrator; TH-nnnn.
