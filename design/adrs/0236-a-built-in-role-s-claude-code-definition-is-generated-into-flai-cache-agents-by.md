---
id: ADR-0236
title: "A built-in role's Claude Code definition is generated into .flai-cache/agents by flai role definition --write, a session started by hand reads it through a committed link .claude/agents/<name>.md the operator makes once, and the template ships no copy"
status: accepted
date: 2026-11-02
supersedes: []
superseded_by: []
refines: [ADR-0207]
---

# ADR-0236 A built-in role's Claude Code definition is generated into .flai-cache/agents by flai role definition --write, a session started by hand reads it through a committed link .claude/agents/<name>.md the operator makes once, and the template ships no copy

## Context

[ADR-0207](0207-a-strategic-role-is-a-markdown-file-under-design-workflows-roles-whose-front.md) § 4 has flai generate a strategic role's harness definition from its role file: `name`, `definition.description`, `definition.tools`, `definition.model` with the role's agent's model over it, and the body without the marker line as the prompt. A session `flai serve` starts takes it through `--agents`, so flai writes nothing under `.claude/`. It left open how a session the operator starts by hand reads the same definition, since Claude Code reads `.claude/agents/`, a folder Claude Code protects ([ADR-0106](0106-a-story-whose-branch-changes-a-path-claude-code-protects-is-accepted-by-the.md)): if Claude Code follows a symbolic link there, the file is a link to a definition flai generates outside `.claude/`; if not, a hand-started session takes `--agents` too.

Until S-0553 the planner's, the orchestrator's, the analyzer's, and the curator's `.claude/agents/<name>.md` were hand-written copies, here and in the template, and `flai serve` read their description, tools, and model, so a change to a strategic agent was made twice (D-0411).

S-0553 tried it with Claude Code 2.1.311, in a scratch repository, with `--model haiku`:

- `.claude/agents/linkbot.md`, a relative symbolic link to `defs/linkbot.md`, was found and run by `claude -p --agent linkbot`, and as a sub-agent through the Agent tool with `subagent_type: linkbot`, as `.claude/agents/copybot.md`, a regular file beside it, was.
- With the link's target moved away, `claude -p --agent linkbot` was refused, `--agent 'linkbot' not found`, and the list of available agents it printed left linkbot out. Nothing else failed.

So Claude Code follows the link, and a link whose target is missing only hides the agent.

## Decision

A built-in role's definition is generated into `.flai-cache/agents/<name>.md`, and a session started by hand reads it through `.claude/agents/<name>.md`, a relative symbolic link to `../../.flai-cache/agents/<name>.md` that is committed and that the operator makes once.

- `flai role definition <name>` prints the definition generated from the role's file, the project's or the embedded baseline, as a Claude Code agent file: front matter with `name`, `description`, `tools` as one comma-separated line, and `model`, then the prompt. The model is the file's `definition.model`; `flai serve` puts the role's agent's model over it at start, as ADR-0207 § 4 says. It takes a role a project defines as well as a built-in one.
- `flai role definition --write` writes the definition of each named role, or of every built-in role when none is named, to `.flai-cache/agents/<name>.md` in the checkout it runs in, which is git-ignored. For each built-in role whose `.claude/agents/<name>.md` is not the link, it prints the command that makes it, `ln -sf ../../.flai-cache/agents/<name>.md .claude/agents/<name>.md`, for the operator to run. flai never writes under `.claude/` itself.
- The generated file is never edited. After a role file changes, `flai role definition --write` writes it again; `flai serve` does not read it, since it generates the definition at every start.
- The template ships no `.claude/agents/<name>.md` for the planner, the orchestrator, the analyzer, or the curator. The explorer's and the verifier's stay hand-written there, being sub-agents' and not roles (ADR-0207 § 2). In a project made from the template, the operator runs `flai role definition --write` and the four commands it prints once, and commits the links.
- `flai check` reports `roles.definition` for a built-in role's `.claude/agents/<name>.md` that exists and is neither that link nor a file whose description, tools, model, and prompt equal the definition generated from the role file, naming the first that differs. A missing file is no finding, and neither is a link whose target is missing, as in a fresh clone or a story's worktree, since only a session started by hand reads it and `flai check` must pass on a fresh clone.
- `flai guard` reads a built-in role's tools from its generated definition, never from `.claude/agents/`.

## Consequences

- A strategic agent's definition is written in one place, its role file. The double edit of D-0411 ends.
- A session started by hand in a fresh clone, or in a story's worktree, has no planner, orchestrator, analyzer, or curator until `flai role definition --write` runs there. It is refused `--agent` with Claude Code's own message, which names the agents it has.
- A hand-started session reads the definition as it was last written. A role file changed since then reaches it only after `flai role definition --write` runs again.
- A project made from the template before S-0553 keeps its four copies, since `flai upgrade` removes no file the template dropped. `flai check` holds each to its role file, and the operator replaces it with the link once.

## Alternatives considered

- A link to the role file itself: Claude Code would read the role file's front matter, whose description and tools are under `definition`, so the agent would have none.
- Commit the generated definitions at a tracked path: a second copy of every role file in the repository, kept in step by a command and a check, which is the double edit made mechanical.
- Pass `--agents` to a session started by hand: works, but the operator must build the JSON for every session, and a session started from an IDE cannot take it.
- Have `flai role definition --write` make the link too: flai would then write under `.claude/`, and an agent running the command would reach a protected path through the shell.
- Teach the template's renderer to render a symbolic link: the link would dangle in every new project until `--write` runs, as it does now, for a change to the renderer, the lock file, and `flai upgrade`.
