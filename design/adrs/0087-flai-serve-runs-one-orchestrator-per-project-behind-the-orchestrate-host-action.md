---
id: ADR-0087
title: "flai serve runs one orchestrator per project behind the orchestrate host action, as a long-running session that flai guard holds to the operator's orchestration.permissions and that logs each decision"
status: accepted
date: 2026-10-05
supersedes: []
superseded_by: []
refines: [ADR-0060, ADR-0082]
topics: [cli, template, orchestration]
---

# ADR-0087 flai serve runs one orchestrator per project behind the orchestrate host action, as a long-running session that flai guard holds to the operator's orchestration.permissions and that logs each decision

## Context

The orchestrator keeps a project's work moving: it asks the planner to plan, finalizes drafts, promotes and orders stories, answers threads, accepts stories, and publishes releases (E-0016). S-0207 gave it a pack (`flai prime --role orchestrate`, ADR-0075), S-0206 an activity document (ADR-0079), and S-0217 the figures it decides by: `flai order --by`, `flai promote --candidates`, and `flai release --evaluate`. Nothing started it. The planner, the one strategic agent `flai serve` runs, plans one item when asked and ends (ADR-0082). The orchestrator has no item: it reacts to whatever changes in the project, for as long as the operator wants it to. Each of its acts is one the operator does today, and some, accepting and publishing, act with the operator's credentials, so the operator must be able to give them one at a time, and an agent's prompt is not a fence. `flai guard` held sub-agents (ADR-0060) and the planner (ADR-0082) to fixed rules; it had no rules that depend on the project.

## Decision

`flai serve` runs one orchestrator per project for as long as the operator has the `orchestrate` host action on, and `flai guard` holds that session to the permissions the operator sets in `orchestration.permissions`.

1. **Behind a host action.** `orchestrate` is off by default. `flai serve enable orchestrate`, or the dashboard's Settings page while `settings` is on, turns it on for a project. While it is on, `flai serve` starts the orchestrator at its next look, starts it again when a run ends, no sooner than a minute after a failed run, and stops it when the action is turned off. A run stopped so is logged as `stopped: orchestrate turned off`. A start refused, for no harness and no command, or no orchestrator definition, is recorded once, not every minute.
2. **Its agent.** `orchestration.agent` in `system-flow.yaml`, merged over the project's `agent` as `planning.agent` is (ADR-0082), gives its harness, model, config, and roles.
3. **How it runs.** In the project's main checkout, as `orchestrator`, with `FLAI_ROLE=orchestrate` and no `FLAI_ITEM` or `FLAI_STORY`. On `claude-code` it runs as the project's `.claude/agents/orchestrator.md`, which the template ships, selected with `--agent orchestrator`. The definition lists no `Edit`, `Write`, or `NotebookEdit`.
4. **A loop, not a task.** Its prompt has it prime, read the inbox and the board, act only within its permissions and by `orchestration.policy`, take every figure from flai's commands, log each decision with `activity_log`, kind `orchestrator`, with what it did, why, and the policy figure behind it, then hold `wait_for_events` and decide again, without ending. Where the guard refuses it, it asks the operator on a thread and never works around the refusal.
5. **Permissions, each off by default.** `orchestration.permissions` has six booleans, `plan_backlog_epics`, `finalize_drafts`, `promote_to_ready`, `order_ready`, `accept_reviews`, and `publish`, and `answer_threads`, which is `off`, `recommend`, or `autonomous`. An unknown key or a bad value is a `manifest.orchestration` finding of `flai check`, so a misspelt permission is not silently off.
6. **Held by the guard.** With `FLAI_ROLE=orchestrate`, `flai guard` reads the permissions from the project's manifest at each call and checks the session's own calls. Reads, `inbox`, `activity_log`, `wait_for_events`, `thread_open`, `thread new`, and `issue new` and `bump` always pass. Each permission allows its own calls. Every other write is refused, whatever the permissions: file edits, git's writes, and every other flai tool and command. A refusal a permission would allow names it, `needs orchestration.permissions.<name>`; any other says the orchestrator never does it.
7. **Refusals are logged.** The guard appends each refusal of the orchestrator's own call to `wip/agents/orchestrator.md`, under `## Refusals`: its time, the call, and the permission it needs, or `none`. A refusal has no seconds or cost and is not an activity. Writers of the document take a lock per activity kind, so that flai serve logging a run and the guard logging a refusal do not lose each other's entry.
8. **Recorded beside the planner's runs.** The run is kept in `serve/agents.json` under `orchestrator`, the newest alone, with its log `<key>-orchestrator-<start>.log`. It does not count against the in-progress limit and holds no story back. When it ends, `flai serve` logs its activity since the last `activity_log` in `orchestrator.md`. Starts, ends, failures, and stops are journal entries with action `orchestrate`.
9. **The settings run the guard on its edits.** The `Edit|Write|NotebookEdit` entry in `.claude/settings.json` runs the guard when `FLAI_ROLE` is `plan` or `orchestrate`.

## Consequences

- The operator turns the orchestrator on and off as a whole, and gives it one power at a time in the manifest, which is committed and reviewed like any other change. A permission takes effect at the orchestrator's next call, with no restart.
- An orchestrator with no permission on reads, logs, and asks. It costs a session for as long as the action is on, though it mostly waits.
- The figures it acts on are flai's, so the operator and the dashboard see the same order, candidates, and release answer the orchestrator acted on.
- `orchestrator.md` holds both what it decided and what it was refused. A refusal that names a permission tells the operator what to turn on if they want it done.
- The answer mode, `recommend` or `autonomous`, is told to the orchestrator by its prompt. The guard lets it reply on a thread with either and refuses it with `off`.
- A project with no `.claude/agents/orchestrator.md` cannot start the orchestrator on `claude-code`; the refusal names `flai upgrade`. A project that keeps its old settings has an orchestrator whose file edits the guard does not see; the definition leaves those tools out, which is the only fence then.
- The analyzer will need the same: a host action, a role, a definition, and guard rules of its own.

## Alternatives considered

- **Starting the orchestrator per event**, as the replanner starts the planner (ADR-0084). Every event would pay for a session's start and its priming again, and a decision that spans events, such as waiting for an answer, would be lost between runs.
- **Permissions enforced by the prompt alone.** An agent may misread or ignore its prompt, and accepting and publishing act with the operator's credentials. The guard is the fence; the prompt only says what the fence allows.
- **One switch for everything the orchestrator does.** An operator who trusts it to order the ready column may not trust it to accept or publish.
- **Refusals as activity entries under `## Log`.** A refusal has no seconds or cost, and counting it as an activity would inflate `tasks_completed` and the totals. A section of its own keeps the totals honest and the refusals easy to read.
- **Restarting at once after any end.** An orchestrator that fails at start, for a bad model name or a missing definition, would be started again every few seconds. A minute's wait after a failure, and one record of a refused start, bound that.
- **Leaving a run to end by itself once the action is off.** It holds `wait_for_events` and does not end; turning the action off ends its process group at flai serve's next look, within a minute.
