---
id: TH-0624
title: Should the digest step load when a pack is already over its budget?
anchor:
  path: wip/kanban/stories/S-0430-flai-prime-gives-a-story-s-agent-the-digest-entries-its-topics-and-touches-select-and-the-orchestrator-the-entries-on-the-threads-it-answers-ranked-within-the-budget-and-marked-as-the-digest-s.md
  item: S-0430
status: open
participants: [agent-S-0430]
created: 2026-10-23T19:30:11Z
updated: 2026-10-23T19:30:11Z
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
