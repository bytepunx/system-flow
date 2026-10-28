---
id: ADR-0211
title: "The orchestrator's trigger vocabulary gains the events moved, thread, planned, ready-empty, and release-due and the conditions to-review, to-backlog, draft-changed, and awaits-orchestrator; flai serve and flai accept append the three new events to .flai-cache/events.jsonl and the inbox derives thread from the thread files"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0207, ADR-0209]
---

# ADR-0211 The orchestrator's trigger vocabulary gains the events moved, thread, planned, ready-empty, and release-due and the conditions to-review, to-backlog, draft-changed, and awaits-orchestrator; flai serve and flai accept append the three new events to .flai-cache/events.jsonl and the inbox derives thread from the thread files

## Context

[ADR-0209](0209-the-orchestrator-runs-for-a-reason-and-ends-started-by-flai-serve-on-the-events.md) §1 names the orchestrator's reasons in [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) §3's trigger vocabulary. That vocabulary has seven events, `accepted`, `cancelled`, `ready`, `review`, `edited`, `reordered`, and `thread-opened`, and three trigger conditions, `stale-figures`, `replan-agent`, and `each-ready-story`. ADR-0209 §1 writes events and conditions it lacks: `moved` with `to-review` and `to-backlog`, `edited` with `draft-changed`, `thread` with `awaits-orchestrator`, `planned`, `ready-empty`, and `accepted` with `release-due`. `flai check` refuses each as `roles.file`, so `orchestrator.md` cannot name them.

Three of them are not events flai has. `inbox.CatchUp`, which `wait_for_events` and `inbox` answer from, derives events from the work item files (`workitem.Changes`), from two logs under `.flai-cache` (`edits.jsonl`, `overlaps.jsonl`), and from the messages. A planner's run ending, the ready column draining, and an acceptance changing the release answer happen in `flai serve` and `flai accept` and leave no trace in any of them. A thread coming to await someone is in the thread files, which threads' many writers (the MCP server, the CLI, the dashboard, `flai serve`) change, but no event reports it.

## Decision

**The trigger vocabulary gains five events and four trigger conditions, each implemented by flai (ADR-0207 §3, §5); `flai serve` and `flai accept` append the three events that happen in them to `.flai-cache/events.jsonl`, and the inbox derives `thread` from the thread files.** This refines ADR-0207 §3 and ADR-0209 §1; the rest of each stands.

1. **Events.** `roles.Events` gains these, and `review` and `thread-opened` stay, so that no role file written to ADR-0207 breaks.

   | Event | Happens when | Fields | Kept in |
   |-------|--------------|--------|---------|
   | `moved` | An item moves to a state | the item, `to` the state, `by` | the item's transitions (`workitem.Changes`) |
   | `thread` | A thread's newest entry makes it await the reader (`you`), or, for the orchestrator, the operator, under [ADR-0181](0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md) | the thread, its story or path, `to` whom it awaits (`you` or `operator`), `by` the entry's author | the thread files, read at catch-up |
   | `planned` | A planner's run ends having written items beyond the item it planned | the item planned, `to` the IDs written, `by` the planner | `events.jsonl` |
   | `ready-empty` | The number of ready stories, held ones counted, falls below the in-progress limit's room (the limit less the stories in progress), once per draining: not again until ready reaches the room again. No in-progress limit, no event. | `to` the ready count, the room, and the limit, `by` `flai serve` | `events.jsonl` |
   | `release-due` | An acceptance changes what `flai release --evaluate` answers: whether the policy is met, or the figures it reads (the pending stories' count and value, a theme's stories not accepted); what changes with the clock alone does not count. A failed evaluation emits nothing and fails no acceptance. | the story accepted, both answers, `by` the acceptor | `events.jsonl` |

2. **`release-due` is an event of its own,** where ADR-0209 §1 wrote `{event: accepted, when: release-due}`. Whether the answer changed is known only at the acceptance, by evaluating before and after it; an `accepted` event read later cannot tell. A role names it as `{event: release-due}`.

3. **Trigger conditions.** Each is a match flai implements against an event that happened, in `roles`' own terms, so that `roles` imports nothing of the inbox.

   | Condition | Qualifies | Holds when |
   |-----------|-----------|------------|
   | `to-review` | `{event: moved}` | the item moved to `review` |
   | `to-backlog` | `{event: moved}`, and the calls `item_move` and `flai move` as before | the item moved to `backlog` |
   | `draft-changed` | `{event: edited}` | the item edited is a story that is a draft after the edit |
   | `awaits-orchestrator` | `{event: thread}` | the thread awaits the orchestrator or the operator |

   `to-backlog` is both a call condition and a trigger condition; on a call it holds the move to backlog, on a trigger it matches one.

4. **The log.** `.flai-cache/events.jsonl` in the main checkout, beside `edits.jsonl` and `overlaps.jsonl`: one JSON object a line, appended, never rewritten, outside git, as the cursors are. Each line has `at`, `kind`, `id`, `type`, `title`, `to`, `by`, and the kind's fields. `inbox.CatchUp` reads the lines after each reader's cursor once, under the 50-event cap with the rest, and `wait_for_events` wakes when the log changes. A lost log loses news, never state: the board, the threads, and `flai release --evaluate` still answer.

5. **`thread` has no log.** Threads are written in too many places for one writer's log to see each entry; the catch-up works out, from each open thread's newest entry after the cursor, whom it awaits, with the inbox's own rule, and reports it once per such entry.

## Consequences

- `orchestrator.md` can name every reason ADR-0209 §1 lists, with `release-due` as an event, and `flai check` accepts it; S-0536 writes it.
- A story's agent sees the new kinds in `wait_for_events` and `inbox` too, and may let `planned`, `ready-empty`, and `release-due` pass; `thread` is a thread now awaiting it.
- Each acceptance evaluates the release twice, which reads git; an acceptance already does more.
- `flai serve`'s reason queue reads the conditions' matches; this decision only implements them.

## Alternatives considered

- **Keep `{event: accepted, when: release-due}`.** The condition could only be judged at the acceptance, so it would still need a record of its own; an event says it plainly.
- **A log written by each thread writer.** It misses the writers that do not call it, the dashboard's and an older flai's.
- **Replace `review` and `thread-opened` with `moved` and `thread`.** It breaks role files written to ADR-0207 for no gain; both stay.
- **`ready-empty` when the ready column has no story.** It misses a column that still holds one story while three agents could start; the room is what the orchestrator fills.
