---
id: ADR-0051
title: "Work items record the tokens and cost their agents spent, measured from the agents' logs and summed up the hierarchy when an item is done"
status: accepted
date: 2026-09-29
supersedes: []
superseded_by: []
topics: [cli, dashboard]
---

# ADR-0051 Work items record the tokens and cost their agents spent, measured from the agents' logs and summed up the hierarchy when an item is done

## Context

The designer asked (S-0143, E-0011) that epics, stories, and tasks record the tokens and the cost their work took, and that charts show token rates per hour, cost, and completion against time and cost, each grouped by the model that did the work.

flai starts agents for ready stories (ADR-0043) and keeps each run's output in a log under its serve directory. For the `claude-code` harness that output is stream-json, and it says what was spent. Read from the 56 logs on the development host on 2026-09-29:

- Each run ends with a `result` event whose `modelUsage` gives, per model, input, output, cache read, and cache write tokens and `costUSD`, subagents included. A session resumed after an answer reports its totals again in its next run, cumulative: S-0122's second run reports 20402200 cache reads, its first 16935243, and the second run's own calls 3466957.
- Each call to a model appears as `assistant` events, one per content block, with a timestamp, the message id, the model, and its usage. Deduplicated by id, their input and cache tokens equal `modelUsage` exactly; their output tokens are the first streamed chunk's and fall far short.
- The reported cost per token of one model is steady enough to estimate with: 0.33 to 0.59 US dollars per million tokens for claude-opus-5-5 in 47 of 50 sessions.

Nothing else records tokens: the work items hold only transitions, and an agent a person runs by hand leaves nothing flai can read.

## Decision

Epics, stories, and tasks may carry `usage` in their front matter: `source` (`log` or `sum`), the `seconds` of agent work, `estimated` when any cost in it was not reported, and `models`, one entry per model with `input`, `output`, `cache_read`, and `cache_write` tokens and `cost` in US dollars. No command takes it as an argument; flai writes it.

1. **Measured from the logs.** `flai serve` measures a story from every log it keeps for the story's agents: each session's newest `result` is its total, and calls after it (a run still going, or one that died) are counted from their `assistant` events and priced at the model's blended rate, reported cost over reported tokens, across the logs it keeps, and marked estimated. A task is measured over the intervals it was in progress: each session's totals in the share of its calls' input and cache tokens that fall in them, marked estimated. The seconds are each run's first to last event, cut to the task's intervals. `flai serve` writes this when an agent it started ends, and when a task of a story whose agent runs enters done; `flai serve agent usage --write` does it by hand, for stories worked before.
2. **Summed up the hierarchy.** Whenever an item enters done, by any path, each item above it whose usage was not measured from a log is given the sum of its children's usage, cancelled and archived children included, up to its epic. A story measured from its log keeps its measurement, which holds its tasks' and the work between them.
3. **Charted by model.** `flai stats` reports, per item, tokens, cost, and tokens per hour of agent work, per model; and per model and type, the totals, and items done cumulatively against time and against cumulative cost. The dashboard charts them, as `design/system/metrics.md` defines.

## Consequences

- Stories worked by agents flai serve starts get tokens and cost without anyone entering them; others get what their tasks sum to, which is nothing unless a task was measured.
- A task's numbers are an apportionment and say so; a story's are what the harness reported whenever the session ended normally.
- Only the `claude-code` harness is read. Another harness that logs usage needs a reader of its own in `flai/internal/usage`.
- Front matter is parsed strictly, so a flai older than this refuses an item that carries `usage`, and with it the listing that reads the item. The flai on the host must be upgraded before any item is measured.
- `flai serve` now writes work items in the main checkout, as flai commands run by agents and the dashboard already do; the write re-reads the item and changes only its `usage`.
- `design/system/metrics.md` gains the usage metrics and charts, the contract between `flai stats` and the dashboard.

## Alternatives considered

- A price table in flai or its configuration: it goes stale with every price change and needs every model named; the logs already report the price paid.
- Output tokens from the streamed calls: they are the first chunk's only, a small fraction of what `result` reports.
- Measuring at every move into done, from any flai: only `flai serve` knows where its logs are and which project key names them, and a task that enters done mid-session has no `result` yet. The move sums front matter, which any flai can do; `flai serve` measures.
- Recording usage only on stories: the designer asked for tasks and epics too.
