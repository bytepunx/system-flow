---
id: ADR-0059
title: "A story's agent hands search, test runs, and verification to an explorer and a verifier sub-agent that can only read, primed by role, and asks the designer itself"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
topics: [cli, conventions, template]
---

# ADR-0059 A story's agent hands search, test runs, and verification to an explorer and a verifier sub-agent that can only read, primed by role, and asks the designer itself

## Context

An agent `flai serve` starts with the `claude-code` harness (ADR-0043) works its whole story in one context. Every search, test run, lint log, and file it reads stays there and is read again, from the cache, on every later turn: S-0118's run cost 5.18 US dollars with 10.5M cache-read tokens. Claude Code's `Task` tool (`Agent` in its log) is already in these headless sessions, and S-0118 used it once on its own, but nothing in the prompt, the conventions, or the template says when to use it.

A sub-agent shares its parent's MCP servers. Every call it makes to flai arrives on the parent's connection under the parent's name, `agent-S-nnnn`: left alone it can move the story, edit items, and consume the parent's inbox, whose cursor is per name.

A probe on 2026-10-01 (Claude Code 2.1.286, `claude -p` with flai's `--mcp-config` and `--strict-mcp-config`, as `flai serve` starts it) found:

- Agent definitions in the project's `.claude/agents/` load in a headless session started in the project, as they do interactively.
- A definition's `tools` restricts MCP tools as well as built-in ones: a sub-agent listed with `mcp__flai__item_get` and `mcp__flai__doc_search` had those two of flai's and no others.
- A definition's own `mcpServers` did not connect under `--strict-mcp-config`, so a sub-agent cannot be given a flai server of its own with its own name.
- In the stream-json log `flai serve` keeps, every event of a sub-agent carries `parent_tool_use_id`, the ID of the `Agent` call that started it, whose input names the definition (`subagent_type`).

Of flai's MCP tools, `thread_get`, `item_get`, `doc_get`, `doc_search`, `prime`, `who_touches`, and `board` read and use no agent identity; `inbox`, `wait_for_events`, and `wait_for_work` advance the caller's cursor; the rest write under the caller's name.

## Decision

A story's agent hands noisy work and pre-review verification to sub-agents of two roles, an explorer and a verifier, that the template defines with tools that only read, that prime with a pack for their role, and that return questions to the story's agent, which alone asks the designer.

1. **Roles.** The explorer searches and reads: code, design, and logs, and returns what it found with paths. The verifier reads and runs the project's tests, lint, and `flai check`, and returns what fails against the story's criteria and the conventions. Neither edits a file.
2. **Definitions.** The template ships `.claude/agents/explorer.md` and `verifier.md`. Their `tools` list is an allowlist: the built-in read tools (and `Bash` for the verifier, to run tests and lint), and of flai's tools only the reads above. Nothing that moves, creates, or edits an item, opens, replies to, or resolves a thread, or reads the inbox or waits for events or work. A session `flai serve` starts in the project loads them as it loads any project agent; the adapter adds no flag.
3. **Priming.** `flai prime --story S-nnnn --role explore|verify`, and `role` on the MCP `prime` tool, return a pack for the role: the conventions whose front matter `roles` lists it, with the sections the story's topics leave out taken out; the story's goal and acceptance criteria; and briefs of what the story names and what its topics and one link step select, in that order, while a budget of half the project's has room. Briefs that do not fit are counted, not printed; `doc_search` finds them.
4. **Questions.** A sub-agent that needs the designer says so in its final message, with the question and its recommended answer. The story's agent asks it with `thread_open` on the story or the task, and waits with `wait_for_events` as for its own questions.
5. **Attribution.** A sub-agent's flai calls are told apart from its parent's in the agent's log by `parent_tool_use_id`. The MCP server is not changed: the calls a sub-agent can make change nothing and are attributed to no one, so their name does not matter.
6. **When.** `harness.Prompt` for `claude-code` and a baseline convention say when to delegate (code search, test and lint runs, long logs, verification before review), what a sub-agent's prompt holds (the worktree, the story and task IDs, the question), and that it returns a summary, not raw output.

## Consequences

- The story's agent keeps decisions and edits in its context and a sub-agent's reading in another, which ends with the sub-agent's final message. What that saves is measured, not assumed: `agent-context.md` records it against runs that did not delegate, from ADR-0051's usage, which includes sub-agents.
- Conventions gain a front matter key, `roles`. A convention without it is read by the story's agent only. The baseline lists `explore` and `verify` on the files each needs; a project narrows or widens them as it does topics.
- The guarantee that a sub-agent cannot act for the story is the harness's tool allowlist. An operator who edits the definitions to add a write tool takes it away; the MCP server would not notice.
- Only the `claude-code` adapter is told about sub-agents. The `command` harness runs a program flai cannot see into, and other harnesses get definitions of their own when they have sub-agents.
- A sub-agent starts from its prompt alone: it sees none of the parent's pack, and a prompt without the worktree, the story, and the question gets a worse answer. The prompt and the convention say what to put in it.
- Forked sub-agents, which inherit the parent's context and cache, are not used here; they are what the task fan-out experiment after this story tries.

## Alternatives considered

- A sub-agent opens its own threads, attributed to the parent's name (option (b) in S-0175): the server reports to an agent only what others changed, so the parent's `wait_for_events` would never see a thread opened under its own name, and two voices would ask the designer about one story.
- The harness's own messaging between agents (`SendMessage`, agent teams) (option (c)): harness-specific, and a sub-agent that asked could not wait for the designer's answer anyway; the parent can.
- A flai MCP server per sub-agent, with its own name and only read tools: the strongest guarantee, but a definition's `mcpServers` does not connect under `--strict-mcp-config`, which `flai serve` needs to keep the operator's own servers out.
- A role flag on `flai mcp` that hides write tools from every caller: the parent shares the server and needs them.
- No definitions, only guidance to use the built-in Explore and general-purpose agents: they inherit every flai tool, so a sub-agent could move the story.
- A scout sub-agent that reads the whole pack for the parent: ADR-0049 rejected it; sub-agents here get less than the parent, not more.
