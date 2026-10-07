---
title: Agent narrative
updated: 2026-10-07
status: active
topics: [all]
---

# Agent narrative

`wip/agents` exists so that a crash, a context reset, or a hand-off between agents costs minutes, not hours. It is the curated memory of each active work stream, written by the agent doing the work, for the next agent who picks it up.

## Files

```text
wip/agents/
├── README.md        # this convention in brief, shipped by the template
├── index.md         # table of active streams, maintained by flai and agents
├── S-0004.md         # one narrative per active story, named by story ID
├── planner.md        # the planner's activity document, written by flai
├── orchestrator.md   # the orchestrator's activity document, written by flai
└── analyzer.md       # the analyzer's activity document, written by flai
```

A work stream is a story. Tasks report into their story's narrative through the `stream` key. Epics have no narrative; their story narratives are enough.

When a story is archived, its narrative moves with it to `wip/archive/agents/`.

## Narrative structure

```markdown
---
stream: S-0004
title: CLI scaffold and config
updated: 2026-09-15T17:02:00Z
agent: claude-fable-5-1          # last agent to write, free text
session: 5da6af50                # opaque, helps correlate with tool logs
host: build-box                  # the host whose flai stream open last opened it
---

# S-0004 CLI scaffold and config

## Context
Two or three paragraphs a fresh agent needs before touching anything. What the story is, what is already true in the repo, what constraints apply. Rewritten as understanding improves.

## Current state
Where things stand right now. Which tasks are done, what is half done, what files are dirty. Rewritten on every meaningful step with `flai stream state`. This is the first thing a resuming agent reads.

## Next steps
Ordered list. The first item is what to do next. Rewritten on every meaningful step with `flai stream state`.

## Decisions
Bullet list of decisions made during this stream with a one-line rationale each. Anything architectural also gets an ADR; link it.

## Open questions
Questions for the human. Each has a date. Answered questions move to Decisions: by hand, or with `flai stream answer` (also reachable from the dashboard's inbox, S-0090), which removes the matching bullet and records the answer under Decisions itself. A story whose narrative still has a hand-written question here is refused into review (S-0089): the agent cannot finish without the answer, so `flai move` (and the MCP and dashboard paths behind it) names the question and how to answer it; answering it, or removing it, lets the move through. Mirrored threads do not count toward this rule, since they are tracked, and close, as threads. Between `<!-- threads:start -->` and `<!-- threads:end -->` flai keeps a generated list of unresolved threads anchored to the story or its tasks (ADR-0020, S-0038); do not edit that block by hand. Agents connected to `flai mcp` see the same threads through `inbox` and are woken by `wait_for_events` (S-0039).

## Log
Append-only. One entry per meaningful step, newest last.

### 2026-09-15T16:12:00Z
Started. Pulled S-0004 to in-progress. Read design/system/flai-cli.md.

### 2026-09-15T16:40:00Z
T-0021 done. Config read/write with tests. Decided on plain encoding/json over viper, see Decisions.
```

An agent writes `## Current state` and `## Next steps` with `flai stream state S-nnnn --current "<text>" --next "<text>"`, the MCP tool `stream_state`, or the host channel's `stream.state` (S-0271), never by editing the file. It replaces the two sections, or the one it is given, leaves every other section as it was, and appends nothing to the log. It stamps `updated`, `agent`, and `session` as `flai stream log` does, and writes `index.md` again. The same text again writes nothing. It refuses a story that is not in progress or in review, a story with no narrative, and text the project's markdown lint rejects; a text holding a `#` or `##` heading is an error. `flai check` warns, under the rule `narrative.state`, when the narrative of a story in progress or in review has an empty `## Current state` or `## Next steps`, or the template's placeholder in one, or a `## Next steps` that is not a list. The warning is advisory, so `--strict` passes over it. The other sections are written as before: `## Context`, `## Decisions`, and `## Open questions` by hand, and the `## Log` by `flai stream log`.

## Strategic agents' activity documents

The planner, the orchestrator, and the analyzer work above a story and keep no narrative. Each has one activity document per project instead: `planner.md`, `orchestrator.md`, and `analyzer.md` ([ADR-0079](../adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)). flai writes them; no agent edits them by hand. A document exists once its agent has logged its first activity.

```markdown
---
kind: planner
accrued_cost: 0.4213
accrued_seconds: 723
tasks_completed: 1
last_run: 2026-10-03T18:00:00Z
---

# Planner activity

flai writes this document, one entry per activity, newest last; do not edit it by hand.

## Log

### 2026-10-03T18:00:00Z

- Summary: Drafted five stories for E-0016.
- Trigger: asked
- Items: E-0016, S-0230
- Seconds: 723
- Cost: 0.4213 USD, estimated
```

| Front matter | Meaning |
|--------------|---------|
| `kind` | `planner`, `orchestrator`, or `analyzer`; the file's name |
| `accrued_cost` | Sum of the entries' cost, in US dollars |
| `accrued_seconds` | Sum of the entries' wall-clock seconds |
| `tasks_completed` | Number of activities logged |
| `last_run` | When the newest activity ended |

The front matter holds these five keys and no others; it is parsed strictly, so an unknown key is an error.

Each log entry is one activity. Its heading is when the activity ended; a second entry ending in the same second gets `(2)` after the time, so no two headings are the same. The summary is one line the agent gives. `Trigger` is one line saying what started the run, written for a planner run `flai serve` started: `asked` when the operator asked for it, otherwise the replanner's triggers, separated by semicolons ([ADR-0084](../adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)). It is absent from older entries and from those an agent logs with `activity_log`. `Items` lists the items the activity touched, or `none`. `Cost` has four decimals and is marked `estimated` when it was apportioned or priced rather than reported. Appending an entry adds its cost and seconds to the totals, counts it, and moves `last_run` to its end when that is later.

An activity ends in one of two ways. An agent whose run spans activities, as the orchestrator's does, reports each with the MCP tool `activity_log`, giving its kind, summary, and items; flai measures it and writes the entry. When a strategic run ends, `flai serve` logs the time since the last entry as one activity, with the run's last result as its summary, unless nothing was spent in it. Either way the activity's seconds and cost come from the run's stream-json log, apportioned to the activity's span as a task's are to its intervals ([ADR-0051](../adrs/0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)).

`index.md` lists the documents that exist under `## Strategic agents`, after the active streams, with each agent's activities, cost, seconds, and last run, or that it is unreadable when it does not parse. `flai check` validates their front matter and log and does not treat them as narratives. `flai stats --json` reports each kind's totals and its log entries in the window as `strategic`, defined in [metrics.md](metrics.md).

## What an agent is told through MCP

`flai mcp` is the agent's view of the repository ([ADR-0020](../adrs/0020-files-plus-mcp.md)); the files stay the record. An agent on the host starts it over stdio from `.mcp.json`; one that cannot start a process there reaches the same server over HTTP after `flai mcp start` ([ADR-0030](../adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)), named by its `X-Flai-Agent` header as `FLAI_AGENT` names it locally, so the cursor and everything below are the same on both. Since S-0058 the view covers work as well as threads, and it serves an agent that ends its turn between the designer's messages as well as one that stays running.

- `inbox` reports three things. Threads awaiting this agent. The stories that are ready to pull, in pull order, with whether a pull is allowed (`can_pull`: false while the in-progress limit is full or review is at or over its limit, with which in `pull_hold`, S-0243, [ADR-0073](../adrs/0073-a-full-review-holds-the-pull-and-flai-check-strict-passes-over-review-over-its.md)): this is state, listed on every call, so an agent that has forgotten everything still sees the work. And the changes others made since this agent last looked: an item created, moved, blocked, or unblocked, with what happened, to what, by whom, and when, and a note when the pull order changed.
- Changes are derived from the files, not recorded anew: `transitions` carry `to`, `at`, and `by`; `blocked` intervals carry `from`, `until`, and `reason`; `board.md` carries `order`. A change whose `by` is this agent is not reported back to it, which is why `flai move`, `block`, and `unblock` record `FLAI_AGENT` when it is set.
- "Since this agent last looked" is a cursor per agent name under `.flai-cache/mcp/`, outside git. It is a read marker, a cache in the sense of the overview's principle that tooling never owns state: losing it repeats or skips the report of a change and loses nothing else, because ready work and open threads are state and are always listed. With no cursor, the last 24 hours are reported.
- Edits are the one change that is not derived from the items' files (S-0085). When someone changes an item's title, fields, or body with `flai edit` or from the dashboard, flai notes who and what in `.flai-cache/edits.jsonl`, and `inbox` and `wait_for_events` report the entries by others as `edited` changes whose `to` names what changed. An agent told so reads its story again before going on. The note is a log beside the cursors, not a key in the front matter, because items are parsed strictly and an older flai reading the same repository would refuse an item with a key it does not know. Losing the log loses only the telling; a hand edit of a file is, as before, not reported.
- Overlaps at acceptance are the other (S-0132, [ADR-0046](../adrs/0046-a-ready-story-whose-claim-overlaps-an-open-story-s-is-held-yellow-and-with-its.md)). `flai accept` notes, in `.flai-cache/overlaps.jsonl`, each open story whose claim covers a path the accepted story's merge changed. `inbox` and `wait_for_events` report each note once as an `overlapped` change on the open story, with `cause` the accepted story and `to` the paths, and with a summary that names at most ten of them. They report it whoever accepted, since the news is for the open story's agent. The log is kept apart from `edits.jsonl` because an older flai reads every line there as an edit. An agent told so syncs its story and runs its tests again before going on.
- A claim that grows into another's is told the same way (S-0244, I-0059). When `flai task new`, `flai edit --touches`, `flai touches`, `item_new`, or `item_edit` adds paths to the claim of a story in progress that another story in progress claims, flai notes it in `overlaps.jsonl` for both stories, and each gets one `overlapped` change with `cause` the other story, `to` the paths, and a summary saying the claims grew to overlap and to agree with the other story's agent in the two stories' conversation (`message_reply`) who changes them first. Beside the change, the story whose claim grew messages the other, by the writer, `about` the paths gained, in the pair's open conversation or a new one, asking which of the two changes them first (S-0332, [ADR-0121](../adrs/0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md)). An agent told so agrees with that agent on the conversation who changes the paths first, and narrows its touches if it can, before it changes them. The note leaves `accepted` empty, so an older flai passes over it.
- Agents of two open stories agree through messages, kept apart from the operator's threads (S-0331, [ADR-0120](../adrs/0120-agents-of-two-open-stories-message-each-other-in-conversations-kept-under-wip.md)). `message_send` (`to`, `text`, `about`) starts a conversation, `message_reply` (`id`, `text`) answers one, and `message_get` (`id`) reads one; each answers the conversation as `flai message show --json` prints it. Send and reply write as the calling session's own story: `FLAI_STORY`, else the story in an agent name of the form `agent-S-nnnn`, else the story branch checked out. No argument names another story, and a session with none is refused with nothing written; `message_get` needs no story. `inbox` lists the open conversations of the agent's story under `messages`, each with `id`, `title`, `with` (the other story), `about`, `updated`, `entries`, `last` (the newest entry), and `awaiting` (`you` or `other`); it is absent with no story or no open conversation, and `awaiting_you` counts threads only. `wait_for_events` watches `wip/messages` and reports a message from the other story once, as a change of kind `message` whose `cause` is the sender's story, `by` its agent, `to` the paths it is about, and `id` and `title` the conversation's; the agent's own entries are not reported back. An agent answers a message awaiting its story with `message_reply` before it goes on with those paths. When the two do not agree it escalates with `message_escalate` (`id`, `reason`), as `flai message escalate` does, rather than `thread_open` (S-0332, [ADR-0121](../adrs/0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md)): it opens a thread on its own story for the operator that names both stories, gives the conversation's path, and quotes the reason, and adds a message naming the thread to the conversation, which awaits the other story and stays open. It answers `thread` and `conversation`. `flai guard` refuses a sub-agent `message_send`, `message_reply`, and `message_escalate`.
- A conflict the trial merge at `flai stream sync` finds with another open story's branch arrives as a message too (S-0332, [ADR-0121](../adrs/0121-flai-tells-two-stories-of-a-trial-merge-conflict-or-a-grown-overlap-with-a.md)). The message comes from the story that synced, by `flai`, `about` the conflicting paths, in the pair's open conversation or a new one; it names both branches and the paths and asks the two to agree who changes what. The same paths again add nothing. A sync that finds the pair merging cleanly, or the other story no longer open, closes the conversation when flai's newest message in it tells of the conflict. The agents settle it on the conversation and escalate only when they do not agree; no sync opens a conflict thread.
- A task's change reaches the other open stories as a message (S-0333). When `flai task done` or `task_done` closes a task of another story with a commit that changes paths this story's claim covers, the shared paths included, it sends this story a message from that story, on their open conversation or a new one, `about` those paths. Its first line says which task of which story changed paths this story's claim covers. The rest names the task's title, the commit's short hash, the story branch, the commit's subject, and the paths, says the change reaches the main branch when that story is accepted and `git show <hash>` shows it until then, and asks for a reply if it breaks the work. `inbox` lists it under `messages`, and `wait_for_events` reports it as a change of kind `message`. The agent checks the change against its work, replies, and adjusts early where it can; a sync does not bring the change until that story is accepted, when an `overlapped` change follows ([workflow.md](workflow.md#branches-and-collisions-adr-0019)).
- `wait_for_events` returns at once when the cursor is already behind, so a change made between two calls is not lost; otherwise it blocks until something changes or `timeout_seconds` passes, 60 s by default. It returns events in the same shape as `inbox` does, with the changed paths, and advances the cursor. It reports work items, threads, narratives, and messages to the agent's story, not sub-agents: the story's agent waits for a sub-agent it started in the background through the harness's own notice that it ended.
- A story's agent does not hold `wait_for_events` for the designer's answer (S-0272). When nothing is left but the answer it writes the narrative's `## Current state` and `## Next steps` with `stream_state`, naming the question and what each answer leads to, and ends, because flai serve starts it again on the answer and the restarted agent reads them, beside its first `inbox`, which holds the answer. When flai serve started it, a thread on its story awaits an answer to its own question, and no task of the story is in progress, `wait_for_events` answers at once with `end: true` and a `why` naming the threads, and the agent writes the same two sections and ends ([workflow.md](workflow.md#branches-and-collisions-adr-0019)). Since S-0335 it ends the same way on another story's agent's reply: `wait_for_events` answers `end: true` also when a conversation of its story awaits the other story's reply and none awaits its own, the `why` naming each conversation and the story it waits on, and flai serve starts it again when that agent replies or a new message to its story comes, which its first `inbox` lists under `messages`. An agent run by hand, and the planner, the orchestrator, and the analyzer, which flai serve does not start again on an answer or a reply, still hold it.
- `wait_for_events` and `wait_for_work` hold a call for the `timeout_seconds` asked, up to 30 minutes (S-0244, I-0059; five minutes until then). Every return is a model turn that rereads the agent's whole context, about 0.05 USD at 150k cached tokens, so a five-minute cap made an idle agent pay that every five minutes. While held, a wait sends a progress notification every minute to a call that carries a progress token, because Claude Code drops a tool call that sends nothing for 30 minutes over stdio and 5 over HTTP (`CLAUDE_CODE_MCP_TOOL_IDLE_TIMEOUT`), and a notification starts that count again.
- `board` returns what `flai board --json` prints: columns, limits, pull order, breaches.

An agent that ends its turn calls `inbox` at the start of every turn. An agent that stays running, and has no story of its own in progress, holds `wait_for_work` (S-0097) and does what it answers: go back to its own story in progress (`resume`; a story is its own when its narrative was last written under its name), answer threads awaiting it written to since `wait_for_work` last answered (`thread`), or pull the first ready story in pull order once the in-progress limit leaves room and review is under its limit (`pull`). Otherwise the call waits, re-deciding whenever a thread, item, or narrative file changes, for 5 minutes unless `timeout_seconds` says otherwise, and a timeout says what it is waiting for (`room`, `review`, `held`, or `ready`), after which the agent calls it again. Either way a ready story is pulled without waiting to be told. Two agents told to pull the same story race on `story_start` (S-0274), which refuses the second ("in progress already"), and that one waits again. A thread awaiting the agent that was not written to since the last answer does not wake it, so a thread it cannot answer does not make it spin; `inbox` still lists every one.

## Obligations

An agent working in a conforming repo must:

1. Open the narrative before making the first change for a story.
2. Rewrite `## Current state` and `## Next steps` with `flai stream state` (or `stream_state`) after every task transition and before any long-running operation.
3. Append a `## Log` entry at every task transition, every decision, and every blocker.
4. Never store secrets, full transcripts, or raw tool output in the narrative. Summaries only, with paths to artefacts.
5. Update `index.md` when opening or closing a stream.

The baseline `CLAUDE.md` in the template states these obligations so every agent session in a conforming repo is instructed the same way.

## Recovery procedure

A fresh agent resuming after a crash:

1. Read `wip/agents/index.md`.
2. For each stream marked active, read `## Current state` and `## Next steps`.
3. Check `git status` for uncommitted work and reconcile with the narrative.
4. Append a log entry stating that recovery happened and what was reconciled.
5. Continue from the first next step.

The measure of a good narrative is that step 5 needs no human input.
