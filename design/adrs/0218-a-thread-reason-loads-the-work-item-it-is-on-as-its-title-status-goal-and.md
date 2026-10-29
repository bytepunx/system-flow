---
id: ADR-0218
title: "A thread reason loads the work item it is on as its title, status, goal, and acceptance criteria, and a reason pack's digest is capped at a sixteenth of the budget"
status: accepted
date: 2026-10-29
supersedes: []
superseded_by: []
refines: [ADR-0194, ADR-0217]
---

# ADR-0218 A thread reason loads the work item it is on as its title, status, goal, and acceptance criteria, and a reason pack's digest is capped at a sixteenth of the budget

## Context

[ADR-0217](0217-a-pack-primed-for-the-orchestrator-s-reasons-loads-the-role-s-conventions-with.md) §3 loads, for a `thread` reason, the thread and the item it is on, whole. Its §4 takes the digest's answers within an eighth of the budget, as [ADR-0194](0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md) does for every pack.

Measured on 2026-10-29, packs for real threads here were 31 to 33% of the general orchestrate pack of 176 KB. S-0537's second criterion asks for under a quarter. A pack for TH-0719 was 58.1 KB:

| Part | Size |
|---|---|
| conventions, cut by reason | 27.4 KB |
| thread | 6.3 KB |
| its story, whole | 11.0 KB |
| digest answers | 8.4 KB |
| header | about 5 KB |

A planned story's notes, tasks, and planning are most of its size, and an answer to a thread on it rarely needs them. On TH-0753 the operator confirmed trimming both parts.

## Decision

**A `thread` reason loads the work item it is on as its title, status, goal, and acceptance criteria, and a reason pack's digest is capped at a sixteenth of the budget.** This refines ADR-0217 §3 and §4, and ADR-0194 for the reason pack alone; the rest of each stands.

1. **A thread's item.** The pack gives its ID, title, type, and status, and its `## Goal` and `## Acceptance criteria` sections. It leaves out the rest of its body: notes, tasks, and planning. `item_get` gives the item whole. A `review`, `planned`, or `draft` reason still loads its items whole, since deciding on them needs the whole item.
2. **Digest.** A reason pack takes the digest's answers and recommendations ranked against the reasons' threads within a sixteenth of the budget, about 5 KB at the default 80 KB. Every other pack keeps ADR-0194's eighth.
3. **The measure.** S-0537's test measures a pack for a copy of a real thread of this repository on its story, not a fixture of one entry.

## Consequences

- A pack for a real thread here falls under a quarter of the general pack.
- The orchestrator calls `item_get` when a thread's answer turns on the item's notes or tasks.
- Fewer digest answers reach a reason pack. Those ranked highest against its threads come first, so the cap drops the weakest.

## Alternatives considered

- **Keep the pack as it was, and tick the criterion on the fixture.** That meets the criterion's words but not its purpose: a run for one thread reading much less than the general pack.
- **Cut the conventions further.** The remaining sections, such as safety, communication, and the orchestrator's rules for threads, are the ones a thread's answer rests on.
