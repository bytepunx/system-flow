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

The planner plans one epic or one story and ends (S-0208, [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)). For an epic with no stories it drafts them; for a story it adds touches, a forecast, and a cost of delay value, and drafts its tasks or revisits those it has (S-0255); for an epic with stories it revisits each one not done or cancelled, enriches it again, and proposes on the epic's plan thread what it would split, merge, add, or drop, drafting the additions (S-0209). It writes through flai alone, every story it creates is a draft, and it moves nothing past backlog.

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

`harness.Prompt` gives the planner a prompt of its own (`planPrompt`): plan the item and nothing else, as the convention's section "As the planner" says; prime with `--role plan` and the item; call `inbox`; read the item and what it links; hand wide search to the explorer. Size stories, and make each story or task pass `flai check --strict` and the markdown lint. Write only through flai; never edit a file, move past backlog, finalize a draft, or overwrite the operator's inputs. When an input the operator owns is missing, ask with `thread_open` on the item and hold `wait_for_events` until it is answered. End with a one-line summary, which becomes the run's activity entry.

What it does with the item depends on its kind (S-0209, S-0255):

| Item | Work | Thread | Summary |
|------|------|--------|---------|
| Epic with no stories | Draft the stories that deliver its outcome, each created with `draft: true` in the backlog: `item_new`'s `draft`, or `flai story new --draft` | One thread on the epic: the stories, their order (their `after`), and the assumptions made | Names the stories created |
| Epic with stories | Revisit each story not done or cancelled and enrich it again: touches, forecast, cost of delay value | The same one thread also proposes each story to split, merge, add, or drop; only the additions are created, as drafts. A finalized story is never cancelled or rewritten without asking | Names the stories created and the stories revisited |
| Story with no tasks | Enrich it: touches, a forecast, a cost of delay value. Draft the tasks that deliver its outcome, each with `## Work` and `## Done when`, a nature, tags, `touches`, and `after`, so that they form layers, created in the backlog: `item_new` with type task, or `flai task new` | One thread on the story: the tasks, their order and layers, and the assumptions made | Names by ID the tasks created |
| Story with tasks | Enrich it, then revisit each task not done or cancelled: its `touches` and `after` enriched again, the tasks the outcome lacks created | The same one thread proposes each task to split, merge, or drop. A task the planner did not write is never cancelled, nor its title, Work, or Done when rewritten, without asking | Names by ID the tasks created and the tasks revisited |

Tasks carry no topics, so a topic a task reaches goes on the story. `itemnew.Create` refuses a task that fails `flai check --strict` or the markdown lint and leaves nothing. The story's agent reviews the planner's tasks when it pulls the story ([work-management.md](../conventions/work-management.md)).

`planner.md` says the same in short and lists the planner's tools: `Read`, `Grep`, `Glob`, `Bash`, `Agent`, and flai's MCP tools for reading, `inbox`, `item_new`, `item_edit`, `item_move`, `thread_open`, `thread_reply`, `activity_log`, and `wait_for_events`. It lists no `Edit`, `Write`, or `NotebookEdit`.

### The guard

`flai guard` reads `FLAI_ROLE` from the hook's environment. With `plan`, it checks the session's own calls, which carry no agent ID, against the planner's rules ([flai-cli.md](flai-cli.md#commands), `flai guard`):

| Kind | Passes | Refused |
|------|--------|---------|
| flai MCP tools | the reads a sub-agent has; `inbox`, `item_new` (of a story only with `draft: true`), `item_edit`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`; `item_move` to `backlog` | every other tool; `item_new` of a story without `draft: true`; `item_move` to any other state |
| flai commands | the reads; `story new` with `--draft`, `epic new`, `task new`, `edit` without `--no-draft`, `touches`, `thread new` and `reply`, `issue new` and `bump`, `move` to `backlog` | every other command that writes; `story new` without `--draft` (S-0209), since the planner's stories are drafts for the operator to finalize; `edit --no-draft`, since finalizing is the operator's |
| git | the reads a sub-agent has | every other git command |
| File tools | | `Edit`, `Write`, `NotebookEdit` |

Each refusal names the rule broken and says to ask the operator on the item, or put it in the final summary. Other shell commands pass, as a sub-agent's do (ADR-0060). The planner's sub-agents are held as every sub-agent is. With `story new`, the last `--draft` decides: bare, or with a value that parses as true.

Claude Code runs a hook only for the tools its matcher names. `.claude/settings.json` has two `PreToolUse` entries: `Bash|mcp__flai__.*` runs the guard in every session; `Edit|Write|NotebookEdit` runs it only when `FLAI_ROLE` is `plan`, and exits at once otherwise, so a story's session pays one shell test per edit. The template's file runs the installed `flai guard`; this repository's runs `scripts/flai.sh guard`.

### Checked creation

A story the planner writes through `item_new` with a body is checked as `flai story new --body-stdin` is (S-0209): `itemnew.Create` writes it, runs `flai check` with it in place and the markdown lint on it, and a finding it introduces, such as a missing template section like `## Tasks` or `## Notes`, refuses it and leaves nothing. A creation with `after` was already checked this way (S-0176). An `item_new` with neither gets the template's own body and is not checked.

### The run and how it ended

The run is recorded in `serve/agents.json` under `plans`, by item, the newest run for each, with the fields a story's run has and `item` in place of `story`. `Running` and `Last` never name a planner run, and the in-progress limit does not count it. Its output is in `serve/agents/<key>-planner-<start>.log`, beside the story agents' logs; each planner run gets a log of its own, at the first free second from its start.

When the process ends, or the next look finds it gone, flai serve settles the run (`planEnded`):

| Outcome | When |
|---------|------|
| `asked` | a thread on the item is open and its last entry is the planner's; `thread` names it |
| `failed` | it exited with a code other than 0, or the project or its threads could not be read |
| `worked` | otherwise, an exit nobody saw included |

It then logs the activity in `wip/agents/planner.md` with `serve.LogRunEnd`: the time since the last entry, the run's final reply as the summary, the items it planned, and the seconds and cost measured from the log ([metrics.md](metrics.md#strategic-agents-s-0206)).

The items are named in this order (`plannedItems`, S-0209): the planned item; then the items under it created since the run started; then those under it created before and changed (`updated`) since. Under an epic are its stories and their tasks; under a story, its tasks. Archived items are left out, and each group is in ID order. When the items cannot be read, flai warns and the entry names the planned item alone. An `asked` planner is not started again when the thread is answered; the operator asks for another run.

The usage apportioned to the activity's span, the same share the entry's cost is, is also charged to the item the run was started for and to every item above it, under `usage.strategic` as the planner's entry, with the entry's seconds (S-0225, [ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md), [work-hierarchy.md](work-hierarchy.md)). An activity the planner logs with `activity_log` during the run is charged the same way, its own span's share. flai serve finds the item by the run's log: the item whose newest run under `plans` kept its log in the run being measured. A planner a person runs by hand, a run that a newer run for its item has replaced, and an activity outside any run charge nothing. The charge on the items and the entry in `planner.md` are two views of one spend; nothing adds them together. A charge that fails leaves the entry logged and is warned of. Measuring a story from its agents' logs keeps what the planner spent on it.

Every start and failure is a journal entry with action `plan`. `agent.status` carries the runs as `plans`. When a planner run starts, cannot start, or ends, `flai serve` tells the dashboard with an `agent` notification that names the `item` in place of a story.

### Measured on E-0016 (S-0209)

On 2026-10-04 the operator ran the planner on E-0016 with S-0209's prompt, over flai 1.30.0 on the host. The guard's draft rule, the checked `item_new`, and the entry's items were not yet in that flai.

| Measure | Result |
|---------|--------|
| Run | 04:34Z to 04:53Z; 1184 s of agent time, 4.49 USD (estimated), one activity entry |
| Revisited | All 22 open stories; touches and a forecast on the 19 not in progress (S-0210 to S-0224, S-0226 to S-0229); the three in progress left alone |
| Asked | Cost of delay inputs and the hour rate, on TH-0106; then a value for each story, its share of the epic's 1500 USD a week by forecast duration |
| Plan thread | TH-0107: the order its forecasts assume in two lanes, its assumptions, its risks, and four proposals: rewrite S-0226 and S-0227, make S-0213 wait for S-0217 and S-0223 for S-0211, split S-0228 |
| Applied on the operator's word | All four; the split created S-0259 as a draft in the backlog, with every section and its `after` |
| Created otherwise | None: it judged the open stories to cover the outcome |
| Found | `item_edit` refuses a touch that starts with a dot ([I-0071](../issues/I-0071-item-edit-refuses-a-touch-that-starts-with-a-dot-though-flai-touches-accepts-it.md)); it used `flai touches` instead |

Its first count of the open stories was wrong (20 and 17); it corrected itself on the same thread. Its summary, the entry's, named the story it created and the stories it revisited.

The designer's judgement, on TH-0103: "I found this useful. I am looking forward to using the planner and refining it over time."

An epic with no stories was not measured: the only open epics, E-0015 and E-0016, both had stories.

### Measured on S-0217 (S-0255)

On 2026-10-04 the operator ran the planner on S-0217, a backlog story with no tasks, with S-0255's prompt, over flai 1.30.0 on the host. The guard's `task new`, the checked `item_new`, and the entry's items were not yet in that flai; tests cover them.

| Measure | Result |
|---------|--------|
| Run | Ended 19:10Z; 333 s of agent time, 1.80 USD (estimated), one activity entry |
| Created | Six tasks in the backlog, T-0809 to T-0814, in four layers: T-0809 and T-0810, then T-0811 and T-0812, then T-0813, then T-0814. Each has `## Work` saying what it changes and why it waits, `## Done when` with its tests, a nature, tags, file-level `touches`, and `after`; no two tasks of a layer share a path. `flai check --strict` reports nothing on them |
| Story | Touches widened by the check package one task needs; forecast duration kept at 2h and delivery moved for the stories ahead of it; cost of delay value kept |
| Plan thread | TH-0108: the tasks and layers as a table, nine assumptions to confirm, no split, merge, or drop, and two notes for later stories (S-0213, S-0222) |
| Summary | The entry's summary named the six tasks created, said none were revisited, and named TH-0108 |

The designer resolved TH-0108 and TH-0102 without asking for a change: the drafts stand as written.

Revisiting a story's open tasks was not measured: the second run on S-0217 was not made. The prompt and the convention cover it.

### On the dashboard

An epic's or a story's page shows **Plan** while the `plan` action is on, the dashboard can write, and the item is open: not done, cancelled, or archived. It says when the item's newest planner run has not ended. See [flaiover-dashboard.md](flaiover-dashboard.md).

## The orchestrator and the analyzer

Each primes with its role and has its activity document; `activity_log` takes its kind. No host action, guard rule, definition, or command starts either yet. When one does, it follows the planner's shape: a host action, a role in `FLAI_ROLE`, a definition in `.claude/agents/`, and guard rules of its own.
