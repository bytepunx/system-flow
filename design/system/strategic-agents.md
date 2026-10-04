---
title: Strategic agents
updated: 2026-10-04
status: active
topics: [planning, orchestration, analysis]
---

# Strategic agents

The planner, the orchestrator, and the analyzer are agents that work above a story (E-0016). What each does and never does is the convention [strategic-agents.md](../conventions/strategic-agents.md). This document is how flai runs them. So far the planner runs; the orchestrator and the analyzer have their packs and their activity documents, and nothing starts them yet.

| Part | Planner | Orchestrator | Analyzer |
|------|---------|--------------|----------|
| Pack: `flai prime --role` ([ADR-0075](../adrs/0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md)) | `plan`, with `--epic` or `--story` | `orchestrate` | `analyze` |
| Activity document ([ADR-0079](../adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)) | `wip/agents/planner.md` | `wip/agents/orchestrator.md` | `wip/agents/analyzer.md` |
| Host action | `plan` | not yet | not yet |
| Started by `flai serve` | on the operator's word ([ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)) | not yet | not yet |

## The planner

The planner plans one epic or one story and ends (S-0208, [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)). For an epic with no stories it drafts them; for a story it adds touches, a forecast, and a cost of delay value; for an epic with stories it revisits each one not done or cancelled and drafts what the outcome still lacks. It writes through flai alone and moves nothing past backlog.

### Starting it

The operator asks for it. Nothing starts it on a move or a timer.

| Way | For | What it runs |
|-----|-----|--------------|
| `flai plan <E-nnnn\|S-nnnn> [--json]` | a shell on the host | `serve.Plan` |
| Plan on an epic's or a story's page | the dashboard | the host API's `plan.run` `{id}`, which runs `flai plan <id> --json` |
| The MCP tool `plan` `{id}` | the operator's own agent | `serve.Plan`, as `flai plan` does; journalled as `mcp.plan` with the agent that asked |

Each needs the host action `plan` on for the project: `flai serve enable plan`, or the dashboard's Settings page while `settings` is on. It is off by default. A holder of the dashboard token can then start the planner for any epic or story that is open.

`serve.Plan` refuses, saying why:

- while the `plan` action is off for the project;
- a task, or an ID that is neither an epic's nor a story's;
- an item that is archived, done, or cancelled;
- while a planner runs for the item: one item has one planner at a time;
- when the planner's agent names no harness and no command is set on the host.

The MCP tool `plan` also refuses an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve`): planning is the operator's to ask for. `flai guard` refuses it to a sub-agent and to the planner.

### Its agent

`planning.agent` in `system-flow.yaml`, merged over the project's `agent` field by field, config key by config key, and role by role (`Manifest.PlanningAgent`), gives the planner's harness, model, config, and roles. With neither set, the host's command starts it. `flai check` reports a `planning.agent` that is not valid as `manifest.planning`. See [project-manifest.md](project-manifest.md).

### How it runs

`serve.Plan` starts the planner through the harness adapter with a `harness.Request` whose `Role` is `plan` and whose `Item` is the epic or story, and hands the process over: the command does not wait, and the serving flai settles the run.

- **Where.** In the project's main checkout, where the items are. No branch, no worktree.
- **As whom.** `FLAI_AGENT` is `planner-<item>`, such as `planner-E-0016`. The environment carries `FLAI_ROLE=plan` and `FLAI_ITEM=<item>`, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`, and no `FLAI_STORY`. `flai serve` drops `FLAI_STORY`, `FLAI_ROLE`, and `FLAI_ITEM` from its own environment before it adds a run's, so a planner started from a story's session is not taken for that story.
- **claude-code.** The session is named `<key> plan <item>`. It runs with `--agent planner` over the project's `.claude/agents/planner.md`, which is passed in `--agents` beside the explorer and the verifier, with the planner's agent's model when it names one. A project with no `planner.md` is refused, naming `flai upgrade`, which adds it from the template.
- **The operator's command.** It gets `FLAI_ROLE` and `FLAI_ITEM`, and `{story}` is empty.

### What it is told

`harness.Prompt` gives the planner a prompt of its own (`planPrompt`): plan the item and nothing else, as the convention's section "As the planner" says; prime with `--role plan` and the item; call `inbox`; read the item and what it links; hand wide search to the explorer. For an epic, draft its stories or revisit them; for a story, enrich it. Size stories, make each pass `flai check --strict`, and summarise the plan in one thread on the item. Write only through flai; never edit a file, move past backlog, finalize a draft, or overwrite the operator's inputs. When an input the operator owns is missing, ask with `thread_open` on the item and hold `wait_for_events` until it is answered. End with a one-line summary, which becomes the run's activity entry.

`planner.md` says the same in short and lists the planner's tools: `Read`, `Grep`, `Glob`, `Bash`, `Agent`, and flai's MCP tools for reading, `inbox`, `item_new`, `item_edit`, `item_move`, `thread_open`, `thread_reply`, `activity_log`, and `wait_for_events`. It lists no `Edit`, `Write`, or `NotebookEdit`.

### The guard

`flai guard` reads `FLAI_ROLE` from the hook's environment. With `plan`, it checks the session's own calls, which carry no agent ID, against the planner's rules ([flai-cli.md](flai-cli.md#commands), `flai guard`):

| Kind | Passes | Refused |
|------|--------|---------|
| flai MCP tools | the reads a sub-agent has; `inbox`, `item_new`, `item_edit`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`; `item_move` to `backlog` | every other tool, and `item_move` to any other state |
| flai commands | the reads; `story new`, `epic new`, `edit` without `--no-draft`, `touches`, `thread new` and `reply`, `issue new` and `bump`, `move` to `backlog` | every other command that writes; `edit --no-draft`, since finalizing is the operator's |
| git | the reads a sub-agent has | every other git command |
| File tools | | `Edit`, `Write`, `NotebookEdit` |

Each refusal names the rule broken and says to ask the operator on the item, or put it in the final summary. Other shell commands pass, as a sub-agent's do (ADR-0060). The planner's sub-agents are held as every sub-agent is.

Claude Code runs a hook only for the tools its matcher names. `.claude/settings.json` has two `PreToolUse` entries: `Bash|mcp__flai__.*` runs the guard in every session; `Edit|Write|NotebookEdit` runs it only when `FLAI_ROLE` is `plan`, and exits at once otherwise, so a story's session pays one shell test per edit. The template's file runs the installed `flai guard`; this repository's runs `scripts/flai.sh guard`.

### The run and how it ended

The run is recorded in `serve/agents.json` under `plans`, by item, the newest run for each, with the fields a story's run has and `item` in place of `story`. `Running` and `Last` never name a planner run, and the in-progress limit does not count it. Its output is in `serve/agents/<key>-planner-<start>.log`, beside the story agents' logs; each planner run gets a log of its own, at the first free second from its start.

When the process ends, or the next look finds it gone, flai serve settles the run (`planEnded`):

| Outcome | When |
|---------|------|
| `asked` | a thread on the item is open and its last entry is the planner's; `thread` names it |
| `failed` | it exited with a code other than 0, or the project or its threads could not be read |
| `worked` | otherwise, an exit nobody saw included |

It then logs the activity in `wip/agents/planner.md` with `serve.LogRunEnd`: the time since the last entry, the run's final reply as the summary, and the seconds and cost measured from the log ([metrics.md](metrics.md#strategic-agents-s-0206)). An `asked` planner is not started again when the thread is answered; the operator asks for another run.

Every start and failure is a journal entry with action `plan`. `agent.status` carries the runs as `plans`. When a planner run starts, cannot start, or ends, `flai serve` tells the dashboard with an `agent` notification that names the `item` in place of a story.

### On the dashboard

An epic's or a story's page shows **Plan** while the `plan` action is on, the dashboard can write, and the item is open: not done, cancelled, or archived. It says when the item's newest planner run has not ended. See [flaiover-dashboard.md](flaiover-dashboard.md).

## The orchestrator and the analyzer

Each primes with its role and has its activity document; `activity_log` takes its kind. No host action, guard rule, definition, or command starts either yet. When one does, it follows the planner's shape: a host action, a role in `FLAI_ROLE`, a definition in `.claude/agents/`, and guard rules of its own.
