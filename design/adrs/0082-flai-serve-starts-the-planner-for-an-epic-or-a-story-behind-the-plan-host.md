---
id: ADR-0082
title: "flai serve starts the planner for an epic or a story behind the plan host action, as a claude-code session that flai guard holds to planning"
status: accepted
date: 2026-10-04
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0065, ADR-0079]
topics: [cli, template, planning]
---

# ADR-0082 flai serve starts the planner for an epic or a story behind the plan host action, as a claude-code session that flai guard holds to planning

## Context

The planner drafts an epic's stories and enriches a story with its touches, a forecast, and a cost of delay (E-0016). S-0207 gave it a pack (`flai prime --role plan`, ADR-0075) and S-0206 an activity document (ADR-0079), but nothing started it. `flai serve` starts only a story's agent, in the story's worktree, when the story enters ready (ADR-0038, ADR-0043). The planner has no story of its own and works in the main checkout, where items live. It must write items and threads and nothing else: no code, no move past backlog, no acceptance, no release. `flai guard` held only sub-agents (ADR-0060), and Claude Code runs a `PreToolUse` hook only for the tools its matcher names, which were `Bash` and flai's MCP tools.

## Decision

The planner is a session `flai serve` starts for one epic or one story when the operator asks, behind the host action `plan`, and `flai guard` holds that session to planning.

1. **Asked for, never started by itself.** `flai plan <E-nnnn|S-nnnn>`, the host API's `plan.run` (`{id}`, for the dashboard's Plan), and the MCP tool `plan` (for the operator's own agent; refused to an agent `flai serve` started) start it. Each needs the `plan` host action on for the project, off by default. Nothing starts it on a move or a timer.
2. **Its agent.** `planning.agent` in `system-flow.yaml`, merged over the project's `agent` as a story's agent is (ADR-0037, ADR-0065), gives its harness, model, config, and roles.
3. **How it runs.** In the project's main checkout, as `planner-<item>`, with `FLAI_ROLE=plan` and `FLAI_ITEM=<item>` and no `FLAI_STORY`, and a prompt of its own. On `claude-code` it runs as the project's `.claude/agents/planner.md`, passed with `--agents` beside the explorer and verifier and selected with `--agent planner`, so the definition is the source of its tools and instructions as ADR-0065 has it for sub-agents. The template ships the definition. The operator's command gets `FLAI_ROLE` and `FLAI_ITEM`, with `{story}` empty.
4. **One run per item, recorded as a story's is.** The run is kept in `serve/agents.json` under `plans`, by item, with its log `<key>-planner-<start>.log` beside the story agents' logs. A second run for an item while one runs is refused. When it ends, flai serve judges it `asked` (a question of its own is open on the item), `failed` (a failed exit), or `worked`, and logs its activity in `wip/agents/planner.md` (ADR-0079). It does not count against the in-progress limit and holds no story back.
5. **Refusals.** The action off, a task or an ID that is neither an epic's nor a story's, an item archived, done, or cancelled, a planner already running for the item, and nothing to start it with.
6. **Held by the guard.** With `FLAI_ROLE=plan`, `flai guard` checks the session's own calls as well as its sub-agents'. It passes reads, the MCP tools `inbox`, `item_new`, `item_edit`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`, and `item_move` to backlog; and the commands `story new`, `epic new`, `edit` without `--no-draft`, `touches`, `thread new` and `reply`, `issue new` and `bump`, and `move` to backlog. It refuses every other flai tool and command, git's writes, and `Edit`, `Write`, and `NotebookEdit`. The planner's sub-agents are held as every sub-agent is.
7. **The settings run the guard on its edits.** `.claude/settings.json`, the template's and this repository's, keeps the `Bash|mcp__flai__.*` entry for every session and gains a second `PreToolUse` entry on `Edit|Write|NotebookEdit` whose command exits at once unless `FLAI_ROLE` is `plan`. The operator chose this on TH-0096.

## Consequences

- The operator plans an item from its page, a shell, or their own agent, and pays for a planner only when they ask.
- A planner writes items only through flai, so every write is checked by `flai check` and stamped with who set it (ADR-0074). Editing an item's words goes through `item_edit`, not the file.
- A story's session runs no hook per edit: the second entry is one shell test outside a planner session.
- A project whose `.claude/agents/` has no `planner.md` cannot start the planner on `claude-code`; the refusal names `flai upgrade`, which adds it from the template. A project that keeps its old settings has a planner whose file edits the guard does not see; `planner.md` leaves those tools out, which is the only fence then.
- ADR-0079 left the strategic agents' names to this story: the planner's is `planner-<item>`.
- The orchestrator and the analyzer will need the same: a host action, a role in `FLAI_ROLE`, a definition, and guard rules of their own.

## Alternatives considered

- **`--settings` on the planner's run**, a hook on `Edit|Write|NotebookEdit` passed only to the planner's session. It was built first, because the settings files could not be written from the story's session, and replaced on the operator's word (TH-0096): the settings files are where every project's hooks live and where the operator reads them.
- **Widening the one matcher to `Bash|Edit|Write|NotebookEdit|mcp__flai__.*` for every session.** Every story's agent would run the guard on every edit, and in this repository rebuild `bin/flai` for it, to refuse nothing.
- **Starting the planner on a move or on a timer**, as story agents start on ready. Planning spends money and writes drafts the operator then reads; the operator asks for it.
- **The planner in a worktree of its own.** Items live in the main checkout, and the planner changes nothing else.
- **A planner run counted against the in-progress limit.** It works no story, and holding a story back for it would trade delivery for planning.
