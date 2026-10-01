# Changelog

## 1.0.32 - 2026-10-02

- S-0194 An experiment story is accepted with its results document under design/experiments (patch).
- S-0176 The story's agent plans which tasks can run in parallel and works them with sub-agents (patch): `work-management.md` says the agent plans a story's tasks as it writes them, giving each its `touches` and the tasks of the story it waits for (`after`), so that tasks with no `after` between them and no path in common form layers that run together; it records the plan's reasoning in the narrative's `## Decisions` and works the plan layer by layer. `delegation.md` adds the task sub-agent: the story's agent runs a layer's tasks at once, one each (a fork where the harness offers one), sharing the story's worktree, or each in a worktree of its own when one builds or tests what another changes, waits for all, checks the files each changed against its `touches`, reviews each one's work, and commits and moves each task itself; a task sub-agent edits only what its task touches, runs only its own tests, never commits or writes through flai, and returns its question for the designer in its final message.

## 1.0.31 - 2026-10-01

- S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs (patch).

## 1.0.30 - 2026-10-01

- S-0189 Sub-agents run on a cheaper model, and the verifier's run replaces the agent's own full-suite runs (patch): `.claude/agents/explorer.md` runs `haiku` and `verifier.md` runs `sonnet`, cheaper than the story's agent, since neither decides anything. `delegation.md` says the story's agent runs only the tests for what it changed and leaves the whole suite, the lint, and `flai check` to one verifier before review, through the close-out script, and one more after fixing what that one found; it fixes what a verifier finds itself, never through a sub-agent, and does not repeat a verifier's passing run.

## 1.0.29 - 2026-10-01

- S-0187 The conventions say how agents write shell commands for zsh and gate each step of a chain (patch): `tooling.md` says the host shell may be zsh, so quote globs, URLs, and variables, use `[ a = b ]`, and put long sequences in a script, and that each step of a chain that must succeed is joined with `&&` or runs under `set -e`, never `;` or a pipe that hides its exit code. `scripts/close-out.sh` runs the lint, the three test tiers, `flai check --strict`, the narrative check, and the commit before review, stopping at the first failure, and `work-management.md` points to it; `make lint-md` calls the new `scripts/lint-md.sh`.

## 1.0.28 - 2026-10-01

- S-0179 flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading (patch).

## 1.0.27 - 2026-10-01

- S-0179 flai lints the markdown it writes in the main checkout, and a thread's entries in one second share a heading (patch): `tooling.md` says flai lints what it writes in `wip/` with the project's markdownlint configuration, warning in `flai check` and refusing a write that would bring a finding (ADR-0061).

## 1.0.26 - 2026-10-01

- S-0175 Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story (patch).

## 1.0.25 - 2026-10-01

- S-0175 Agents the claude-code adapter starts hand noisy work and verification to sub-agents that cannot act for the story (patch): `.claude/agents/explorer.md` and `verifier.md`, sub-agents whose tools only read; `.claude/settings.json`, which runs `flai guard` before a sub-agent's shell and flai calls; a baseline `delegation.md` convention; and `roles` on the conventions each sub-agent reads with `flai prime --role` (ADR-0059, ADR-0060).

## 1.0.24 - 2026-09-29

- S-0149 Provide briefs in cases where large files are part of context (patch).

## 1.0.23 - 2026-09-29

- S-0149 Provide briefs in cases where large files are part of context (patch): CLAUDE.md and session-start.md say a large file a story names only by its path written out is briefed (ADR-0050).

## 1.0.22 - 2026-09-29

- S-0148 Agents flai starts, and any agent with a story, prime with the budgeted pack and fetch what it briefs (patch).

## 1.0.21 - 2026-09-28

- S-0134 Conventions, design, tech files, and ADRs carry topics on the file and on headings, and flai check keeps them honest (patch).

## 1.0.20 - 2026-09-27

- S-0140 Agents should not leave uncommitted work in their branch when moving to review (patch).

## 1.0.19 - 2026-09-26

- S-0132 Accepting a story tells every open story that overlaps it which paths changed (patch).

## 1.0.18 - 2026-09-23

- S-0097 Agents should connect to the MCP and consistently monitor for ready tickets (patch).

## 1.0.17 - 2026-09-21

- S-0087 Moving a story to done merges it; a publish button on the board tags and pushes everything accumulated since the last one (patch).

## 1.0.16 - 2026-09-20

- S-0079 A story moved to ready starts an agent on the host, with the command the operator configured (patch).

## 1.0.15 - 2026-09-20

- S-0076 flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes (patch).

## 1.0.14 - 2026-09-20

- S-0070 Cancelling an epic cancels its open stories and their tasks, and cancelling a story cancels its open tasks (patch).

## 1.0.13 - 2026-09-19

- S-0060 Operators create ADRs from the dashboard (patch).

## 1.0.12 - 2026-09-19

- S-0063 An acceptance that is not pushed stays visible on the board and in inbox, and flai push --pending pushes it (patch).

## 1.0.11 - 2026-09-19

- S-0053 Research and experiment stories can be accepted: findings land on main and are pushed, with no release (patch).

## 1.0.10 - 2026-09-19

- S-0057 Cards can be dragged within a column to change the pull order (patch).

## 1.0.9 - 2026-09-19

- S-0058 An agent learns through MCP that the designer moved work to ready (patch).

## 1.0.8 - 2026-09-18

- S-0051 Accepting from the dashboard shows uncommitted changes outside wip and lets the designer include them (patch).

## 1.0.7 - 2026-09-18

- S-0049 A story is ready without tasks; the agent that pulls it writes them (patch).

## 1.0.6 - 2026-09-18

- S-0039 flai mcp: inbox, reply, read, and move tools with the template .mcp.json (patch).

## 1.0.5 - 2026-09-18

- S-0038 Comment threads and answers as files under wip (patch).

## 1.0.4 - 2026-09-18

- S-0037 Story branches with wip on main, stream sync, and the touches flag (patch).

## 1.0.3 - 2026-09-17

- S-0035 flai dashboard runs from a private registry or a local build (patch).

## 1.0.2 - 2026-09-17

- S-0033 Four-digit work item IDs (patch).

## 1.0.1 - 2026-09-17

- S-0021 flai template push publishes template changes to a remote (patch).

## 1.0.0 - 2026-09-16

Epic E-0005 (agent conventions) complete: twelve baseline conventions under
`design/conventions`, `design/issues` with `flai issue`, `scripts/` behind
the Makefile with three test tiers, and a `CLAUDE.md` that primes every
session from the conventions. First major release of the template.

## 0.4.0 - 2026-09-16

- `design/conventions/logging.md` and `telemetry.md` are active: five log levels including fatal, structured events, `/_health`, `/_ready`, `/metrics`, OpenTelemetry, golden signals (S-0030).

## 0.3.0 - 2026-09-16

Note: under the refined release rule (incidental additive touches are a patch) this would have been 0.2.1; the version stands as cut.

- `design/conventions/logging.md` and `telemetry.md` added as drafts, indexed at 110 and 120 (S-0025).
- Markdownlint excludes test fixtures and `bin/`.

## 0.2.0 - 2026-09-16

- `design/conventions/`: ten baseline agent norms with a project-additions marker (S-0023, S-0026).
- `design/issues/` skeleton and `scripts/` with test tier stubs and Makefile delegation (S-0026).
- `CLAUDE.md` opens with a priming section and is a map; norms live in the conventions (S-0024).
- Rendered projects are markdownlint-clean; lint config aligned with the item format.

## 0.1.0 - 2026-09-15

Initial prototype. Layout, work item schema, baseline CLAUDE.md, devex files.
