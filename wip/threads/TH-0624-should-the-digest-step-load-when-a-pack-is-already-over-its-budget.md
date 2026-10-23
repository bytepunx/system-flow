---
id: TH-0624
title: Should the digest step load when a pack is already over its budget?
anchor:
  path: wip/kanban/stories/S-0430-flai-prime-gives-a-story-s-agent-the-digest-entries-its-topics-and-touches-select-and-the-orchestrator-the-entries-on-the-threads-it-answers-ranked-within-the-budget-and-marked-as-the-digest-s.md
  item: S-0430
status: resolved
participants: [agent-S-0430, orchestrator, alex]
created: 2026-10-23T19:30:11Z
updated: 2026-10-23T20:47:07Z
---

# TH-0624 Should the digest step load when a pack is already over its budget?

On wip/kanban/stories/S-0430-flai-prime-gives-a-story-s-agent-the-digest-entries-its-topics-and-touches-select-and-the-orchestrator-the-entries-on-the-threads-it-answers-ranked-within-the-budget-and-marked-as-the-digest-s.md.

## Entries

### 2026-10-23T19:30:11Z agent-S-0430
Every pack on this repository is over its budget before the digest step runs. The story packs I measured (S-0250, S-0300, S-0350, S-0400, S-0428, S-0429, S-0468) are about 145 KB of conventions against 80 KB (`exceeded: conventions`). The orchestrator's pack is 143 KB (`exceeded: briefs`). ADR-0194, as written on this branch, loads nothing on a pack that is over its budget, as the ranked step does. So no agent here would get a digest entry, and the story's goal would not be met on this repository.

**Recommendation: A.** Run the digest step whatever `exceeded` says, within its own share, an eighth of the budget (10 KB of 80 KB), after the briefs. This is how every brief is kept when a pack is over its budget (TH-0032). A pack within its budget still never goes over because of the step. A pack already over grows by at most that eighth, and parts (ADR-0104) carry it. I would amend ADR-0194 on this branch, which no other story has seen yet, and test both cases.

- **B.** Keep ADR-0194 as written. The digest loads only once the conventions fit, by narrowing their topics or raising `prime.budget`. Here that means no pack gets entries until that happens.
- **C.** A, and also raise this project's `prime.budget`. That is a manifest change of yours, separate from this story.

With A, I measure criterion 6 on three archived stories with and without the step. With B, I record that the step loads nothing at 80 KB, measure at a budget the conventions fit, and keep the code as it is.

### 2026-10-23T19:48:59Z orchestrator (recommendation)
Recommended for the operator: **A.** The digest step loads within its own share, an eighth of the budget, whatever `exceeded` says. ADR-0194 is amended on S-0430's branch, with tests for a pack under its budget and a pack over it.

**Why A.**
- E-0025's outcome, which you confirmed on TH-0476 and TH-0480, is that past answers reach the agent that meets the question again. Under B no pack on this repository gets an entry, because the conventions alone are about 145 KB against 80 KB. S-0431 to S-0435 would then build on a step that does nothing here.
- ADR-0049 already exempts what the agent must see from the budget: conventions and named documents are never cut, and briefs are kept when the pack is over. A digest entry is the same kind of thing, a short pointer that stops a repeat question. The ranked step is different: it is filler and rightly yields.
- The cost is bounded. A pack grows by at most 10 KB, and parts (ADR-0104) carry it.

**Why this is yours.** It changes what every pack holds, and ADR-0194 has not been accepted yet, so it is a decision not yet recorded. C, raising `prime.budget`, is a separate manifest change and stays open either way.

Source: design/adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md § Decision

### 2026-10-23T19:50:18Z agent-S-0430
Measured since, with and without the step (this branch's flai against main's):

| Pack | Budget | Over budget on | Size | Digest entries |
|------|--------|----------------|------|----------------|
| S-0429, S-0400, S-0350 | 80 KB | conventions | 145–146 KB | 0 |
| S-0429, S-0400, S-0350 | 240 KB | briefs | 312–349 KB | 0 |
| orchestrator | 80 KB | briefs | 144 KB | 0 |
| orchestrator | 240 KB | nothing | 175 KB, from 144 KB | 16 |

So under B, a story pack here gets no entry even at 240 KB. Its briefs alone pass that budget, so raising `prime.budget` (C) would not be enough for stories either. A is still my recommendation. Under A, each pack above would carry its digest within 10 KB at 80 KB, or 30 KB at 240 KB.

### 2026-10-23T20:30:50Z alex
A

### 2026-10-23T20:47:07Z agent-S-0430
Resolved: Answered A and applied in T-2090: the digest step loads within its eighth on a pack already over its budget. ADR-0194 is amended on story/S-0430, and agent-context.md records the measurement.
