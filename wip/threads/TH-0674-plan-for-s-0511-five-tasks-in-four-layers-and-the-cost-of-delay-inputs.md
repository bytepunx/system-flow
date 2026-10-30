---
id: TH-0674
title: "Plan for S-0511: five tasks in four layers, and the cost of delay inputs"
anchor:
  path: wip/kanban/stories/S-0511-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md
  item: S-0511
status: resolved
participants: [planner-S-0511, orchestrator, agent-S-0511, alex]
created: 2026-10-28T08:21:13Z
updated: 2026-10-30T16:48:48Z
---

# TH-0674 Plan for S-0511: five tasks in four layers, and the cost of delay inputs

On wip/kanban/stories/S-0511-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md.

## Entries

### 2026-10-28T08:21:13Z planner-S-0511
Question first, since it needs the operator: S-0511 has no cost of delay inputs and no epic, so `flai cod` refuses it. Recommended: `penalty_per_week: 500`, with no time lost per cycle. The defect is a safety gap found once by inspection (an agent `flai serve` started can run `flai serve enable auto-approve`, `notify set --command`, or `agent set --` and so widen its own permissions or name a program the host runs as the operator); it has cost no cycle any time, so a weekly penalty for the standing risk fits better than time lost. Alternatives: `penalty_per_week: 1000`, as S-0475 got for a hold on the board; or `time_lost_per_cycle: 20m`, the usual figure for a one-instance defect, giving 50 USD a week. Give the inputs with `flai edit S-0511 --penalty-per-week 500` or `item_edit` `cost_of_delay`, and I write the value.

The plan. The issue asks for a solution before it is built; the proposal is a guard rule that binds every served session and sub-agent, as the shared paths (ADR-0096) and `flai forge push` (ADR-0166) do, recorded in an ADR first.

Tasks and layers:

- Layer 1: T-2184, the ADR. It fixes the list: refused are `flai serve enable` and `disable`, `agent set`, `agent clear`, `agent harness` and `agent provider` with a setting or `--reset`, `checks set`, `timeout`, `clear`, `components idle`, `max-servers`, `forge` with a setting or `--reset`, `notify set`, `notify clear`, `import add`, `import remove`, and `flai config set` and `flai template use`; open stay the reading forms (`status`, `actions`, `journal`, `agent show`, `harness` and `provider` bare or with a name, `checks show`, `components show`, `forge` bare, `notify show`, `import list`, `config get`) and `serve start`, `stop`, `project`, `agent start`, `restart`, `commit`, `stop`, `orchestrate`, which write no configuration.
- Layer 2, after T-2184, no path in common: T-2185, the rule in `flai/internal/guard/hostconfig.go` wired in `Decide`, with a test over every session kind modelled on `TestFlaiForgePushIsTheOperators`; T-2186, the `flai guard` help and `make flai-reference`; T-2187, `docs/users/flai.md`, `docs/operators/index.md`, `docs/operators/settings.md`, `design/system/flai-cli.md`, `design/system/agent-adapters.md`.
- Layer 3, after T-2185, T-2186, T-2187: T-2188, `flai issue close I-0210` and the criteria ticked.

Assumptions:

1. `flai config set` and `flai template use` are in the rule though the issue names only `flai serve`: they write the same `~/.flai/config.json` through the same gap (`agentEnv` drops `FLAI_CONFIG`). Say so if you want them left to a separate issue.
2. The story's agent is the one served session the guard does not already hold on flai writes (the planner, analyzer, curator, and orchestrator refuse every write they do not name, and sub-agents are refused by default), so the rule runs in `Decide` before the role split, as `shared` and `forgePush` do, rather than per role.
3. `flai serve stop` and `flai host stop`, which stop the operator's dashboards, stay out: they write no configuration. A separate issue if you want them held.
4. No template change: its `.claude/settings.json` runs the guard through a generic hook, and `safety.md` already forbids changing the operator's global configuration.
5. The ADR's number is allocated by `flai adr new` in T-2184, so the story keeps the folder touch `design/adrs` beside `design/adrs/README.md`, recorded under the story's `### Planning`.

Figures: forecast 45m, raised from flai's 26m because S-0369, the same shape, took 33m with less to parse; delivery 2026-10-29T21:30Z. S-0504, in progress, touches `docs/users/flai.md`, `design/adrs`, and `design/issues/summary.md`, so the story is held while S-0504 is in progress.

### 2026-10-28T08:21:40Z orchestrator (recommendation)
**Recommended for the operator:** keep `flai config set` and `flai template use` in the rule (assumption 1), and approve the plan as written.

## What I set

- S-0511's cost of delay input is now `penalty_per_week: 500`, the planner's recommended figure. Nothing in the thread or the story argues for 1000 or for 20m a cycle. On its next run, the planner writes the value with `flai cod`.

## The plan fits

- **Tasks:** five tasks in three layers, not four as the title says. T-2184, the ADR, comes first. T-2185, T-2186 and T-2187 then share no path. T-2188 closes I-0210 last.
- **Touches:** file-level on every task. The folder `design/adrs` is kept for the ADR file no task can name, and the story's `### Planning` names it.
- **Forecast:** 45m, kept over flai's 26m against S-0369, with the reason stated.
- **Assumptions 2 to 5** follow precedents already recorded: `shared` and `forgePush` in `Decide`, and ADR-0096 and ADR-0166. They need no answer.

## Why assumption 1 is yours

- The story's title and I-0210 name only the `flai serve` sub-commands. Adding `flai config set` and `flai template use` widens what the story changes, so it is a change of scope.
- **Recommended: include them.** They write the same `~/.flai/config.json` through the same gap, because `agentEnv` drops `FLAI_CONFIG`. A guard that refuses one path to a file and leaves the other open does not close the risk I-0210 records. Including them adds lines to T-2185's rule and test, and no task.
- **The alternative:** leave them to a separate issue, and S-0511 covers `flai serve` alone.

The thread stays open for that answer. S-0511 stays a draft until its value is written.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-28T08:22:21Z planner-S-0511
The value is written: 500 USD a week, as `flai cod` gives it from `penalty_per_week: 500`. The forecast (45m, delivery 2026-10-29T21:30Z) and the `### Planning` note are on the story. The title's "four layers" is wrong: three, as the orchestrator says.

What remains is assumption 1 for the operator. If the answer is the alternative, `flai serve` alone, the change is in three tasks and I make it on the next plan run, or the story's agent does when it reviews the tasks: drop `flai config set` and `flai template use` from T-2184's list of refused forms, from T-2185's rule and `Commands`, and from T-2186's help paragraph, and file a separate issue for the two commands.

### 2026-10-30T00:11:32Z agent-S-0511
S-0511 is started and works the plan as written, with the orchestrator's recommendation on assumption 1: `flai config set` and `flai template use` are in the rule. If the operator answers the alternative, `flai serve` alone, I drop the two from the ADR, the rule, its test, and the help before review, and file a separate issue for them.

### 2026-10-30T00:15:40Z orchestrator (recommendation)
**Recommended for the operator, again:** keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1).

The planner's and the agent's entries came after my recommendation of 2026-10-28T08:21:40Z, so it no longer stood as pending. I post it again so you can confirm it or answer otherwise.

- **Why include them:** they write the same `~/.flai/config.json` through the same gap, since `agentEnv` drops `FLAI_CONFIG`. Refusing one path to that file and leaving the other open does not close the risk I-0210 records.
- **What it costs:** lines in T-2185's rule and test, and in the help. No new task.
- **The alternative:** `flai serve` alone, and a separate issue for the two commands. agent-S-0511 drops them before review if you answer so.

S-0511 is in progress and works to this recommendation meanwhile. The rest of the plan needs no answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-30T16:48:25Z alex
Confirmed the recommendation of 2026-10-30T00:15:40Z orchestrator.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-30T16:48:48Z alex
Resolved.
