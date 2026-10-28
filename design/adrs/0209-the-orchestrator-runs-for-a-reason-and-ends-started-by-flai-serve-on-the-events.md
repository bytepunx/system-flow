---
id: ADR-0209
title: "The orchestrator runs for a reason and ends, started by flai serve on the events it can act on, coalesced, with a pack primed for them, and flai serve holds every start while the host's budget is exhausted, circulating the exhaustion as an event"
status: accepted
date: 2026-10-28
supersedes: []
superseded_by: []
refines: [ADR-0087, ADR-0108, ADR-0154, ADR-0207]
---

# ADR-0209 The orchestrator runs for a reason and ends, started by flai serve on the events it can act on, coalesced, with a pack primed for them, and flai serve holds every start while the host's budget is exhausted, circulating the exhaustion as an event

## Context

[ADR-0087](0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md) runs the orchestrator as one long session per project: it primes, decides, logs, holds `wait_for_events`, and decides again, without ending. [ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) writes that as the role's trigger, `continuous`, the only role with one. The other three roles start for a reason, `asked`, an event, or a schedule, do their work, and end.

The orchestrator's session logs from 2026-10-14 to 2026-10-28 (`~/.flai/serve/agents/<key>-orchestrator-<start>.log`, Claude Code's stream, with each turn's usage) show what that loop costs and why.

**What a wake ingests is tiny; what it pays for is the transcript.** A `wait_for_events` result is about 560 bytes, some 300 tokens, and the orchestrator's reply is about 30 tokens, usually the next `wait_for_events`. But each turn sends the whole conversation: the role prompt, the system prompt, and the schemas of 26 tools (34k tokens on turn one); the pack from the last prime, three or four parts of some 35k characters each; and every tool result since the last compaction (`promote_candidates` about 27k bytes a call, `activity_log` and `item_get` about 11k, `thread_get` about 7k). The context grows about 1k tokens a turn, reaches 850k to 960k after 700 to 900 wakes, compacts to 12k, re-primes, and climbs again. Over the fortnight the average context read per turn stayed between 385k and 468k in every system-flow session. Nearly all of it is a cache hit, which keeps the dollar figure low per wake, and all of it counts against the account's usage limit.

**Most wakes do nothing.** In the three largest system-flow sessions, 75 to 80 percent of `wait_for_events` calls were followed by another `wait_for_events`. By cause, with the share of context tokens each cost:

| Wake carried | Oct 14 to 21 | Oct 23 | Oct 28 |
|---|---|---|---|
| A path change under `wip/` and no event (`changed` alone: another agent's narrative sync, `index.md`, a thread file) | 35 % | 49 % | 42 % |
| Only forecast edits by `flai`, the replan after an acceptance, one event per open story | 19 % | 25 % | 27 % |
| Only task moves by stories' agents | 13 % | 11 % | 9 % |
| The 30-minute timeout | 6 % | 0 % | 1 % |
| Events the orchestrator may act on | 27 % | 15 % | 21 % |

`wait_for_events` returns whenever the repository snapshot differs, and the catch-up then finds no event for the orchestrator in half the wakes. The replan after each acceptance ([ADR-0135](0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md)) edits `forecast` on every open story with `by: flai`, and each edit arrives as its own wake. Idle, the two timeouts an hour read about 1M tokens; busy, 400 to 500 turns an hour read about 200M.

**The cost scales with projects and with activity, not with decisions.** One session read 3.3B context tokens over Oct 14 to 21 ($503 reported). The next week read 3.6B in four active days, with five projects' orchestrators running by Oct 26 against one on Oct 14, and busy-hour wake rates roughly doubled as more stories' agents ran in parallel. Per action taken (a thread reply, a promotion, a plan asked, an acceptance), the Oct 28 session read about 7M tokens.

**The transcript holds little the orchestrator uses.** It compacts to 12k several times a day and re-primes, and the orchestrator carries on from flai's state: the inbox, the board, the threads, its activity log, its own thread ([ADR-0189](0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-sends-back-a-story-whose.md)). What the transcript retains between compactions is stale tool results, read again on every wake.

**The usage limit turns the loop into a storm.** The account's weekly limit was reached on Oct 14, Oct 20, and Oct 25; the five-hour session limit on Oct 28 after three busy hours. On the 429 the session ends with exit 1, and `flai serve` starts it again a minute later (ADR-0087 §1, `orchestrateRetry`), with no backoff and no knowledge of why it failed. Oct 26 and 27 produced 8,986 one-turn sessions across five projects, a start a minute per project for two days, each a request the API refused. Stories' agents fail the same way under [ADR-0108](0108-flai-serve-restarts-a-story-s-agent-that-ended-with-its-story-in-progress-up-to.md), up to `agent.auto_restarts` each, and nothing tells the other projects' serve loops, the planner, the curator, or the analyzer that the host has no budget.

The question is whether the orchestrator should be a loop at all, and what `flai serve` should do when the API says the account is out of budget.

## Decision

**The orchestrator runs for a reason and ends, like every other strategic role: `flai serve` starts it when an event it can act on happens, coalesced over a window, with the reason in its start prompt and a pack primed for that reason, and holds every strategic and story start while the host's budget is exhausted, circulating the exhaustion as an event instead of retrying each minute.** This refines ADR-0087, ADR-0108, ADR-0207, and [ADR-0154](0154-a-reply-may-settle-a-conversation-which-then-awaits-neither-story-and-flai.md); the rest of each stands.

1. **The orchestrator's trigger is a list of events, not `continuous`.** `orchestrator.md` names them under `trigger`, in ADR-0207's vocabulary: `asked`; `{event: thread, when: awaits-orchestrator}` (a thread the inbox marks awaiting it or the operator, [ADR-0181](0181-a-thread-awaits-an-agent-only-when-the-agent-takes-part-in-it-or-it-is-on-the.md)); `{event: moved, when: to-review}`; `{event: moved, when: to-backlog}` and `{event: edited, when: draft-changed}` for drafts to finalize; `{event: planned}` (a planner's run ended, its stories or tasks written); `{event: accepted, when: release-due}` (the release policy's figures changed); `{event: ready-empty}` (the ready column drained below the in-progress limit); `{event: reordered}`; and `{schedule: orchestration.schedule}`, unset by default, for a periodic look. `flai check` rejects `continuous` on any role. A file change with no event, a task move, and a replan's forecast edits start nothing.
2. **Starts are coalesced.** `flai serve` queues the reasons it sees and starts one run per project no sooner than `orchestration.coalesce` after the first queued reason (a Go duration, `30s` by default), with every reason seen meanwhile, as the curator's acceptances are coalesced ([ADR-0199](0199-the-curator-runs-after-acceptances-coalesced-and-on-digest-schedule-marks.md)). Reasons that arrive while a run is going start the next run when it ends. One run, one journal entry with its reasons.
3. **The start prompt names the reasons, and the pack is primed for them.** The prompt lists each reason with its item or thread. `flai prime --role orchestrate --reason <reason>...` (MCP `prime` with `reasons`) loads the role's conventions once and, per reason, what it needs: for a thread, the thread, its item, and the conventions on threads; for a review, the story, its last `verify` result, and the acceptance rules; for a plan, the planner's thread and the items it wrote; for a release, `release_evaluate`'s figures. The general pack of three to four parts is not loaded. The prompt's step 15 (ADR-0087 §4) becomes: act on each reason, log each action, then end; `flai serve` starts you again when there is more.
4. **The run ends when its reasons are done.** `wait_for_events` answers the orchestrator `end: true` with `why: reasons done` as soon as nothing queued remains for it, as it answers a story's agent that asked (ADR-0154), so that a prompt that holds it anyway still ends. A run that holds past `orchestration.run_timeout` (`20m` by default) is stopped and logged.
5. **Budget exhaustion is an event, and it pauses the host.** When a run of any role, or a story's agent, ends with the API's usage-limit refusal (a 429 whose body names the limit, Claude Code's `api_error: usage_limit_reached` with the reset time it reports), `flai serve` records a `budget` state on the host, in the serve folder beside `upgrade-asks.json`, with the provider, the limit named, the reset time when the API gives one, and else a backoff that doubles from `5m` to `1h`. While the state is set, no run of any role and no story's agent is started on any project the host serves with that provider; the reasons queue and the stories wait where they are, and none of this counts as a failed run against `agent.auto_restarts` or `orchestrateRetry`. At the reset time, or after the backoff, one probe run is started; its success clears the state, its refusal extends it. The state is a journal entry, an `inbox` field (`budget`, with the reset time) for every agent on the host, a dashboard banner, and a reminder through the channel of [ADR-0205](0205-flai-serve-reminds-the-operator-of-a-thread-that-awaits-them-past-threads.md) when one is set.
6. **Restarts are events too.** Each start, end, and failure of a role's run or a story's agent is an event of kind `agent` (`started`, `ended`, `failed`, with `why`), reported by `wait_for_events` and `inbox` as the others are, so that the orchestrator, the analyzer, and the operator see an agent restarting three times without reading the serve journal. The thread ADR-0108 opens after `auto_restarts` cites them.
7. **What stays.** The guard's hold on the orchestrator's calls, its permissions, its activity log and charge, its own thread, `flai serve orchestrate stop|start`, and the journal entries are as ADR-0087 and ADR-0189 have them. The `orchestrate` host action still gates it: off, nothing is queued.

## Consequences

- An idle project costs nothing: no timeout wakes, no reads of a 450k transcript to learn that a narrative synced.
- A decision costs its own context: the 34k base, a pack of one reason's size, the reads it makes, and a few turns, in the order of 0.3M to 0.6M tokens against about 7M today. The base is identical across starts, so the one-hour cache keeps it hot in a busy period.
- No compaction, no re-prime of a stale pack: each run reads the current pack.
- The orchestrator loses the transcript between runs. What it needs across runs lives in flai already, and §6 gives it the one thing it could notice only by watching: an agent that keeps restarting.
- A reason arrives one coalescing window late. A thread awaiting the orchestrator waits `orchestration.coalesce` longer than today's second; for a reply the operator reads hours later, nothing.
- The host stops at the first refusal instead of a start a minute per project, and every agent and the operator learn why and until when.
- `wait_for_events`'s `changed` field, the replan's events, and task moves still reach a story's agent that holds the call; nothing changes for it but the `agent` events and the `budget` field.
- `flai serve` gains a reason queue per project, the budget state, the probe, and `orchestration.coalesce`, `orchestration.run_timeout`, and `orchestration.schedule` in the manifest; `flai check` validates them. The template's `orchestrator.md`, the orchestration convention, `flai-cli.md`, and the serve docs change in the same story.
- The usage limit still binds the five projects together: pausing is per provider on the host, not per project, because the limit is the account's.

## Alternatives considered

- **Keep the loop and filter the wakes.** Return from `wait_for_events` only on an event the orchestrator can act on, and batch the replan's edits. It removes the empty wakes, about 60 to 75 percent of them, but every remaining wake still reads the whole transcript, the idle timeouts stay, and the compaction cycle and its re-prime stay. Half the saving for the same change to the wait.
- **Keep the loop and shrink the transcript.** Trim tool results, prime fewer parts, compact sooner. Each wake still costs the transcript; the lever is the number of wakes.
- **A fixed schedule alone.** Start the orchestrator every N minutes. Simple, but a thread awaiting it waits N minutes and an idle project still pays N times an hour. Kept as `orchestration.schedule` for projects that want a periodic look.
- **Back off the restart on a 429 and change nothing else.** Stops the storm, and nothing tells the other projects, the other roles, or the operator; each learns by failing in turn. §5 is this with the state shared and circulated.
- **One orchestrator for the host.** Fewer sessions, but one run would hold several projects' permissions and checkouts, against ADR-0087's one orchestrator per project; the per-provider budget state (§5) gets the sharing that is needed without it.
