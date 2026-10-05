# Changelog

## 1.0.49 - 2026-10-05

- S-0268 The story agent's prompt says to batch independent edits and commands in one turn, to take a new task through ready in one command, and to prove a test fails without the shared stash (patch).

## 1.0.48 - 2026-10-05

- S-0249 flai check --strict stops a story's close-out on wip/ findings outside the story (patch).
- S-0266 The verifier reads a close-out run's exit status without piping it away, and the story agent tells it where the run is expected to stop (patch).

## 1.0.47 - 2026-10-05

- S-0249 flai check --strict stops a story's close-out on wip/ findings outside the story (patch): `scripts/close-out.sh` exports the story as `CLOSE_OUT_STORY`, and `scripts/check.sh` then runs `flai check --strict --story "$CLOSE_OUT_STORY" --record-issues`. A finding outside the story, and every `wip.overlap`, is printed as a note that does not stop the close-out, and each rule's are recorded in one open issue under `design/issues`, opened or bumped, which the close-out's commit step commits. Without `CLOSE_OUT_STORY`, as in CI and `make`, the check is unscoped and every finding counts. `work-management.md`'s close-out rule says the same. The scripts need a flai that has `--story` and `--record-issues`.

## 1.0.46 - 2026-10-04

- S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value (patch).

## 1.0.45 - 2026-10-04

- S-0210 The planner enriches a story with predicted touches, a forecast, and a cost of delay value (patch): `strategic-agents.md`'s section "As the planner" says to predict a story's `touches` from `flai touches suggest`, its goal and criteria, the design it links, and the code layout, keeping every touch it declares; to take its forecast from `flai forecast` and its cost of delay `value` from `flai cod`; to review each figure and adjust it with a stated reason; and to record where each touch came from and why each figure stands under a `### Planning` heading in the story's Notes, which is the planner's to rewrite. `.claude/agents/planner.md`, the planner's definition, says the same. The three commands print and write nothing, and need a flai that has them.

## 1.0.44 - 2026-10-04

- S-0255 The planner drafts a story's tasks into the backlog and revisits the children it already has (patch): `strategic-agents.md` says the planner drafts a story's tasks in layers, each with `## Work`, `## Done when`, a nature, tags, `touches`, and `after`, revisits the open ones, proposes splits, merges, and drops on the plan's thread, never cancels or rewrites a task it did not write without asking, adds a task's topic to the story, and names the items it created and revisited in its summary. `work-management.md` says the planner may draft a story's tasks, and that the agent that pulls the story reviews them before it works them.

## 1.0.43 - 2026-10-04

- S-0209 The planner drafts an epic's stories into the backlog and revisits the children it already has (patch): `strategic-agents.md`'s section "As the planner" says to create each story as a draft (`draft: true`), since flai refuses a planner's story that is not one, and to write every section of the item's template, since flai check refuses a story missing one; to summarise the plan in one thread on the item, naming on an epic the stories, their order, and the assumptions made; for an epic with stories, to enrich each one not done or cancelled again and propose in that thread each story to split, merge, add, or drop, drafting the additions only; and to end with one line naming the stories created and revisited. `.claude/agents/planner.md`, the planner's definition, says the same.

## 1.0.42 - 2026-10-04

- S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action (patch).

## 1.0.41 - 2026-10-04

- S-0208 The planner is an agent flai serve starts for an epic or a story, behind the plan host action (patch): the template ships `.claude/agents/planner.md`, the definition `flai plan` runs the planner's session as, which writes work items and threads through flai only. `.claude/settings.json` keeps its hook on `Bash` and flai's MCP tools and gains a second `PreToolUse` entry that runs `flai guard` before `Edit`, `Write`, and `NotebookEdit` in a planner's session alone (`FLAI_ROLE=plan`), so the guard can refuse the planner a file edit while a story's session runs no hook per edit (ADR-0082, the operator's choice on TH-0096).

## 1.0.40 - 2026-10-03

- S-0206 Planner, orchestrator, and analyzer runs are logged in an activity document, with cost, seconds, and work completed in its front matter (patch): `wip/agents/README.md` names the three activity documents flai writes, and `strategic-agents.md` says `activity_log` takes the kind, the summary, and the items touched, that flai measures the duration and cost, and that a run's final reply summarises what it did after its last reported activity (ADR-0079).

## 1.0.39 - 2026-10-03

- S-0200 An epic follows its stories: to ready and in-progress with the first, to review and done with the last (patch): `git.md`'s release rule says an epic contributes no bump of its own, its stories carry theirs, now that an epic is accepted with its last open story (ADR-0078, the designer's choice on TH-0082).

## 1.0.38 - 2026-10-03

- S-0207 The conventions and prime know the planner, orchestrator, and analyzer as roles (patch): `strategic-agents.md` says what the planner, the orchestrator, and the analyzer do, never do, log, and ask, read by the roles `plan`, `orchestrate`, and `analyze`. The baseline's `roles` name the new roles where they apply, as the designer decided on TH-0081: `work-management.md` adds `plan` and `orchestrate`, `decisions.md` `orchestrate`, `continuous-improvement.md` and `code-quality.md` `analyze`, and `telemetry.md`, read by every agent until now, lists `story`, `explore`, `verify`, and `analyze`.

## 1.0.37 - 2026-10-03

- S-0199 Work items carry planning data: draft, cost of delay inputs and value, and a forecast with who set it (patch): `work-management.md`'s definition of ready says a draft story is finalized before it is ready (ADR-0074).

## 1.0.36 - 2026-10-03

- S-0243 Review column exceeds its WIP limit while acceptance is batched (patch).

## 1.0.35 - 2026-10-03

- S-0243 Review column exceeds its WIP limit while acceptance is batched (patch): `work-management.md` says `wait_for_work` names a story to pull only while review is under its limit as well as the in-progress limit, and that on a timeout it may be waiting for a story in review to be accepted or sent back (ADR-0073, with the designer's consent on TH-0079).

## 1.0.34 - 2026-10-03

- S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve (patch).
- S-0242 A publish refused because the remote moved drops its tags, takes the remote by rebase or merge, and works conflicts through tasks and threads (patch).
- S-0198 Give the operator the option to have all issues turned into stories (patch).

## 1.0.33 - 2026-10-03

- S-0197 Agents commit each task after flai stream sync, and stream sync makes conflicts easy to resolve (patch): `git.md` says to use `flai stream sync` for the branch's git operations and never start a `git rebase` or `git merge` by hand, finishing a rebase the sync stopped on with `git add` and `git rebase --continue`, or undoing it with `git rebase --abort`; to commit each task when it is done, then sync and resolve, then run its tests and commit any fix; and before review to commit what is outstanding, sync again, and close out. `work-management.md` points to that cycle for each task, before review, and on an `overlapped` change, and its definition of done asks for the branch synced after its last commit. `delegation.md` syncs again before the verifier runs. `scripts/close-out.sh` stops while a rebase is unfinished, and when the branch does not contain the main branch, saying to run `flai stream sync`: before the tests when there is nothing to commit, after the commit otherwise.

## 1.0.32 - 2026-10-02

- S-0194 An experiment story is accepted with its results document under design/experiments (patch).

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
