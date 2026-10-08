---
id: ADR-0127
title: "flai guard lets a story's sub-agent write an ADR with flai adr new, topics, and accept, or adr_new, and refuses it the commit, which stays the story's agent's"
status: proposed
date: 2026-10-08
supersedes: []
superseded_by: []
refines: [ADR-0060]
topics: [cli, conventions, template]
---

# ADR-0127 flai guard lets a story's sub-agent write an ADR with flai adr new, topics, and accept, or adr_new, and refuses it the commit, which stays the story's agent's

## Context

[ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) has `flai guard` refuse a sub-agent any flai command that does not read, and every flai MCP tool but the reads, so that only the story's agent changes work items, threads, and history. An ADR is none of those: it is a document under `design/adrs`, with a row in its index, like the design and the code a task sub-agent edits. But `flai adr new`, `flai adr topics`, `flai adr accept`, and `adr_new` are not reads, so the guard refused them, and task sub-agents that had to record a decision worked around it ([I-0062](../issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md)):

- S-0249's sub-agent was refused `flai adr new`, even with `--print-body`, and wrote ADR-0085 by hand, copying the number, the file name, `refines`, the index row, and the topics that `internal/adr` writes.
- S-0220's sub-agent was refused it too, and cited ADR-0090 by guessing the number before the story's agent made it.
- S-0207's sub-agent had `flai adr new` let through by an older guard, then `flai adr topics` on the same ADR refused, so the story's agent finished the ADR.

Each of these commands can also commit: `--commit` commits what it wrote on the story's branch and widens the story's touches, `--autocommit` commits it in the checkout, and `adr_new` does the first with `commit`. A commit is history, and the touches are a work item's front matter.

## Decision

`flai guard` lets a story's sub-agent run `flai adr new` (with `--print-body` too), `flai adr topics`, and `flai adr accept`, and call `adr_new`, so long as none commits; it refuses each with `--commit` or `--autocommit`, in any form, and `adr_new` with `commit` true, saying the commit is the story's agent's.

- What the sub-agent wrote stays uncommitted in the worktree, and the story's agent reviews it and commits it with the task, as it does the task's other files.
- The planner, the orchestrator, and the analyzer are not changed: their rules still refuse every `flai adr` command and `adr_new`.

## Consequences

- A task sub-agent records the decision its task makes with flai, which numbers it from the whole repository, names the file, writes the front matter, and adds the index row, instead of copying `internal/adr` by hand or guessing the number.
- The new ADR and its index row are not in the task's `touches` when it was planned; `flai task done` widens the touches to them when the story's agent commits them.
- A sub-agent can now change an accepted ADR's `topics` or accept a proposed ADR, the two edits flai allows to them. These are document edits the story's agent reviews in the diff before it commits, as any file a sub-agent changes.
- `delegation.md`, its template copy, `docs/users/flai.md`, and `design/system/agent-context.md` say what a sub-agent may run.

## Alternatives considered

- Refuse every `flai adr` command and `adr_new` to a sub-agent, and say so in `delegation.md`: it keeps the hand-copied ADRs and guessed numbers I-0062 records.
- Let a sub-agent commit the ADR too: a commit is history and widens the story's touches, both the story's agent's under ADR-0060.
