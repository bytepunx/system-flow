---
id: ADR-0186
title: "flai guard lets every agent run flai components list, and a story's task sub-agent, but no read-only sub-agent and no strategic agent's own call, run flai components serve"
status: proposed
date: 2026-10-22
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0141]
topics: [cli, components]
---

# ADR-0186 flai guard lets every agent run flai components list, and a story's task sub-agent, but no read-only sub-agent and no strategic agent's own call, run flai components serve

## Context

[ADR-0060](0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md) has `flai guard` refuse a sub-agent every `flai` command that is not in its list of reads. `flai components` is not in that list, so every sub-agent is refused `flai components list` and `flai components serve`, through `bin/flai` and `scripts/flai.sh` alike.

[I-0165](../issues/I-0165-flai-guard-refuses-a-story-s-sub-agent-flai-components-list-and-flai-components-serve-so-a-headless-pass-of-the-harness-falls-to-the-story-s-agent.md) and [I-0156](../issues/I-0156-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md) record what this costs. S-0449's task sub-agent wrote a pass of the component harness and could not run it. The story's agent had to run the pass in its own context.

The three subcommands do different things:

- `flai components list` reads. It runs the sub-project's own `svelte/compiler` on its components and runs no component.
- `flai components serve` runs the sub-project's own Vite dev server on 127.0.0.1, as [ADR-0163](0163-flai-runs-a-sub-project-s-own-vite-dev-server-on-127-0-0-1-behind-the.md) decides. It writes a Vite config under `.flai-cache/components/`. With `--story`, it installs the worktree's dependencies when its `vite` is missing. It changes no tracked file.
- `flai components config` rewrites `system-flow.yaml`.

[ADR-0141](0141-flai-guard-refuses-a-read-only-sub-agent-a-shell-command-that-writes-a-file.md) has the guard refuse a read-only sub-agent, one whose definition grants neither `Edit` nor `Write`, any shell command that writes a file outside the temporary folder.

## Decision

`flai guard` lets every agent run `flai components list`, and lets a story's task sub-agent run `flai components serve`, but refuses serve to a read-only sub-agent and to a strategic agent's own call.

- **List.** `flai components list` is a read, and so is `flai components` with no subcommand, which prints its help. Every agent the guard holds to reads may run them: every sub-agent, read-only ones among them, and the planner's, the orchestrator's, and the analyzer's own calls.
- **Serve.** A story's task sub-agent may run `flai components serve`, with or without `--story`. A read-only sub-agent is refused it, since serving writes under `.flai-cache` and may install dependencies, which ADR-0141 forbids it. The refusal says why and tells it to ask the story's agent, or a task sub-agent, to run it.
- **Strategic agents.** The planner's, the orchestrator's, and the analyzer's own calls are still refused `flai components serve`. Running the project's code is not part of planning, orchestrating, or analysing.
- **Config.** `flai components config` stays refused to every sub-agent, since it rewrites the manifest. The story's agent runs it.

## Consequences

- A task sub-agent can write a pass of the component harness and run it, so the pass stays out of the story's agent's context.
- A task sub-agent's dev server runs the project's code and its dependencies on 127.0.0.1. A task sub-agent could already do this by running `vite` or a package manager through the shell, which the guard does not judge.
- The guard decides `flai components serve` by whether the sub-agent is read-only, which it reads from the definition as ADR-0141 does. A definition it cannot read is not read-only, so the guard fails open, as ADR-0060 says.
- Only Claude Code's sessions are covered, as with ADR-0060.

## Alternatives considered

- Let every sub-agent serve: a read-only sub-agent would write under `.flai-cache` and install dependencies, which ADR-0141 forbids.
- Let only `flai components serve --story` through: serving the main checkout's sub-project writes nothing more than serving a worktree's, and a task sub-agent may need either.
- Leave serve to the story's agent: this is the cost I-0165 records.
