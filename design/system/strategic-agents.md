---
title: Strategic agents
updated: 2026-10-06
status: active
topics: [planning, orchestration, analysis]
---

# Strategic agents

The planner, the orchestrator, and the analyzer are agents that work above a story (E-0016). What each does and never does is the convention [strategic-agents.md](../conventions/strategic-agents.md). This document is how flai runs them. The planner runs when the operator asks for it and, behind the `plan` host action, on its own. The orchestrator runs, one per project, for as long as the `orchestrate` host action is on. The analyzer runs, one per project at a time, behind the `analyze` host action, when the operator asks for it and on a schedule, and writes one report under `design/analysis/`.

| Part | Planner | Orchestrator | Analyzer |
|------|---------|--------------|----------|
| Pack: `flai prime --role` ([ADR-0075](../adrs/0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md)) | `plan`, with `--epic` or `--story` | `orchestrate` | `analyze` |
| Activity document ([ADR-0079](../adrs/0079-the-planner-the-orchestrator-and-the-analyzer-each-log-their-activities-in-one.md)) | `wip/agents/planner.md` | `wip/agents/orchestrator.md` | `wip/agents/analyzer.md` |
| Host action | `plan` | `orchestrate` | `analyze` |
| Started by `flai serve` | on the operator's word ([ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)), on the orchestrator's for an epic behind `plan_backlog_epics` (S-0219), and on its own behind `plan` ([ADR-0084](../adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md), [Planning again](#planning-again)) | while `orchestrate` is on, one per project, started again when it ends ([ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md), [The orchestrator](#the-orchestrator)) | on the operator's word, and on `analysis.schedule` while `analyze` is on, one per project at a time ([ADR-0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md), [The analyzer](#the-analyzer)) |

## The planner

The planner plans one epic or one story and ends (S-0208, [ADR-0082](../adrs/0082-flai-serve-starts-the-planner-for-an-epic-or-a-story-behind-the-plan-host.md)). For an epic with no stories it drafts them; for a story it adds touches, a forecast, and a cost of delay value ([Enriching a story](#enriching-a-story-s-0210)), and drafts its tasks or revisits those it has (S-0255); for an epic with stories it revisits each one not done or cancelled, enriches it again, and proposes on the epic's plan thread what it would split, merge, add, or drop, drafting the additions (S-0209). It writes through flai alone, every story it creates is a draft, and it moves nothing past backlog.

### Starting it

The operator asks for it, in one of the ways below. While the `plan` host action is on, `flai serve` also starts it on its own: on an edit, when work ahead completes or the order changes, and on a schedule ([Planning again](#planning-again)).

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

The MCP tool `plan` also refuses an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve`): planning is the operator's to ask for. The orchestrator is the one exception: it may ask for an epic `flai plan --candidates` lists while the operator gives it `plan_backlog_epics` ([What it does with each permission](#what-it-does-with-each-permission-s-0219)). `flai guard` refuses the tool to a sub-agent, to the planner, and to the analyzer.

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

### Enriching a story (S-0210)

To enrich a story, the planner runs three commands, reviews what they print, and writes what it decides. The commands write nothing. Each prints its result, and `--json` gives the structure. The guard lets every sub-agent and the planner run all three (`cliReads`).

| Command | Prints | Code |
|---------|--------|------|
| `flai touches suggest <S-nnnn> [path...] [--min 2] [--limit 20]` | The files most often changed together with the story's touches | `cmd/touches.go`, `storygit.CommitFiles`, `planning.CoChanged` |
| `flai forecast <S-nnnn>` | A duration, a delivery, and a one-sentence basis | `cmd/forecast.go`, `planning.Forecast` |
| `flai cod <E-nnnn\|S-nnnn>` | A cost of delay per week, and its basis | `cmd/cod.go`, `planning.CostOfDelay` |

#### Touches from co-change

The seeds are the story's touches, its own and those of its tasks not cancelled, read as a claim reads them (`Holds.Claim`: a sub-project's name or tag means its path), and the paths given. A story with no touches and no path given is refused.

`CommitFiles` reads the main branch's history once, in the main checkout: commits that are not merges, renames not followed, paths from the project's root. Story branches are not read, so a commit is not counted on its branch and again after acceptance. Files under the wip folder are left out, since flai writes them. So are flai's own bookkeeping commits, whose files would otherwise seem to change with everything:

| Left out | Subject |
|----------|---------|
| Acceptance | `chore: [ID] accept and archive`, with an older release in it or not |
| Creation | `chore: [ID] create …` |
| Edit by fields | `chore: [ID] edit …` |
| Release | `chore: publish …` |
| Default agent | `chore: set the project's default agent`, `chore: clear the project's default agent` |
| Import | `chore: bring … under system-flow (flai import)` |

A commit with no file left is dropped. A seed commit is one that changed a file under a seed: the file itself, or a folder that holds it. Every other file is counted by the seed commits it is in, with that count's share of them. Files counted at least `--min` times are kept, ordered by count, most first, then by path, and cut at `--limit` (0 keeps every one). The header names the seeds and how many of the commits are seed commits.

#### Forecast

A story's size is its acceptance criteria's checkboxes and its own touches counted together, at least 1. Its band is small up to 6, medium up to 12, and large from 13.

The history is the done stories with `usage.seconds` above zero and a cycle time, from first in progress to done. Each gives three ratios: agent seconds per unit of size; seconds in progress, over every spell, per agent second; and cycle time per agent second. The match takes the first rung with at least three stories:

| Rung | The done stories that share the story's |
|------|------------------------------------------|
| 1 | Nature, model (`agent.model`), and band |
| 2 | Nature and band |
| 3 | Nature |
| 4 | Nothing: all of them |

| Figure | From the matched stories |
|--------|--------------------------|
| Duration | The median seconds per unit times the size, rounded up to the minute |
| Busy factor | The median seconds in progress per agent second, at least 1: how long a story holds its place in progress |
| Cycle factor | The median cycle time per agent second, at least 1: how long from the start of progress to done |

With fewer than three stories on every rung, the duration is `planning.default_duration` (1h unset, [project-manifest.md](project-manifest.md)), both factors are 1, and the basis says so.

The delivery plays out the board from now. The lanes are the board's in-progress limit; with none, every story starts as soon as what it waits for allows. The stories are played out in this order:

| Stories | Played out |
|---------|------------|
| The story itself, in progress or in review | Delivered at its start plus its duration times the cycle factor, never before now. Nothing else is played out |
| In progress | Each holds a lane until its start plus its duration times the busy factor, and is delivered at its start plus its duration times the cycle factor, neither before now. More than the limit free a lane only once enough have finished to bring the count under it |
| In review | Not played out: they hold no lane |
| Ready, then backlog, in pull order | Each starts at the earliest free lane, but not before the delivery of any story in its `after` already played out; holds the lane for its duration times the busy factor; and is delivered at its start plus its duration times the cycle factor. The play-out stops at the story |

A story ahead with its own `forecast.duration` is played out with it; the others, and the story itself, with the duration worked out from history. Times are UTC, rounded up to the minute. `--json` adds `size`, `criteria`, `touches`, `default`, `history` (`match`, `stories`, `seconds_per_unit`, `busy_factor`, `cycle_factor`), `position` (its place among the ready and backlog stories, 0 in progress or in review), `limit`, and `ahead`, each story played out before it with its duration, `from_forecast`, start, and delivery.

#### Cost of delay

From an epic's or a story's own inputs, the value per week is:

`revenue_per_week` + `penalty_per_week` + hours of `time_lost_per_cycle` × `planning.hour_rate` × (168h ÷ `planning.cycle`)

Time lost with no `planning.hour_rate` is refused. E-0016's 10h lost per 168h cycle at 150 USD an hour is 1500 USD a week.

A story without inputs takes a share of its epic's value: the epic's worked out from its inputs, else the value recorded on it. The share is the story's duration over the sum of the durations of the epic's open stories without inputs, each its own `forecast.duration`, else the one `flai forecast` works out. On 2026-10-04 S-0210's was 2h of 30h45m over 20 stories, 97.56 USD a week.

Refused: an epic without inputs; a story without inputs that has no epic, or whose epic has neither inputs nor a value; a task; an item closed or archived. Each refusal names the `flai edit` flags that give the inputs, which are the operator's. Amounts are rounded to two decimals, in `planning.currency`. `--json` adds `from` (`inputs` or `epic`), the `inputs` with the hour rate, cycle, and cycles a week when time lost is among them, and the `epic` share with the stories it was apportioned over.

#### What the planner writes

The planner predicts the touches from what `touches suggest` lists, the story's goal and criteria, the design it links, and the code layout, and keeps every touch the story declares. It reviews the forecast and the value, adjusts a figure only with a stated reason, and writes the touches, the forecast's duration, delivery, and basis, and the cost of delay value through flai: `item_edit`, or `flai edit` and `flai touches`, which stamp `by` and `at` ([ADR-0074](../adrs/0074-work-items-carry-planning-data-a-story-s-draft-flag-an-epic-s-or-story-s-cost.md)). It never changes the operator's inputs.

It records why under a `### Planning` heading in the story's Notes: where each touch came from (declared, co-change, design, or layout) and why each figure stands or was adjusted. That heading is the planner's, rewritten on each run; the rest of the Notes is left as it was. `planPrompt`, the convention's "As the planner", and `planner.md` say the same.

### The guard

`flai guard` reads `FLAI_ROLE` from the hook's environment. With `plan`, it checks the session's own calls, which carry no agent ID, against the planner's rules ([flai-cli.md](flai-cli.md#commands), `flai guard`):

| Kind | Passes | Refused |
|------|--------|---------|
| flai MCP tools | the reads a sub-agent has; `inbox`, `item_new` (of a story only with `draft: true`), `item_edit`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`; `item_move` to `backlog` | every other tool; `item_new` of a story without `draft: true`; `item_move` to any other state |
| flai commands | the reads, `forecast`, `cod`, and `touches suggest` among them; `story new` with `--draft`, `epic new`, `task new`, `edit` without `--no-draft`, `touches`, `thread new` and `reply`, `issue new` and `bump`, `move` to `backlog` | every other command that writes; `story new` without `--draft` (S-0209), since the planner's stories are drafts for the operator to finalize; `edit --no-draft`, since finalizing is the operator's |
| git | the reads a sub-agent has | every other git command |
| File tools | | `Edit`, `Write`, `NotebookEdit` |

Each refusal names the rule broken and says to ask the operator on the item, or put it in the final summary. Other shell commands pass, as a sub-agent's do (ADR-0060). The planner's sub-agents are held as every sub-agent is. With `story new`, the last `--draft` decides: bare, or with a value that parses as true.

Claude Code runs a hook only for the tools its matcher names. `.claude/settings.json` has two `PreToolUse` entries: `Bash|mcp__flai__.*` runs the guard in every session; `Edit|Write|NotebookEdit` runs it only when `FLAI_ROLE` is `plan`, `orchestrate` (S-0218, [Its permissions and the guard](#its-permissions-and-the-guard)), or `analyze` (S-0223, [The analyzer's guard](#the-analyzers-guard)), and exits at once otherwise, so a story's session pays one shell test per edit. The template's file runs the installed `flai guard`; this repository's runs `scripts/flai.sh guard`.

### Checked creation

A story the planner writes through `item_new` with a body is checked as `flai story new --body-stdin` is (S-0209): `itemnew.Create` writes it, runs `flai check` with it in place and the markdown lint on it, and a finding it introduces, such as a missing template section like `## Tasks` or `## Notes`, refuses it and leaves nothing. A creation with `after` was already checked this way (S-0176). An `item_new` with neither gets the template's own body and is not checked.

### The run and how it ended

The run is recorded in `serve/agents.json` under `plans`, by item, the newest run for each, with the fields a story's run has, `item` in place of `story`, and `trigger`, what started it: `asked` when the operator asked, `orchestrator` when the orchestrator did (S-0219), otherwise the replanner's triggers ([What a run records](#what-a-run-records)). `Running` and `Last` never name a planner run, and the in-progress limit does not count it. Its output is in `serve/agents/<key>-planner-<start>.log`, beside the story agents' logs; each planner run gets a log of its own, at the first free second from its start.

When the process ends, or the next look finds it gone, flai serve settles the run (`planEnded`):

| Outcome | When |
|---------|------|
| `asked` | a thread on the item is open and its last entry is the planner's; `thread` names it |
| `failed` | it exited with a code other than 0, or the project or its threads could not be read |
| `worked` | otherwise, an exit nobody saw included |

It then logs the activity in `wip/agents/planner.md` with `serve.LogRunEnd`: the time since the last entry, the run's final reply as the summary, the run's trigger as the entry's `- Trigger:` line ([agent-narrative.md](agent-narrative.md#strategic-agents-activity-documents)), the items it planned, and the seconds and cost measured from the log ([metrics.md](metrics.md#strategic-agents-s-0206)).

The items are named in this order (`plannedItems`, S-0209): the planned item; then the items under it created since the run started; then those under it created before and changed (`updated`) since. Under an epic are its stories and their tasks; under a story, its tasks. Archived items are left out, and each group is in ID order. When the items cannot be read, flai warns and the entry names the planned item alone. An `asked` planner is not started again when the thread is answered; the operator asks for another run.

The usage apportioned to the activity's span, the same share the entry's cost is, is also charged to the item the run was started for and to every item above it, under `usage.strategic` as the planner's entry, with the entry's seconds (S-0225, [ADR-0083](../adrs/0083-a-planner-activity-s-usage-is-charged-to-the-item-it-planned-and-the-items.md), [work-hierarchy.md](work-hierarchy.md)). An activity the planner logs with `activity_log` during the run is charged the same way, its own span's share. flai serve finds the item by the run's log: the item whose newest run under `plans` kept its log in the run being measured. A planner a person runs by hand, a run that a newer run for its item has replaced, and an activity outside any run charge nothing. The charge on the items and the entry in `planner.md` are two views of one spend; nothing adds them together. A charge that fails leaves the entry logged and is warned of. Measuring a story from its agents' logs keeps what the planner spent on it.

Every start and failure is a journal entry with action `plan`. `agent.status` carries the runs as `plans`. When a planner run starts, cannot start, or ends, `flai serve` tells the dashboard with an `agent` notification that names the `item` in place of a story.

### Planning again

While the `plan` host action is on, `flai serve` plans again without the operator asking (S-0211, [ADR-0084](../adrs/0084-flai-serve-plans-again-on-its-own-behind-the-plan-host-action-on-an-edit-when.md)). The operator's own asking is unchanged. A replanner per served project (`internal/serve/replan.go`) looks whenever the project's launcher does: when its work items or threads change, when an agent ends, and every minute. It acts on three triggers.

| Trigger | When | What `flai serve` does |
|---------|------|------------------------|
| An edit | Someone other than a planner edits the goal, criteria, or touches of a backlog or ready story whose forecast or cost of delay value is older than the edit; or edits its cost of delay inputs, leaving them newer than its value (`CostOfDelay.Stale`) | Queues the planner for the story |
| Work ahead completes or the order changes | A story is accepted or cancelled, or the pull order changes | What `planning.replan` says ([The replan policy](#the-replan-policy)) |
| The schedule | `planning.schedule` comes round | Queues the planner for every ready story, in pull order |

Edits are read from the edit notices in `.flai-cache/edits.jsonl`, which `flai edit`, `flai touches`, the dashboard, and the MCP tool `item_edit` leave ([flai-cli.md](flai-cli.md#commands)); a hand edit of a file leaves none. An edit by an agent named `planner-…` never triggers, so a planner's own writes do not start it again. A story never planned, with neither a forecast nor a cost of delay value, is not planned on an edit: replanning keeps planning fresh, and planning a story first is the operator's to ask. A cost of delay value written alone triggers nothing.

#### The replan policy

`planning.replan` in `system-flow.yaml` ([project-manifest.md](project-manifest.md)) says what happens when a story is accepted or cancelled or the pull order changes:

| Value | What happens |
|-------|--------------|
| `never` | Nothing |
| `deterministic`, the default | Every ready and backlog story with a forecast duration is played out again, in pull order, with `planning.Replay`: the [forecast](#forecast)'s play-out, keeping the story's own duration. Where the delivery moved, the new delivery and basis are written; the duration stays. No agent runs |
| `agent` | The same, then the planner is queued for each story written, with the triggers |

The replay writes each story through `itemedit.Apply`, the path `flai edit` takes, stamped `by: flai`, with no commit per story. The files written, the items and `wip/agents/index.md`, are then committed together in the main checkout as `chore: replan forecasts after <triggers>`, unless `dashboard.autocommit` is false. flai serve logs `forecasts replanned` with the policy, how many moved, and the triggers. A replay or a write that fails is warned of, and the others go on.

A known cost: `itemedit.Apply` runs `flai check` before and after each write, as `flai edit` does, so a replan that moves many forecasts takes about a second per story it writes.

#### The schedule

`planning.schedule` is a five-field cron expression in UTC (minute, hour, day of month, month, day of week), such as `0 6 * * 1-5`, or `daily`, which is 00:00 UTC. flai parses it itself (`internal/cron`, no library): `*`, values, ranges, comma lists, and steps `/n`; Sunday is 0 or 7; with both day fields given, a day either names matches. An expression that never comes round, such as February 31st, is refused. `flai check` reports one it cannot parse as `manifest.planning`.

The replanner sees the schedule come round at the next look after its time, so at most a minute late. A schedule set or changed comes round first at its next time from then. Times missed between two looks come round once.

#### The queue

The queue holds a story once: a trigger for a story already queued is added to its entry, not queued again. flai serve drains it one planner at a time per project, with its own launcher, so it sees the run end at once and starts the next. It passes over an entry whose item has a planner running, such as one the operator asked for, and comes back to it. Each entry goes through `planCheck`, the checks `serve.Plan` makes ([Starting it](#starting-it)); an entry they refuse, such as a story accepted meanwhile, is dropped and logged at info as `queued planner dropped`, with the reason. A run the operator asks for is not queued.

#### What a run records

A run's triggers, joined by semicolons, are its `trigger` in `serve/agents.json` and the `- Trigger:` line of its activity entry in `wip/agents/planner.md` ([agent-narrative.md](agent-narrative.md#strategic-agents-activity-documents)), such as `- Trigger: accepted S-0210; reordered`.

| Trigger | Said as |
|---------|---------|
| The operator asked | `asked` |
| The orchestrator asked (S-0219) | `orchestrator` |
| An edit | `edited <fields> by <who>`, the fields among `goal`, `criteria`, `touches`, and `cost_of_delay` |
| A story accepted | `accepted <ID>` |
| A story cancelled | `cancelled <ID>` |
| The pull order changed | `reordered` |
| The schedule | `schedule <expression>`, as written in the manifest |

#### What it keeps

What the replanner has seen is kept in memory, from flai serve's start: an edit, a move, or a scheduled time that passed while flai serve was down is not acted on. While `plan` is off it acts on nothing and drops its queue, but keeps up with what happens, so turning `plan` on does not act on the past. A `planning.replan` or `planning.schedule` that is not valid is warned of once in flai serve's log, and nothing is done for it until it is fixed.

#### On the settings page

The Settings page shows the triggers read-only, under "Planning again, for this project" (`planningTriggers` in `cmd/serve_actions.go`): whether edits to a planned story start the planner, which is whether `plan` is on; the replan policy, marked when it is the default; and the schedule with its next run in UTC. A value flai cannot read is shown with its error. The operator sets both keys in `system-flow.yaml` by hand, as the other `planning` keys.

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

## The orchestrator

The orchestrator keeps a project's work moving, and nothing else (S-0218, [ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)). It asks the planner to plan an epic, finalizes drafts, promotes stories to ready, orders the ready column, answers threads, accepts stories, and publishes releases, each only while the operator's permission for it is on. It takes every figure from flai's commands, logs each decision, and waits on events between decisions. Unlike the planner, it has no item and does not end: `flai serve` runs one per project for as long as the operator wants it.

### Starting and stopping it

The host action `orchestrate`, off by default, turns it on for a project: `flai serve enable orchestrate`, or the dashboard's Settings page while `settings` is on. No command or method starts it. `flai serve` reads the action at every look of the project's launcher: when its work items or threads change, when an agent ends, and every minute (`internal/serve/orchestrate.go`).

| At a look | What `flai serve` does |
|-----------|------------------------|
| The action is on and no run is going | Starts the orchestrator |
| The action is on and the last run failed less than a minute ago | Nothing until the minute has passed, so that a run that fails at once does not spin |
| The action is on and a run is going | Nothing |
| The action is off and a run is going | Stops it as `flai serve agent stop` stops a story's agent: marks it stopped, ends its process group, and kills the group once the grace has passed. The stop is journalled, and the run's activity is logged with the summary `stopped: orchestrate turned off` |

A run that ends while the action is on is started again at the next look. The orchestrator does not end by itself, so a run ends when it fails, when it is stopped, or when its harness ends it.

A start is refused when the orchestrator's agent names no harness and no command is set on the host, and when the harness refuses it, as `claude-code` does in a project with no `.claude/agents/orchestrator.md`, naming `flai upgrade`. A refusal is recorded as a failed run and journalled once: the same refusal at the next looks is not recorded again until a start gets past it or the action is turned off.

### Its agent

`orchestration.agent` in `system-flow.yaml`, merged over the project's `agent` as `planning.agent` is (`Manifest.OrchestrationAgent`), gives the orchestrator's harness, model, config, and roles. With neither set, the host's command starts it. `flai check` reports an `orchestration.agent` that is not valid as `manifest.orchestration`. See [project-manifest.md](project-manifest.md).

### How it runs

The launcher's `orchestrate` starts it through the harness adapter with a `harness.Request` whose `Role` is `orchestrate`, with no `Item` and no `Story`.

- **Where.** In the project's main checkout, where the items are. No branch, no worktree.
- **As whom.** `FLAI_AGENT` is `orchestrator`. The environment carries `FLAI_ROLE=orchestrate`, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`, and no `FLAI_ITEM` or `FLAI_STORY`. `roleEnv` refuses an orchestrator request that names an item or a story.
- **claude-code.** The session is named `<key> orchestrate`. It runs with `--agent orchestrator` over the project's `.claude/agents/orchestrator.md`, which is passed in `--agents` beside the explorer and the verifier, with the orchestrator's agent's model when it names one. Its permission handler is flai's `permission_prompt`, as for every session `flai serve` starts ([ADR-0086](../adrs/0086-flai-serve-gives-a-claude-code-agent-flai-s-permission-prompt-as-its-permission.md)).
- **The operator's command.** It gets `FLAI_ROLE`, and `{story}` is empty.

### What it is told

`harness.Prompt` gives the orchestrator a prompt of its own (`orchestratePrompt`): keep the project's work moving and do nothing else, as the convention's section "As the orchestrator" says; prime with `--role orchestrate`; call `inbox` and read the board; hand wide search to the explorer. Act only within `orchestration.permissions` and by `orchestration.policy`, and ask when a permission is unclear. Take every figure from flai and do no arithmetic: the epics to plan from `flai plan --candidates`, the drafts complete enough to finalize from `flai promote --drafts`, the order from `flai order --by` (`order_by_policy`), the stories that could go to ready from `flai promote --candidates` (`promote_candidates`), and whether a release is due from `flai release --evaluate` (`release_evaluate`). With each permission, run its command and act on what it prints ([What it does with each permission](#what-it-does-with-each-permission-s-0219)). Reply on the threads awaiting the operator as `answer_threads` says, read afresh from `system-flow.yaml` each time ([Answering threads](#answering-threads-s-0220)). While `publish` is on, evaluate the release policy after each acceptance and publish through `release_publish` when it is met ([Publishing a release](#publishing-a-release-s-0222)). Log each action with `activity_log`, kind `orchestrator`: what it did, on which items, why, and the policy figure behind it, as flai gave it. Then hold `wait_for_events`, again each time it returns, and decide again when something has changed, without ending. Write only through flai; never edit code or documents, and never work a story. Never work around a refusal: a refusal from `flai guard` or from flai ends that attempt, which it logs and does not try again until something it depends on changes. When a decision needs the operator, ask with `thread_open` on the item, the recommended answer first, and do what needs no answer meanwhile. The convention's "As the orchestrator" says the same (template 1.0.53, its threads 1.0.55).

`orchestrator.md` says the same in short and lists its tools: `Read`, `Grep`, `Glob`, `Bash`, `Agent`, and flai's MCP tools for reading, `inbox`, `board`, `item_edit`, `item_move`, `thread_open`, `thread_reply`, `activity_log`, `wait_for_events`, `plan`, `order_by_policy`, `promote_candidates`, `release_evaluate`, and, since S-0222, `release_publish`. It lists no `Edit`, `Write`, or `NotebookEdit`. It finalizes a draft with `item_edit` giving only the story's `id` and `draft: false`.

### What it does with each permission (S-0219)

Under each of four permissions the orchestrator runs one command that reads, judges what it prints, and makes one kind of write. The reads write nothing and pass the guard whatever the permissions. flai holds the writes itself, as well as the guard: under `FLAI_ROLE=orchestrate`, `flai plan`, `flai edit --no-draft`, `flai move … ready`, and `flai order --by --apply`, and the MCP tools `plan`, `item_edit` with `draft: false`, and `item_move` to `ready`, check the permission and the item before they write (`Repo.OrchestratorPermits`), so that a call the guard does not see is held all the same.

| Permission | It reads | It writes | flai also refuses |
|------------|----------|-----------|-------------------|
| `plan_backlog_epics` | `flai plan --candidates`: each backlog epic with no stories, and each open epic whose stories are all done or cancelled with one done ([flai-cli.md](flai-cli.md#commands)) | The planner for each epic listed, one at a time, with the MCP tool `plan` | An ID that is not an epic's, and an epic the candidates do not list, naming why: left out while a planner runs for it or its question awaits the operator, closed, archived, or neither kind of candidate (`orchestratorPlans`) |
| `finalize_drafts` | `flai promote --drafts`: each backlog draft, complete or with what it lacks | A complete draft it judges consistent, finalized with `item_edit` giving only its `id` and `draft: false`, or `flai edit <S-nnnn> --no-draft` alone. For any other draft it opens one thread on the story saying what is missing or inconsistent, once, and leaves it | A change to anything but the draft flag; a draft that is not complete, naming what it lacks: `S-0300 is not complete, so the orchestrator does not finalize it (flai promote --drafts): no touches; no forecast duration` |
| `promote_to_ready` | `flai promote --candidates` (`promote_candidates`) | Its candidates, moved to `ready` in its order with `item_move`, while the ready column's limit has room | An item that is not a story; a story that is not a candidate, with each reason the candidates give it, such as `draft: finalize it first`, a hold, or `no forecast duration`; a story not in the backlog; any story while the ready column is at its WIP limit: `the ready column is at its WIP limit (5 of 5): the orchestrator moves no story to ready until one leaves it` (`Repo.Promotable`) |
| `order_ready` | `flai order --by <policy>` (`order_by_policy`), with the policy `orchestration.policy` names | The ready column's order, `flai order --by <policy> --apply`, after each change to the ready column, its own or another's | Nothing besides the permission. The order keeps in place a ready story placed by hand in the last day ([The pull order](workflow.md#the-pull-order-s-0057), [ADR-0088](../adrs/0088-board-md-records-who-placed-a-story-by-hand-and-when-and-a-policy-s-order-keeps.md)) |

With a permission off, flai's refusal is `the orchestrator <does what> only with orchestration.permissions.<name>, which is off: ask the operator with thread_open on <item>`, after `rule:` on the command line. A finalized draft's `finalized` block names the orchestrator, as any edit's names who made it. `flai promote --drafts` judges what arithmetic can: every section of the story template, a goal, a checkbox criterion, a touch, a forecast duration and delivery, a cost of delay value, and an open epic or none. Whether the criteria, touches, and forecast describe the same work is the orchestrator's judgement.

A planner the orchestrator starts runs as any planner does, behind the `plan` host action and `serve.Plan`'s refusals ([Starting it](#starting-it)), and is started with `serve.PlanForOrchestrator`: its run records the trigger `orchestrator` in place of `asked`, so `serve/agents.json` and the `- Trigger:` line of its entry in `wip/agents/planner.md` say who asked. The journal names `orchestrator` as who asked. The orchestrator never places a story by hand: the guard refuses it `flai order <story>` whatever its permissions, and under `FLAI_ROLE=orchestrate` `flai order` records any placement as the orchestrator's, which the window does not keep.

Each action is one `activity_log` entry in `wip/agents/orchestrator.md`: what it did, the items, why, and the figure flai gave behind it, such as a story's cost of delay value, its value over its duration, its forecast, or a candidate's rank. A refusal ends the attempt: the orchestrator logs it with the refusal and does not retry until something it depends on changes. The guard's refusals are also appended under `## Refusals` ([Its permissions and the guard](#its-permissions-and-the-guard)); flai's own are in the entry the orchestrator logs. `cmd/orchestrate_permissions_test.go` runs each of the four permissions on and off on a fixture board, through the guard and flai.

### Answering threads (S-0220)

`orchestration.permissions.answer_threads` is `off`, `recommend`, or `autonomous`, and lets the orchestrator reply on the threads awaiting the operator: open, last entry by a story's agent, not opened by the orchestrator, and with no recommendation pending (`pending_recommendation` in `inbox`). The operator may change it while the orchestrator runs, so its prompt states the rules for each value and tells it to read the value from `system-flow.yaml` each time before it acts on threads (`orchestratePrompt`).

| Value | It replies |
|-------|------------|
| `off` | Not at all |
| `recommend` | With a recommendation on each, with `thread_reply`, `recommendation: true`, and a `source`: the ADR, design section, or convention its answer rests on, as `<path>` or `<path>#<heading>`, read with `doc_get` first. The thread stays awaiting the operator, who confirms the recommendation (`flai thread confirm`, or the dashboard) or replies in its place ([ADR-0090](../adrs/0090-a-thread-entry-is-marked-a-recommendation-in-its-heading-and-cites-its-source.md)) |
| `autonomous` | With an answer and a `source` when a source settles the question, which ends the wait as the operator's answer would. With a recommendation instead, escalating to the operator, when no source settles it or the question asks for the operator's judgement: a decision not yet recorded, a change of scope, or money, such as a cost of delay input, an estimate, or spend |

flai enforces only what it can tell. The guard holds its thread calls ([Its permissions and the guard](#its-permissions-and-the-guard)): with `answer_threads` off it refuses any reply on another's thread; with `recommend` it lets only a recommendation through; with `autonomous` it lets an answer through only with a source that names a file, and tells it to cite one or post a recommendation. Whether a source settles the question, and whether the question is the operator's judgement, is the orchestrator's to decide. A recommendation ends no wait, and `flai stats` counts the waits the orchestrator ended apart from the operator's ([ADR-0091](../adrs/0091-a-recommendation-ends-no-thread-wait-and-flai-stats-reports-the-waits-the.md)).

Each reply is logged in its decision log by `thread_reply` itself, not by the orchestrator calling `activity_log`: an entry in `wip/agents/orchestrator.md` saying `Answered TH-nnnn` or `Recommended an answer on TH-nnnn`, `citing <path> § <heading>` or `citing no source`, with the thread's story as its item. A reply posted whose log entry fails is reported as an error naming both, so the orchestrator logs it with `activity_log`.

### Accepting a story (S-0221)

With `orchestration.permissions.accept_reviews` on, the orchestrator accepts a story in review through `flai accept <S-nnnn> --by orchestrator`, the same flow as the operator's: merge, done, archive, commit, and never a release ([ADR-0093](../adrs/0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md), refining [ADR-0032](../adrs/0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md)). For each story in review it:

1. Hands the story's worktree to its verifier, which runs the close-out's tests, lint, and `flai check --strict`, and matches each acceptance criterion to the files the branch changes that meet it.
2. Runs `flai accept <S-nnnn> --by orchestrator --verified <commit> --dry-run`, with the commit the verifier checked, and reads the blockers.
3. With no blocker and every criterion matched, runs `flai accept <S-nnnn> --by orchestrator --verified <commit> --evidence <file>`, and logs the acceptance with `activity_log`, naming the story and the commit.
4. Otherwise leaves the story in review, opens a thread on it saying what is missing, each blocker and each criterion it could not check against the diff, and logs that decision.

The acceptance preview adds four blockers for an acceptance by the orchestrator, beside those it reports for any acceptance, and `flai accept` refuses on any blocker before anything is merged:

| Condition | Blocker when |
|-----------|--------------|
| The verifier passed at the branch head | `--verified` is missing, or names a commit that is not the story branch's head |
| Every criterion is ticked | A criterion under `## Acceptance criteria` is unticked; each is named |
| The diff stays within the touches | A file the branch changes is under none of the story's touches, read as a claim reads them; each is named. The workflow's own records are the story's whatever it touches: the `wip` folder, the issues under `design/issues`, and an experiment's results |
| No thread is open | A thread on the story or one of its tasks is not resolved; each is named |

The evidence, `--evidence <file>` or `-` for standard input, is markdown: a `Verdict:` line from the verifier's report, and one list item per criterion, `- <n>: <files>`, naming the changed files that meet criterion `<n>`. flai refuses an orchestrator's acceptance without evidence, and one whose evidence names no changed file for a criterion. It writes the evidence, with the verified commit, under `### Accepted by the orchestrator` in the story's `## Notes`, in the acceptance commit, and returns it in `--json`. The done transition is recorded `by: orchestrator`, and the dashboard's review and story pages name who accepted.

flai holds the permission itself, as well as the guard: `flai accept --by orchestrator`, and `accept.run` given `by: orchestrator`, are refused while `accept_reviews` is off, naming it, and under `FLAI_ROLE=orchestrate` an acceptance with any other `--by` is refused. `item_move` refuses the orchestrator a move to done as it refuses every agent, and names `flai accept --by orchestrator` and the permission.

### Publishing a release (S-0222)

With `orchestration.permissions.publish` on, the orchestrator publishes what is accepted and not yet released, through the MCP tool `release_publish` alone, by the project's release policy ([ADR-0094](../adrs/0094-with-publish-on-the-orchestrator-publishes-only-through-release-publish-by-its.md), refining [ADR-0067](../adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md) and [ADR-0087](../adrs/0087-flai-serve-runs-one-orchestrator-per-project-behind-the-orchestrate-host-action.md)). The permission is the operator's asking that ADR-0067 requires. While it is off, the orchestrator neither evaluates nor publishes.

After each acceptance, its own or one `wait_for_events` reports as a story moved to done, it calls `release_evaluate` (`flai release --evaluate`, [project-manifest.md](project-manifest.md)) and decides by the policy:

| Policy | It publishes when | Its reason |
|--------|-------------------|------------|
| `threshold` | The evaluation is met: the unreleased stories' summed cost of delay is at or over `value`, or their count at or over `count` | The figure that was met |
| `theme` | The evaluation is met: every story of the epic or tag is accepted, and one at least is not yet released | The figure that was met |
| `judgement` | It judges the unreleased work coherent and complete | Its reasoning, which the tool requires |

Under every policy, with `orchestration.release.whole_epics` set, a batch with a story whose epic is in neither review nor done is not published: the evaluation names those stories in `held_by_epic`, and is not met while one is there. Each release publishes everything accepted and not yet released, as `flai release --pending` does.

`release_publish` (`internal/mcpserver/orchestrate.go`) takes `reason`, one sentence; under `threshold` or `theme`, the evaluation's reason when it is left out. It weighs the batch as `release_evaluate` does, and refuses, changing nothing, saying why and what would allow it (`mayPublish`):

- nothing is accepted and not yet released;
- `whole_epics` holds the batch back, naming each story held, its epic, and the epic's status, or `not found`;
- under `threshold` or `theme`, the policy is not met, with the evaluation's reason;
- under `judgement`, no reason is given;
- the `push` host action is off, which `publish.run` needs too, naming `flai serve enable push`;
- the caller is not the orchestrator, or `publish` is off (`Repo.OrchestratorPermits`), as well as the guard's refusal.

Otherwise it runs what the dashboard's `publish.run` runs, `flai release --pending`, never forced, through `hostapi.Publisher`, which `flai mcp` and `flai mcp http` fill (`app.mcpPublisher`). A refusal of `flai release --pending` with exit 3, for release tags on the remote newer than this clone's (S-0174) or a remote branch that moved, is returned as flai's message unchanged, beginning `conflict:`. A push that fails after the release is tagged here is returned as a refusal too, saying that the operator puts it right and publishes again. On success it returns `policy`, `reason`, `evaluation`, `versions` (each component's level, `from`, `to`, and items), `tags`, `items` (each ID once), `pushed`, `published`, and `warnings`.

What is recorded:

| What | Where |
|------|-------|
| The run | `serve/journal.jsonl`, as `publish.run`'s are: action `push`, method `mcp.release_publish`, `by` the orchestrator, and `<policy>: <reason>` before the outcome in its detail, so a release the orchestrator cut is told apart from the operator's. A run refused for the `push` action is journalled as `disabled` |
| Each release | `activity_log`: the policy, its figures, the versions and tags, and the items bundled, as `release_publish` returned them |
| Each decision not to publish | `activity_log`, with the evaluation's figures |
| Each refusal | `activity_log`, with the refusal; and a thread to the operator with `thread_open` on the most recently accepted story of the batch, with the refusal's words and what fixes it |

After a refusal it does not try again until that thread is answered or the next acceptance, and it never works around the refusal: the guard refuses it every other way to the remote ([Its permissions and the guard](#its-permissions-and-the-guard)).

### Its permissions and the guard

`orchestration.permissions` in `system-flow.yaml` says what the orchestrator may do without the operator. Each is off when unset. `flai check` reports an unknown key or a bad value as `manifest.orchestration`, so that a misspelt permission is not silently off ([project-manifest.md](project-manifest.md)).

`flai guard` reads `FLAI_ROLE` from the hook's environment. With `orchestrate`, it reads the permissions from the manifest of the project the hook runs in, at each call, and checks the session's own calls, which carry no agent ID (`Guard.orchestrate`). A manifest it cannot read leaves every permission off, with a warning. The orchestrator's sub-agents are held as every sub-agent is.

| Permission | Allows |
|------------|--------|
| none needed | The reads a sub-agent has, S-0217's and S-0219's evaluations among them (`order_by_policy`, `promote_candidates`, `release_evaluate`, `flai plan --candidates`, `flai promote --candidates` and `--drafts`, `flai order --by` without `--apply`, `flai release --evaluate`); `inbox`, `activity_log`, `wait_for_events`, and `thread_open`; `thread new`, and `issue new` and `bump`; git's reads |
| `plan_backlog_epics` | The MCP tool `plan` and `flai plan` on an epic; flai then lets it plan only an epic `flai plan --candidates` lists |
| `finalize_drafts` | `item_edit` giving `draft: false` and no field but `id`, `hash`, and `project` (`finalizesDraft`); `flai edit --no-draft` with nothing else to change: only `--by`, `--config`, `--hash`, `--json`, `--verbose`, and `--yes` beside it. flai then finalizes only a complete draft |
| `promote_to_ready` | `item_move` and `flai move` to `ready`; flai then moves only a candidate, while the ready column has room |
| `order_ready` | `flai order --by <policy> --apply` |
| `answer_threads` | On another's thread, `thread_reply` and `flai thread reply` as a recommendation (`recommendation: true`, `--recommend`) while it is `recommend` or `autonomous`, and as an answer citing a source (`source`, `--source`) while it is `autonomous` |
| `accept_reviews` | `flai accept` and `flai move <story> done`, each only with `--by orchestrator`; flai then accepts only with no blocker and evidence for every criterion ([Accepting a story](#accepting-a-story-s-0221)) |
| `publish` | The MCP tool `release_publish`; the tool then publishes only what the release policy allows ([Publishing a release](#publishing-a-release-s-0222)) |

A call a permission would allow is refused while that permission is off, and the refusal names it: `the orchestrator cannot <call>: it needs orchestration.permissions.<name>, which is off`. Every other write is refused whatever the permissions, saying that the orchestrator never does it, and why: `Edit`, `Write`, and `NotebookEdit`; every other flai MCP tool, `item_new` among them; `item_edit` that changes anything but the draft flag; `plan` on anything but an epic; `item_move` to any state but `ready`; `flai move` to any state but `ready`, and `done` but under `accept_reviews` with `--by orchestrator`; `flai edit` that changes anything but the draft flag; `flai order` that places a story by hand, which is the operator's (S-0219); `flai release` other than `--evaluate`, `flai push`, `git push`, and `git tag`, since it publishes through `release_publish` alone (S-0222); every other flai command that writes; and git's other writes. Each refusal ends by telling it to ask the operator with `thread_open` on the item. The guard refuses `release_publish` to every session but the orchestrator's, its sub-agents included. Other shell commands pass, as a sub-agent's do (ADR-0060).

Thread calls are held by who opened the thread as well (S-0220). On a thread it opened, the orchestrator follows up and resolves whatever its permissions, but never recommends or answers its own question. On another's, it replies only as `answer_threads` allows, and never resolves it. It never confirms a recommendation, which is the operator's alone, and never writes as another (`--by`). The guard reads who opened a thread from the thread itself, and takes a thread it cannot read as another's ([Answering threads](#answering-threads-s-0220)). The MCP tool `plan` and `flai plan`, asked by the orchestrator, are also checked in flai (`orchestratorPlans`): the item must be an epic `flai plan --candidates` lists, and `plan_backlog_epics` on ([What it does with each permission](#what-it-does-with-each-permission-s-0219)). The orchestrator is journalled as who asked. The `plan` host action must be on as well, since the planner is started as `serve.Plan` starts any. Every other agent `flai serve` started is still refused the tool.

Each refusal of the orchestrator's own call is appended to `wip/agents/orchestrator.md` under `## Refusals`, after `## Log`, by `Repo.AppendRefusal`: a heading with the time to the second, `- Call:` with the call refused in one code span, and `- Needs:` with the permission, or `none`. A refusal has no seconds or cost and adds nothing to the document's totals. A refusal that cannot be logged is warned of and refused all the same. Each writer of an activity document takes a lock per activity kind in `.flai-cache`, so that flai serve logging a run and the guard logging a refusal at the same time do not lose each other's entry; a lock older than ten seconds is taken as left behind.

`.claude/settings.json` runs the guard on `Edit|Write|NotebookEdit` when `FLAI_ROLE` is `plan`, `orchestrate`, or `analyze` ([The guard](#the-guard)).

### The run and how it ends

The run is recorded in `serve/agents.json` under `orchestrator`, the project's newest run alone, with the fields a story's run has and neither `story` nor `item`. `Running` and `Last` never name it, and the in-progress limit does not count it. Its output is in `serve/agents/<key>-orchestrator-<start>.log`, beside the other agents' logs, each run with a log of its own.

When the process ends, or the next look finds it gone, flai serve settles the run (`orchestrateEnded`):

| Outcome | When |
|---------|------|
| `stopped` | the `orchestrate` action was turned off and flai serve stopped it |
| `failed` | it exited with a code other than 0, or could not be started |
| `worked` | otherwise, an exit nobody saw included |

It then logs the run's activity since the last entry in `wip/agents/orchestrator.md` with `serve.LogRunEnd`: the run's final reply as the summary, or `stopped: orchestrate turned off`, and the seconds and cost measured from the log. Each decision the orchestrator logs with `activity_log` during the run covers its own span, so the end logs only what came after the last one.

### What it costs on the items (S-0226)

Each orchestrator activity is charged to the work items its entry names, when it is logged, through `activity_log`, through `thread_reply`'s own entry, or as the run ends ([ADR-0095](../adrs/0095-an-orchestrator-activity-s-usage-is-charged-evenly-to-the-work-items-it-named.md)). The usage apportioned to its span, the same share its entry's cost is, as for the planner ([The run and how it ended](#the-run-and-how-it-ended)), is split evenly between the epics, stories, and tasks it names that exist, archived ones included, each once; a thread, an issue, or anything else named takes no share. Tokens and seconds are split in whole numbers, the remainder to the first items named, and cost by the count, so that the shares add up to the entry. Each share goes under the `orchestrator` entry of `usage.strategic` on its item and on every item above it, up to its epic, at once ([work-hierarchy.md](work-hierarchy.md)). A charge that fails leaves the entry logged and is warned of.

An activity that names no work item charges no item: its cost is the orchestrator's project strategic total. That total, for every kind, is what the kind's activity document accrued that no item carries: its `accrued_cost` and `accrued_seconds` less that kind's `strategic` figures on the items with no parent, which carry every charge made below them. It is worked out by `flai stats`, which reports it per kind beside what the items carry ([metrics.md](metrics.md#strategic-agents-s-0206)), not written anywhere, so that the items and the project total always add up to the document. It also holds the planner's activities with no planned item, those logged before S-0225, and the analyzer's.

Every start, end, failure, and stop is a journal entry with action `orchestrate` and method `serve.orchestrate`. `agent.status` carries the run as `orchestrator`, null before the first; `agent.stream` with `orchestrator: true` reads its log as it reads a story's agent's (`OrchestratorStream`). When the run starts, cannot start, or ends, `flai serve` tells the dashboard with an `agent` notification that carries the project and `role: orchestrate`.

### On the dashboard

The Activity page lists the orchestrator's run and its stream. The Settings page lists the `orchestrate` action with the others, to turn on and off while `settings` is on. A page of its own, with its decisions and refusals, is S-0228's. See [flaiover-dashboard.md](flaiover-dashboard.md).

## The analyzer

The analyzer looks at the whole project, writes what it finds in one report, and ends (S-0223, [ADR-0099](../adrs/0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md)). It reads the metrics, the design, the code, and the issues, and looks for what its focus names: bottlenecks in the flow of work, gaps between `design/system` and the code (`intent`), or technical and security risks (`risk`), or all three when it is given no focus. It writes one report under `design/analysis/` and adds it to the folder's index. It edits nothing else and authors no stories. Filing its actionable findings as issues, which the convention's "As the analyzer" asks of it, is S-0224's: until then the guard refuses it `flai issue new` and `bump`, and its findings stay in the report.

### Starting an analysis

| Way | For | What it runs |
|-----|-----|--------------|
| `flai analyze [--focus bottlenecks\|intent\|risk] [--json]` | a shell on the host | `serve.Analyze` |
| The host API's `analyze.run` `{focus?}` | a dashboard | `flai analyze [--focus=<focus>] --json`; a focus none of the three is refused before anything runs |
| The MCP tool `analyze` `{focus?}` | the operator's own agent | `serve.Analyze`, as `flai analyze` does; journalled as `mcp.analyze` with the agent that asked |
| `analysis.schedule` | `flai serve`, on its own | the launcher's `analyze` with no focus, after the checks `serve.Analyze` makes ([The analysis schedule](#the-analysis-schedule)) |

Each needs the host action `analyze` on for the project: `flai serve enable analyze`, or the dashboard's Settings page while `settings` is on. It is off by default. A holder of the dashboard token can then start the analyzer whenever none runs. `flai analyze` starts the run itself, as `flai plan` does, and hands the process over: the command does not wait, and the serving flai settles the run. With no `flai serve` running it says so on standard error, and the analyzer still runs. Turning the action off starts no more runs; it does not stop one under way, which ends when the analyzer does.

`serve.Analyze` (`analyzeCheck`) refuses, saying why:

- while the `analyze` action is off for the project;
- a focus that is none of `bottlenecks`, `intent`, and `risk`;
- while an analyzer runs for the project, naming its focus, what started it, its PID, when it started, and its log: one analyzer runs per project at a time;
- when the analyzer's agent names no harness and no command is set on the host.

The MCP tool `analyze` also refuses an agent `flai serve` started (`FLAI_STARTED_BY=flai-serve`): an analysis is the operator's to ask for, or the schedule's. `flai guard` refuses the tool and `flai analyze` to a sub-agent, the planner, the orchestrator, and the analyzer.

### The analyzer's agent

`analysis.agent` in `system-flow.yaml`, merged over the project's `agent` as `planning.agent` is (`Manifest.AnalysisAgent`), gives the analyzer's harness, model, config, and roles. With neither set, the host's command starts it. `flai check` reports an `analysis.agent` or an `analysis.schedule` that is not valid as `manifest.analysis`. See [project-manifest.md](project-manifest.md).

### How the analyzer runs

The launcher's `analyze` starts it through the harness adapter with a `harness.Request` whose `Role` is `analyze` and whose `Focus` is the focus, or empty for all three, with no `Item` and no `Story`.

- **Where.** In the project's main checkout. No branch, no worktree.
- **As whom.** `FLAI_AGENT` is `analyzer`. The environment carries `FLAI_ROLE=analyze`, `FLAI_SESSION`, and `FLAI_STARTED_BY=flai-serve`, and no `FLAI_ITEM` or `FLAI_STORY`. `roleEnv` refuses an analyzer request that names a story or an item, or a focus none of the three, and a focus given to any other agent. `flai serve` drops `FLAI_FOCUS` from its own environment, with `FLAI_STORY`, `FLAI_ROLE`, and `FLAI_ITEM`, before it adds a run's (`runMarks`).
- **claude-code.** The session is named `<key> analyze`, with the focus after it when there is one, such as `sf analyze risk`. It runs with `--agent analyzer` over the project's `.claude/agents/analyzer.md`, passed in `--agents` with the analyzer's agent's model when it names one, beside the sub-agents whose roles the agent gives a model. A project with no `analyzer.md` is refused, naming `flai upgrade`. Its permission handler is flai's `permission_prompt`, as for every session `flai serve` starts.
- **The operator's command.** It gets `FLAI_ROLE`, and `FLAI_FOCUS` when the run has a focus, since a command has no prompt to be told it in. `{story}` is empty.

### What the analyzer is told

`harness.Prompt` gives the analyzer a prompt of its own (`analyzePrompt`): analyze the project and do nothing else, as the convention's section "As the analyzer" says; prime with `--role analyze`; call `inbox`. Look for the findings its focus names and nothing else, or for all three kinds with none:

| Focus | It looks for |
|-------|--------------|
| `bottlenecks` | Bottlenecks in the flow of work, from the cumulative flow, the time items spend in each state, the time they wait, and the holds on them |
| `intent` | Gaps between what `design/system` says and what the code does: a design section the code does not meet, and code the design does not describe |
| `risk` | Technical and security risks in the code and in how it is built, tested, and released |
| none, recorded as `all` | All three |

Read the metrics with `flai stats --json` and take every figure from it. Read the design with `doc_search` and `doc_get`, and the section a finding rests on before relying on it. Read the issues under `design/issues`, `summary.md` first. Hand wide search of the code to the explorer. Write one report, `design/analysis/<date>-<focus>.md` ([The report and its index](#the-report-and-its-index)), and add it to `design/analysis/README.md`. Edit nothing else, no code, design, issue, or work item, author no stories, and never work around a refusal. When an input the operator owns is missing, ask with `thread_open` on the report, the recommended answer first, analyze what needs no answer meanwhile, and hold `wait_for_events` until it is answered. End with a one-line summary that names the report, which becomes the run's activity entry.

`analyzer.md`, which the template ships, says the same in short and lists its tools: `Read`, `Grep`, `Glob`, `Edit`, `Write`, `Bash`, `Agent`, and flai's MCP tools `prime`, `inbox`, `board`, `item_get`, `doc_get`, `doc_search`, `thread_get`, `who_touches`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events`. It lists no `NotebookEdit` and no tool that writes a work item.

### The report and its index

A report is `<date>-<focus>.md` in the design folder's `analysis/` (`analysis.Dir`, `analysis.ReportPath`): the date is the run's UTC day as `YYYY-MM-DD`, and the focus one of the three or `all`. Two runs on one day with one focus name the same file.

| Front matter | What it holds |
|--------------|---------------|
| `title` | The report's title |
| `updated` | A date, `YYYY-MM-DD`, or a timestamp, `YYYY-MM-DDTHH:MM:SSZ` |
| `status` | The documentation standard's: `draft` while the analyzer writes it, `active` once the run has ended, `deprecated` when a later report replaces it |
| `focus` | The focus the file is named for |
| `from`, `to` | The window its metrics cover, as dates, `from` not after `to` |

The body has one section per finding: its evidence (the metric figures as flai gave them, the file paths, and the design sections quoted), its severity, and its estimated impact, the time it loses per cycle, or the revenue or penalty it puts at stake where the design states them.

`design/analysis/README.md`, which the template ships, says what the folder is for and has a row per report: its link, focus, window, and status. Only the analyzer writes in the folder. `flai check` leaves the folder out of the design documents' own rules and checks it by its own (`checker.analysis`):

| Finding | Level | When |
|---------|-------|------|
| `analysis.report` | error | A report not named `<date>-<focus>.md` with a real date and a known focus; with no front matter, or front matter that does not parse; missing a field; with an `updated` that is neither a date nor a timestamp, a `status` none of the three, a `focus` none of the four or not the one the file is named for, a `from` or `to` that is not a date, or a window that ends before it starts. Each is reported on its field's line |
| `analysis.index` | warning | A report the README does not link, and a row of the README that links a report not there |

A project with no `design/analysis/` has nothing checked there. flai checks a report's status but does not change it: the analyzer writes it, `draft` while it writes the report and `active` once it is done, as its prompt and its definition tell it.

### The analyzer's guard

`flai guard` reads `FLAI_ROLE` from the hook's environment. With `analyze`, the command opens the project from the working directory and sets the guard's `Root` and `Reports`, the reports folder, and `Guard.analyze` checks the session's own calls, which carry no agent ID:

| Kind | Passes | Refused |
|------|--------|---------|
| File tools | `Edit`, `Write`, and `NotebookEdit` on a file inside the reports folder: its report and the folder's README | the same on any other file, and on every file when the project cannot be read |
| flai MCP tools | the reads a sub-agent has; `inbox`, `activity_log`, `thread_open`, `thread_reply`, and `wait_for_events` (`guard.MCPAnalyzes`) | `item_new`, `item_edit`, and `item_move`, since the analyzer authors no stories; every other tool, `plan`, `analyze`, `issue_story`, and `release_publish` among them |
| flai commands | the reads, `flai stats` among them | every command that writes, `issue new` and `bump`, `plan`, and `analyze` among them |
| git | the reads a sub-agent has | every other git command |

A file's path is taken from the project's root when it is relative, cleaned, and resolved along the part of it that exists, symbolic links followed, as the folder's is (`inReports`), so that neither `..` nor a link leads out of the folder. A refused edit names the folder. Each refusal says that the analyzer edits nothing but its report and the index, and to put what it found, and the stories it would suggest, in its report, or ask the operator with `thread_open`. Other shell commands pass, as a sub-agent's do (ADR-0060), so a file written through the shell is not seen. The analyzer's sub-agents, the explorer among them, are held as every sub-agent is.

`.claude/settings.json` runs the guard on `Edit|Write|NotebookEdit` when `FLAI_ROLE` is `analyze`, as for `plan` and `orchestrate` ([The guard](#the-guard)). A project that keeps its old settings has an analyzer whose file edits the guard does not see.

### The analysis schedule

`analysis.schedule` in `system-flow.yaml` is a five-field cron expression in UTC, or `daily`, read as `planning.schedule` is ([The schedule](#the-schedule)); unset, there is none. An analysis scheduler per served project (`internal/serve/analysis_schedule.go`) looks whenever the project's launcher does: when its work items or threads change, when an agent ends, and every minute, so it sees the schedule come round at most a minute late. Each time the schedule comes round with the `analyze` action on, it starts the analyzer with no focus and the trigger `schedule <spec>`, the schedule as written, such as `schedule daily`, after the checks `serve.Analyze` makes. A start they refuse, such as while an analyzer runs, is logged at info as `scheduled analyzer not started` with the reason, and is not tried again until the schedule next comes round.

Times missed between two looks come round once. A schedule set or changed comes round first at its next time from then. What the scheduler has seen is kept in memory: a time that passed while `flai serve` was down does not come round. While `analyze` is off it starts nothing and keeps up with the schedule, so turning the action on does not act on the past. A schedule that is not valid is warned of once in flai serve's log, and nothing is started for it until it is fixed. The scheduler follows its schedule as the replanner follows `planning.schedule` (`cadence`).

### The analyzer's run and how it ended

The run is recorded in `serve/agents.json` under `analyzer`, the project's newest run alone, with the fields a story's run has, neither `story` nor `item`, and `focus` (one of the three, or `all`), `trigger` (`asked`, or the schedule's), and, once it has ended, `report`. `Running` and `Last` never name it, and the in-progress limit does not count it. Its output is in `serve/agents/<key>-analyzer-<start>.log`, beside the other agents' logs.

When the process ends, or the next look finds it gone, flai serve settles the run (`analyzeEnded`). Its report is the Markdown file directly in the reports folder, the README aside, changed last, and changed no earlier than the second the run started (`newestReport`); there is none when no such file changed.

| Outcome | When |
|---------|------|
| `failed` | it exited with a code other than 0, or could not be started; or it ended without writing a report, which `why` says, naming the folder |
| `worked` | otherwise: it wrote a report, whether its exit was seen or not |

It then logs the run's activity in `wip/agents/analyzer.md` with `logRunEndSaying`: the time since the last entry, the seconds and cost measured from the log, and the run's trigger as the entry's `- Trigger:` line. The summary is the first line of the run's final reply, the one-line summary its prompt asks for: as it is when it names the report, with `(report <path>)` added when it does not, and with `(no report written under design/analysis)` added when there is no report. The entry names no item, so its cost charges no item: it is the analyzer's project strategic total ([What it costs on the items](#what-it-costs-on-the-items-s-0226)) until S-0227 charges it to the issues the analyzer files.

Every start and failure to start is a journal entry with action `analyze` and method `serve.analyze`; `analyze.run` and the MCP tool `analyze` are journalled as well, under the same action. `agent.status` carries the run as `analyzer`, null before the first. When the run starts, cannot start, or ends, `flai serve` tells the dashboard with an `agent` notification that carries the project and `role: analyze`.

### The analyzer on the dashboard

The documents page lists the reports under `design/analysis/` beside the folder's README and shows each as it shows any design document ([flaiover-dashboard.md](flaiover-dashboard.md)). The Settings page lists the `analyze` action with the others, to turn on and off while `settings` is on. No page has an Analyze button yet: `analyze.run` is there for the analyzer's page, with its runs and activity, which is S-0228's.
