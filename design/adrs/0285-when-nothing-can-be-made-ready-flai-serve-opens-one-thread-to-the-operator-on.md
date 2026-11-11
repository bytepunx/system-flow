---
id: ADR-0285
title: "When nothing can be made ready, flai serve opens one thread to the operator on the board that lists each backlog story's decision, adds an entry only when the list changes, and resolves it once a story can start"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0205, ADR-0211, ADR-0279]
topics: [cli, orchestration, ready, threads]
---

# ADR-0285 When nothing can be made ready, flai serve opens one thread to the operator on the board that lists each backlog story's decision, adds an entry only when the list changes, and resolves it once a story can start

## Context

[I-0288](../issues/I-0288-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md), found by the analysis of 2026-11-07, records that from 2026-11-03 to 2026-11-05 ready, in progress, and review all stood at 0. Every backlog story was a draft, or was held by `after` behind a draft. The drafts S-0406 and S-0412 waited on TH-0532 and TH-0533. Each of those threads awaited the operator with one entry, from 2026-10-20, and both were resolved on 2026-11-07, once the analysis had named them.

Its instance shows three things.

- `ready-empty` fired as designed (`flai/internal/serve/readyempty.go`, [ADR-0211](0211-the-orchestrator-s-trigger-vocabulary-gains-the-events-moved-thread-planned.md) §1, [ADR-0279](0279-ready-empty-counts-only-the-ready-stories-not-held-on-after-and-a-move-of-a.md)), and the orchestrator was started for it. Its `ready` reason pack holds the board and `flai promote --candidates`. `flai promote --candidates` lists a draft only as refused, with the reason `draft: finalize it first`, and names no thread. The orchestrator's run of 2026-11-02 promoted nothing and told the operator nothing new (`wip/agents/orchestrator.md`).
- [ADR-0205](0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md)'s reminder, with `threads.remind_after: 4h` set since S-0498, would have reminded the operator of both threads. But it reaches them only through a channel set with `flai serve notify set`. With none set it only logs, and TH-0660 records that none was set.
- So neither the repository nor the dashboard told the operator that the board had stopped, or which of their decisions would start it again. The two threads were on the drafts' own pages, among the other threads that awaited the operator.

## Decision

**When no story in ready can start, the in-progress limit has room, and `flai promote --candidates` lists no candidate, `flai serve` opens one thread to the operator on `wip/kanban/board.md`. The thread lists each backlog story that cannot be promoted and the decision it waits on. flai serve adds an entry only when that list changes, and resolves the thread once a story can start.** This refines ADR-0205, ADR-0211, and ADR-0279. The rest of each stands.

1. **When.** At its look, after it notes `ready-empty`, flai serve finds the board starved when all three hold:
   - no ready story would start, counted as ADR-0279 counts them: ready stories, blocked ones left out, not held on `after`;
   - the in-progress limit leaves room, the limit less the stories in progress being above zero;
   - `workitem.PromotionCandidates` lists no candidate.

   With no in-progress limit, the board is never starved. A ready story held on overlap or on no touches counts as one that can start, as it does for `ready-empty`.
2. **Where and by whom.** One thread, anchored on the board file, `wip/kanban/board.md`, by flai serve's author, `flai serve`. While one such thread is open, no second one is opened. A restart of flai serve finds the open thread by its anchor, author, and title, as the merge queue finds its preflight thread ([ADR-0234](0234-before-it-builds-a-batch-the-merge-queue-runs-the-landing-s-preflight-on-the.md)).
3. **What it lists.** One line per backlog story that `PromotionCandidates` refuses, in the pull order. Each line names the story and the decision it waits on, the first that applies:
   - **A draft.** Finalize it, with the ID and title of each open thread on it or its tasks that awaits the operator ([ADR-0181](0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md)).
   - **Held by `after`.** The story at the root of its chain: each story named in `after` that is not done, followed until one names none that is not done.
   - **Blocked.** Its block reason.
   - **Anything else,** such as no forecast or no cost of delay. The reasons as `PromotionCandidates` gives them, which a planner run or the operator settles.

   When the host has no reminder channel, the text says so, and names `flai serve notify set`.
4. **When it gains an entry.** Only when the list's text differs from what flai serve last wrote on the thread. What it last wrote is kept beside the `ready-empty` flag in `.flai-cache`, so a restart writes nothing again.
5. **When it is resolved.** At the first look that finds the board not starved, flai serve resolves the thread, saying why: a ready story can start, or a candidate can be promoted, naming it. The next starving opens a new thread.
6. **Who answers it.** The thread is flai serve's, so the orchestrator leaves it alone: `answer_threads` passes over threads opened by neither a story's agent nor a planner. While it is open it awaits the operator like any other thread. So the dashboard's inbox lists it, and ADR-0205's reminder sends it through the channel, when one is set.

## Consequences

- An empty board is visible where the operator looks for what awaits them, on the board's own thread. The thread names each decision that would start work again.
- The thread names the operator's open threads on the drafts, so the decisions buried on item pages come up in one place.
- Nothing changes while any story can start or be promoted.
- A host with no reminder channel is told so on the thread.
- The orchestrator's `ready` reason still runs on `ready-empty`. The board's thread does not replace it.

## Alternatives considered

- **The orchestrator asks on its `ready-empty` reason.** It runs only while the `orchestrate` host action is on, and spends a model run on what flai can work out itself. Its answer would be a thread of its own, with nothing to keep it from asking again at each draining.
- **A reminder without a thread.** ADR-0205 already reminds of each thread, and that failed here for want of a channel. A thread is seen on the dashboard with no channel at all.
- **`flai promote --candidates` names the threads.** That helps the orchestrator's pack, but it still tells nobody. The list this ADR writes is the same join, so the command may gain it later.
