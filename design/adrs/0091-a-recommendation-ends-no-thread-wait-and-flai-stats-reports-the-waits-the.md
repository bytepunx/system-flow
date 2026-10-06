---
id: ADR-0091
title: "A recommendation ends no thread wait, and flai stats reports the waits the orchestrator ended and those the operator's confirmation of a recommendation ended"
status: proposed
date: 2026-10-06
supersedes: []
superseded_by: []
refines: [ADR-0081]
topics: [cli, orchestration, analysis]
---

# ADR-0091 A recommendation ends no thread wait, and flai stats reports the waits the orchestrator ended and those the operator's confirmation of a recommendation ended

## Context

`flai stats` measures how long a story's agent waited on its threads (ADR-0081, S-0205): a thread's wait runs from its first entry to the first later entry by another author, and `items[].wait_threads_seconds` and `waiting.weeks[].threads` sum those waits over the time the item was in progress. S-0220 lets the orchestrator reply on threads that await the operator, as `orchestration.permissions.answer_threads` allows: with `recommend` its reply is a recommendation, which leaves the thread awaiting the operator until they confirm it, and with `autonomous` it answers when it can cite a source (ADR-0090). Read as before, a recommendation would end the agent's wait although nobody has answered, and the waits the orchestrator ends could not be told from the operator's, so nobody could see how much waiting the orchestrator takes off the operator. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, and changes only with an ADR.

## Decision

A recommendation ends no thread wait, and each wait records whether the orchestrator's answer or the operator's confirmation of a recommendation ended it, reported in new fields beside the existing ones.

1. **A recommendation ends no wait.** A thread's wait ends at the first later entry by an author other than the opener that is not a recommendation. While the only replies are recommendations, the wait runs on to the operator's confirmation, to the next answer, or, with neither, to `updated` once resolved and to now while open, as before.
2. **Each wait records who ended it.** The orchestrator, when the entry that ended it is by `orchestrator`, the author name of the orchestrator's run (ADR-0087); confirmed, when that entry confirms a recommendation, as ADR-0090 writes a confirmation; anyone else otherwise, and nobody while the wait is not ended.
3. **`items[].wait_threads_orchestrator_seconds`** is the seconds of the union of the item's waits the orchestrator ended that fall in its `in-progress` intervals, computed as `wait_threads_seconds` is but over those waits alone. It is never more than `wait_threads_seconds`; time in which such a wait overlaps a wait someone else ended counts in it. It is present, 0 or more, whenever `wait_threads_seconds` is.
4. **`waiting.weeks[].threads` gains `orchestrator` and `confirmed`**, each with `count`, the number of the waits of the week's items that the orchestrator, or a confirmation, ended and that have some part in the item's `in-progress` intervals, and `total_seconds`, the sum over the week's items of the seconds of the union of those waits that fall in the item's `in-progress` intervals; for `orchestrator`, the sum of `wait_threads_orchestrator_seconds`. Both are present in every week, 0 without such waits.
5. **The existing fields keep their meaning.** `wait_threads_seconds` and the week's `threads.total_seconds` and `threads.mean_seconds` still cover every wait, whoever ended it; only a wait a recommendation would have ended is now longer.

## Consequences

- The dashboard reads its waiting charts as before and can show, beside them, how much of the agents' waiting the orchestrator ended and how many waits the operator ended by confirming a recommendation.
- A thread with no recommendation and no orchestrator entry gives the same figures as before; the new fields are 0 for it.
- A wait on a thread the orchestrator recommended an answer on runs to the operator's confirmation, which is the time the agent was in fact left waiting.
- The orchestrator's and the confirmed parts are each a union of their own waits, so with overlapping waits they can add up to more than the total; each is at most the total.
- The orchestrator is recognised by its author name alone. An entry written under that name by anyone else counts as the orchestrator's.

## Alternatives considered

- **Let a recommendation end the wait, as any reply did.** The agent has no answer until the operator confirms, so the figure would understate the waiting, most of all where the orchestrator only recommends.
- **Count the time between a recommendation and its confirmation under the orchestrator.** The operator, not the orchestrator, ended that wait; counting it apart as confirmed shows what recommending saves without crediting the orchestrator with the operator's answer.
- **Split the total into exclusive parts, overlaps given to one ender.** Exact sums, but which ender wins an overlap is arbitrary, and the parts would change when an unrelated wait overlaps them.
- **A per-item `wait_threads_confirmed_seconds`.** Not asked for by the charts S-0220 serves; the week's `confirmed` carries it, and an item's field can be added later without changing these.
- **Recognise the orchestrator by the run that wrote the entry rather than its name.** Thread files record an author name only; the orchestrator's run is named `orchestrator`, as the board's placements already rely on (ADR-0087).
