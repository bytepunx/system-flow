---
title: Agent narrative
updated: 2026-09-19
status: active
---

# Agent narrative

`wip/agents` exists so that a crash, a context reset, or a hand-off between agents costs minutes, not hours. It is the curated memory of each active work stream, written by the agent doing the work, for the next agent who picks it up.

## Files

```text
wip/agents/
├── README.md        # this convention in brief, shipped by the template
├── index.md         # table of active streams, maintained by flai and agents
└── S-0004.md         # one narrative per active story, named by story ID
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
---

# S-0004 CLI scaffold and config

## Context
Two or three paragraphs a fresh agent needs before touching anything. What the story is, what is already true in the repo, what constraints apply. Rewritten as understanding improves.

## Current state
Where things stand right now. Which tasks are done, what is half done, what files are dirty. Rewritten on every meaningful step. This is the first thing a resuming agent reads.

## Next steps
Ordered list. The first item is what to do next. Rewritten on every meaningful step.

## Decisions
Bullet list of decisions made during this stream with a one-line rationale each. Anything architectural also gets an ADR; link it.

## Open questions
Questions for the human. Each has a date. Answered questions move to Decisions. Between `<!-- threads:start -->` and `<!-- threads:end -->` flai keeps a generated list of unresolved threads anchored to the story or its tasks (ADR-0020, S-0038); do not edit that block by hand. Agents connected to `flai mcp` see the same threads through `inbox` and are woken by `wait_for_events` (S-0039).

## Log
Append-only. One entry per meaningful step, newest last.

### 2026-09-15T16:12:00Z
Started. Pulled S-0004 to in-progress. Read design/system/flai-cli.md.

### 2026-09-15T16:40:00Z
T-0021 done. Config read/write with tests. Decided on plain encoding/json over viper, see Decisions.
```

## What an agent is told through MCP

`flai mcp` is the agent's view of the repository ([ADR-0020](../adrs/0020-files-plus-mcp.md)); the files stay the record. An agent on the host starts it over stdio from `.mcp.json`; one that cannot start a process there reaches the same server over HTTP after `flai mcp start` ([ADR-0030](../adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md)), named by its `X-Flai-Agent` header as `FLAI_AGENT` names it locally, so the cursor and everything below are the same on both. Since S-0058 the view covers work as well as threads, and it serves an agent that ends its turn between the designer's messages as well as one that stays running.

- `inbox` reports three things. Threads awaiting this agent. The stories that are ready to pull, in pull order, with whether the in-progress limit allows a pull: this is state, listed on every call, so an agent that has forgotten everything still sees the work. And the changes others made since this agent last looked: an item created, moved, blocked, or unblocked, with what happened, to what, by whom, and when, and a note when the pull order changed.
- Changes are derived from the files, not recorded anew: `transitions` carry `to`, `at`, and `by`; `blocked` intervals carry `from`, `until`, and `reason`; `board.md` carries `order`. A change whose `by` is this agent is not reported back to it, which is why `flai move`, `block`, and `unblock` record `FLAI_AGENT` when it is set.
- "Since this agent last looked" is a cursor per agent name under `.flai-cache/mcp/`, outside git. It is a read marker, a cache in the sense of the overview's principle that tooling never owns state: losing it repeats or skips the report of a change and loses nothing else, because ready work and open threads are state and are always listed. With no cursor, the last 24 hours are reported.
- `wait_for_events` returns at once when the cursor is already behind, so a change made between two calls is not lost; otherwise it blocks until something changes. It returns events in the same shape as `inbox` does, with the changed paths, and advances the cursor.
- `board` returns what `flai board --json` prints: columns, limits, pull order, breaches.

An agent that ends its turn calls `inbox` at the start of every turn. An agent that stays running holds `wait_for_events`. Either way a ready story found there is pulled when nothing is in progress, without waiting to be told.

## Obligations

An agent working in a conforming repo must:

1. Open the narrative before making the first change for a story.
2. Rewrite `## Current state` and `## Next steps` after every task transition and before any long-running operation.
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
