---
id: ADR-0105
title: "A story's empty wakes, the wait_for_events calls of its agents that timed out with nothing to report, are counted from their logs into its usage, and flai stats reports them"
status: proposed
date: 2026-10-07
supersedes: []
superseded_by: []
---

# ADR-0105 A story's empty wakes, the wait_for_events calls of its agents that timed out with nothing to report, are counted from their logs into its usage, and flai stats reports them

## Context

A story's agent that has asked the designer a question holds the flai MCP tool `wait_for_events` until something changes. The tool answers when an event arrives or a path it watches changes, or when its timeout runs out, with `timed_out: true`, no events, and no changed paths. Each such answer wakes the agent for a turn that finds nothing to do and calls the tool again, and each turn is a model call over the agent's whole context.

E-0017's classification of the story runs logged on this host found 650 such turns across 38 story runs, 1,235 minutes of agent time. S-0272 removes most of them: an agent whose story has a question of its own open to the designer and no task in progress ends instead of waiting, `wait_for_events` answers it `end: true` at once, and flai serve starts it again when the question is answered. Its criterion 4 asks that `flai stats` count the empty wakes, so that the saving is measured rather than assumed.

Nothing records the turns today. `flai stats` reads work items, threads, issues, and the strategic agents' activity documents; it never reads an agent's log. What agents spent reaches it through each item's `usage` front matter, which flai serve measures from the stream-json logs it keeps of the agents it started for the story ([ADR-0051](0051-work-items-record-the-tokens-and-cost-their-agents-spent-measured-from-the.md)). The question is what counts as an empty wake, where the count is kept, and where `flai stats` reports it.

## Decision

**A story's empty wakes are counted from its agents' logs when flai serve measures the story's usage, and recorded in its `usage` as `empty_wakes`; `flai stats` reports them per item and, under `waiting`, over the window and per week.**

1. **An empty wake** is a call to the MCP tool `wait_for_events` of the flai server, by a story's agent itself, whose result reports `timed_out: true`, no events, and no changed paths. In Claude Code's stream-json log it is a `tool_use` named `mcp__flai__wait_for_events` in an event without `parent_tool_use_id`, paired by its `id` with the `tool_result` whose `tool_use_id` is that ID, not `is_error`, whose JSON has `timed_out` true, `events` empty, `changed` empty, and `end` not true.
   - A sub-agent's call is not counted: its events carry `parent_tool_use_id`.
   - A call refused or failed is not a wake: its result is an error, as flai guard's refusal is.
   - A call that answered `end: true` is not empty: it ends the agent, which is the saving.
   - A call with no result in the log, because the run ended during it, is not counted.
2. **Counted when the story is measured.** When flai serve measures a story from the logs of the agents it started for it, when an agent ends, when a task of a running agent enters done, or by `flai serve agent usage --write` (`serve.Measure`, from `usage.Read(...)`), it counts the empty wakes over every one of those logs, as it sums their tokens, and writes the count with the rest of the story's `usage`. A measurement replaces the count, as it replaces the tokens: it is never added to.
3. **Recorded in `usage` as `empty_wakes`**, a whole number, omitted when 0. Only a story measured from its logs (`source: log`) counts any. A task's usage carries none: a call is the story's agent's, and is not apportioned to tasks. An epic's usage, summed from its stories (`source: sum`), carries the sum of theirs; a story summed from its tasks has none. A strategic charge leaves it as it is. Stories measured before this change carry none, and the count starts with the flai release that ships it.
4. **`flai stats` reports it from the items**, not the logs, as it reports their tokens:
   - `items[].usage.empty_wakes`, the item's `usage.empty_wakes`, 0 when its usage has none;
   - `waiting.empty_wakes`, with `count`, the sum over the items of the report's type completed in the window, cancelled ones left out, as `waiting.weeks[]` takes them, and `mean`, that count over the number of those items that carry agents' usage, absent when none does;
   - `waiting.weeks[].empty_wakes`, the sum over the items completed in the week;
   - in the text report, one line in its waiting section.

The definition is the empty-wake class of S-0293, which classifies every turn of a run: S-0293 takes this count over, built on the same parsing, rather than defining the class again.

## Consequences

- The saving S-0272 promises is measured: once its release is installed, the empty wakes of the stories completed each week show in `waiting.weeks[].empty_wakes`, beside the thread waits that cause them.
- The count is of calls, not minutes or cost: an empty wake's turn is priced in the story's tokens already, and is not split out. S-0293 classifies turns with their cost.
- Weeks before the release, and stories measured before it, show none: they read as 0, not as unmeasured. A mean over a window that spans the release counts stories with usage but no count, so it understates the rate until the window is past it.
- Only agents flai serve starts with the `claude-code` harness are counted, as only they are measured. A story worked by hand counts none.
- `flai/internal/usage` pairs a `tool_use` with its `tool_result` while it reads a log, and `usage.Usage` gains `EmptyWakes`. `Add` and `Sum` add it; the roll-up keeps it.
- `empty_wakes` is listed under `item.usage` in `front-matter-fields.txt`. An older flai ignores a key inside `usage` that it does not know, and a write by it would drop the key, so publishing the release that carries it raises `flai.minimum` to it, as `strategic` did.
- `design/system/metrics.md` and `design/system/work-hierarchy.md` define the values.

## Alternatives considered

- **`flai stats` reads the host's agent logs itself.** The logs are the host's: flai serve prunes them, and another clone of the repository has none, so the same commit would report different figures on different machines. Usage front matter is how agents' figures already reach `flai stats`, and it travels with the item.
- **Leave it to S-0293, which classifies every turn.** Rejected by the operator on TH-0215: the saving is to be measured as soon as S-0272 lands. S-0293 builds on the same parsing.
- **Apportion the calls to tasks, as tokens are.** A wake is the story agent's call between tasks, while it waits for an answer; a share of a call is not a count of anything.
