---
id: ADR-0206
title: "flai serve confirms an orchestrator recommendation the operator left unanswered past orchestration.permissions.confirm_after, on a thread a story's agent or a planner opened, never one marked for the operator alone, and flai stats counts those confirmations apart"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0090, ADR-0091, ADR-0205]
topics: [cli, orchestration, analysis]
---

# ADR-0206 flai serve confirms an orchestrator recommendation the operator left unanswered past orchestration.permissions.confirm_after, on a thread a story's agent or a planner opened, never one marked for the operator alone, and flai stats counts those confirmations apart

## Context

[I-0195](../issues/I-0195-a-thread-with-a-pending-orchestrator-recommendation-waits-on-the-operator-however-long-it-takes-with-no-setting-to-confirm-it-after-a-threshold.md) records TH-0583 and TH-0588: planners' design-choice questions on which the orchestrator recommended within minutes and the story then idled for a day or more until the operator confirmed. [ADR-0090](0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md) makes a recommendation the answer only on the operator's confirmation, and [ADR-0091](0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md) keeps the wait open until then. [ADR-0205](0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md) reminds the operator of such a thread but confirms nothing. On TH-0598 the operator chose a setting after which `flai serve` confirms the recommendation itself, off by default, and never for money, host changes, permission prompts, `.claude/` writes, or a thread the operator wrote on since.

## Decision

**With `orchestration.permissions.confirm_after` set, `flai serve` confirms, in the operator's absence, a pending recommendation of the orchestrator's that has waited that long, on an open thread a story's agent or a planner opened, with no entry by the operator since, and not marked for the operator alone; the confirmation is an entry of its own author, logged in the orchestrator's activity document and counted apart by `flai stats`.** This refines ADR-0090, ADR-0091, and ADR-0205; the rest of each stands.

1. **The setting.** `orchestration.permissions.confirm_after` is a Go duration, unset by default: unset or `0`, nothing is confirmed. It is the one permission that is a duration. `flai check` reports a value that is not a duration, or a negative one, under the manifest's orchestration finding. It is set with `flai manifest set` and shown in the settings catalog beside the other permissions.
2. **When.** On the per-minute look `flai serve` runs for each served project, the same look that reminds (ADR-0205), and on a change to its threads, only while the project's `orchestrate` host action is on: the permission is the orchestrator's, and with the orchestrator off no recommendation is posted and none is confirmed.
3. **What is confirmed.** A thread qualifies when all of these hold:
   - it is not resolved, and its pending recommendation (ADR-0090) is the orchestrator's;
   - the recommendation was written at least `confirm_after` ago;
   - its opener is a story's agent (`agent-S-nnnn`) or a planner (`planner-…`): never a thread the operator, `flai serve`, `flai accept`, the merge queue, or another strategic agent opened;
   - it is not a permission prompt's thread (`Allow <tool> <path>?`, which covers every write under `.claude/`);
   - no entry by the operator, the manifest's owner or the story's owner, follows the recommendation;
   - the recommendation is not marked for the operator alone.
4. **The operator-alone mark.** flai cannot tell money or a host change from a recommendation's words; the orchestrator can, and already classes them under `answer_threads`. It marks such a recommendation with `thread_reply`'s `operator_only`, or `flai thread reply --recommend --operator-only`, which writes the paragraph `For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.` above the `Source:` line. The mark needs a recommendation. A flai older than this reads the entry as a plain recommendation.
5. **The confirmation.** An entry by `flai serve`, written as `flai thread confirm` writes the operator's: `Confirmed the recommendation of <time> <author> after <confirm_after> with no answer from the operator (orchestration.permissions.confirm_after).`, citing the recommendation's source, with the thread's status `answered`. The asking agent is started again as on the operator's confirmation. The operator's answer becomes a digest entry; this confirmation does not, since the operator gave no answer. Each confirmation is logged in `wip/agents/orchestrator.md` as an activity of no cost, naming the thread and its item, and in `flai serve`'s log.
6. **The metrics.** `metrics.md` § Waiting gains `waiting.weeks[].threads.timeout`, a count and seconds as `orchestrator` and `confirmed` are, for the waits a confirmation by `flai serve` ended; `confirmed` keeps the operator's alone. In the longest waits, such a wait reads `awaited` `flai serve`.

## Consequences

- A planner's or a story's agent's question that the orchestrator recommended on no longer holds its story past `confirm_after`; the operator can still overrule it afterwards by answering on the thread.
- What only the operator may decide stays theirs: the threads `flai serve` and the permission prompt open by their opener and title, money and host changes by the orchestrator's mark.
- `flai stats` tells the operator's confirmations from the timeouts, so the share of recommendations confirmed unread is visible.
- A project that leaves the setting unset sees nothing change.

## Alternatives considered

- **The orchestrator confirming its own recommendation.** Rejected: ADR-0090 refuses a recommendation's author its confirmation, and the guard refuses every agent `flai thread confirm`.
- **Excluding every planner's thread.** Rejected: it would have excluded both instances, TH-0583 and TH-0588.
- **Telling money and host changes from the words.** Rejected as brittle: the orchestrator already makes that call.
- **A heading mark such as `(recommendation, operator only)`.** Rejected: a flai older than this would read the author as the whole heading and lose the recommendation.
