---
id: ADR-0261
title: "A story or strategic agent that names no harness falls back, after the project's default and the project's host command, to flai's built-in default agent, claude-code, the template's manifest names that default, the host's agent settings are kept per project over a host-wide fallback, and a refusal that remains names the setting, the dashboard page that sets it, and the command"
status: accepted
date: 2026-11-07
supersedes: []
superseded_by: []
refines: [ADR-0037, ADR-0251]
---

# ADR-0261 A story or strategic agent that names no harness falls back, after the project's default and the project's host command, to flai's built-in default agent, claude-code, the template's manifest names that default, the host's agent settings are kept per project over a host-wide fallback, and a refusal that remains names the setting, the dashboard page that sets it, and the command

## Context

S-0587 found that in a new project, once the `agent` host action is on, the dashboard's Start agent for a new story is refused:

```text
S-nnnn names no harness, and no command is set on the host (flai serve agent set -- <program> [args...])
```

[ADR-0251](0251-a-story-whose-agent-names-no-harness-is-started-with-the-project-s-default.md) starts a story that names no harness with the project's default agent. A new project has none: the template's `system-flow.yaml` ships `agent:` commented out, `flai new` and `flai import` write nothing else, and flai has no default of its own. The `claude-code` the dashboard shows in the agent fields is a placeholder, not a stored default. The planner, the orchestrator, the analyzer, the curator, and a project's roles are refused the same way.

Two more things the operator asked for:

- When something is still missing, the refusal should point at the page where it is set.
- The host's agent settings, the command (`flai serve agent set`), the session name, and the harness programs and arguments (`flai serve agent harness`), live once in `~/.flai/config.json` for every project on the host. Setting a harness's program for one project changes it for all of them. Only the host actions are kept per project.

## Decision

Every agent flai starts names a harness: what names none falls back to the project's default, then to the project's host command, then to flai's built-in default, `claude-code`; the template names that default; the host's agent settings are per project; and each refusal that remains names its fix.

### 1. A built-in default agent

**flai has a built-in default agent: harness `claude-code`, no model, no config.** An agent with no model runs on the harness's own default model; the claude-code adapter passes no `--model` then. The agent a story's agent or a strategic agent is started with resolves in this order, the first that names a harness winning:

1. the story's own agent (ADR-0037), or the role's own block (`planning.agent` and the rest);
2. the project's default, `agent` in `system-flow.yaml` (ADR-0251);
3. the host's command for the project (§3), when one is set: the agent then runs on the `command` harness, as ADR-0251 kept;
4. the built-in default.

What the earlier steps say beside a harness, a model or a config key, is laid over the step that names one, field by field (`Agent.With`), as ADR-0251 lays the story's words over the default.

- Nothing is written. A story made in a project with no default still carries no `agent`, and `flai story new` copies only the project's default (ADR-0037), never the built-in.
- So no start is refused for naming no harness. The refusal `names no harness, and no command is set on the host` goes.
- `harness: none` still names a harness ([ADR-0242](0242-a-story-s-agent-whose-harness-is-none-runs-elsewhere-and-flai-serve-starts.md)): such a story is still run elsewhere.
- The dashboard's story page and `flai serve`'s journal name the harness and model of the agent that will run, as resolved here.

### 2. The template names the default

The template's `system-flow.yaml` ships `agent: { harness: claude-code }` set, not commented out. So `flai new` and `flai import` write it into every new project, and the project states its choice where the operator can see and change it.

The model is left unset, so it is the harness's own default, which does not go stale with a model release. Every story made in such a project carries a copy of it (ADR-0037).

### 3. The host's agent settings are per project

The host configuration's `agent` block gains `projects`, a map from a project's root to that project's own `command`, `name`, and `harnesses`, keyed as `host_actions` is.

- For a project, flai serve reads the project's entry over the host-wide `agent.command`, `agent.name`, and `agent.harnesses`:
  - a `command` or a `name` set for the project wins;
  - a harness entry set for the project replaces the host-wide entry of that harness, whole.
- The host-wide values stay the fallback for every project that sets none. A configuration written before this loads and runs as it did, with nothing moved.
- `flai serve agent set`, `flai serve agent clear`, and `flai serve agent harness` write the entry of the project in the working directory, or of `--project <root>`. They write the host-wide fallback only with `--all-projects`, as `flai serve enable` and `disable` take it. Outside a project, with neither flag, they refuse and name both flags.
- The dashboard's Harnesses page edits the project it is showing, and the host-wide fallback only when the operator chooses every project.
- The providers (`agent.providers`) and `agent.auto_restarts` stay host-wide: a provider is a gateway on the host, not a project's choice.

### 4. A refusal names its fix

Every refusal and `waiting` reason that remains about a story's or a strategic agent's harness says three things: what is missing, the dashboard page that sets it by its path, and the command. An example is a harness named `command` with no command set for the project, or a harness flai cannot start:

```text
S-0042's harness is the host's command, and none is set for this project: set it on the dashboard's Harnesses page, /settings/harnesses, or with flai serve agent set -- <program> [args...]
```

- The pages are `/settings/harnesses`, for the command and the harness programs, and `/settings/agents`, for the project's default agent.
- The dashboard renders a `/settings/<page>` path in such a message, on the story page and in the host agent notice, as a link to that page. The words stay plain text, so the CLI and the journal read the same.

## Consequences

- A new project starts a story's agent with no setting beyond turning the `agent` host action on, provided `claude` is on the host's `PATH` or the claude-code harness's program is set.
- A project that relied on the host's command with no default of its own keeps it: the command still comes before the built-in default.
- A host whose harness program, such as the claude-code program, differs between projects sets it per project. Setting it with `--all-projects` changes every project that has no entry of its own.
- The settings index, the operator guide, and the CLI reference gain `agent.projects` and `--project` and `--all-projects` on the three commands.

## Alternatives considered

- **Built-in default ahead of the host's command.** It is simpler: an agent would then always name a harness without the host. But a host that set a command and no default would start a different program after an upgrade, without a word. ADR-0251 kept the command as the fallback.
- **The template's default alone, with no built-in.** A project imported before this, or one whose operator removed `agent`, would still be refused.
- **Copy the built-in into every story made.** ADR-0251 already declined writing the default into stories. It would freeze a choice flai makes into every story's file.
- **A built-in model.** A fixed model name in flai goes stale with each model release and differs by provider. The harness's own default does not.
- **Per-project settings in the project's `system-flow.yaml`.** A program path and its arguments are facts about the host, not the project, and differ between the operator's machines. A program in the manifest would also run whatever a commit names.
- **Structured links in the refusal.** A refusal is a plain string from the CLI's error through the host API to the dashboard. A path in the words reads the same in a terminal and renders as a link in the dashboard, with no change to the error's shape.
