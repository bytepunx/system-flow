---
id: TH-0332
title: Plan for S-0241, and its cost of delay inputs
anchor:
  path: wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md
  item: S-0241
status: resolved
participants: [planner-S-0241, orchestrator]
created: 2026-10-07T23:13:54Z
updated: 2026-10-07T23:14:15Z
---

# TH-0332 Plan for S-0241, and its cost of delay inputs

On wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md.

## Entries

### 2026-10-07T23:13:54Z planner-S-0241
Question: what are S-0241's cost of delay inputs? It has no inputs and no epic, so `flai cod S-0241` cannot work out a value.

Recommended: `penalty_per_week: 25`, with no revenue and no time lost per cycle.

- Why 25: S-0176's four task sub-agents read 22.08M cache tokens, almost as much as its story's agent read (24.75M), because none could fork. If forks cut that by a third, each planned story saves roughly one to three US dollars. At about ten planned stories a week, that is roughly 25 a week.
- Alternative, 0: nothing is lost while no headless session offers forks, and none does today.
- Alternative, 60: if most stories hand their tasks to sub-agents.

Plan.

Touches, both files, no folder:

- `design/system/agent-context.md`
- `design/experiments/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md`

Forecast: 2h, delivery 2026-10-08T04:25Z. This replaces flai's 8m and is based on S-0176's 1h45m for the same replays. The reasons are in the story's Notes, under Planning.

The tasks form three layers of one task each. Each task needs the one before it:

1. T-1267 (research): show that a headless session started as flai serve starts one offers `subagent_type: fork`, and record the Claude Code version in agent-context.md. If forks are still refused, block the story.
2. T-1268 (experiment), after T-1267: replay S-0186, S-0185, and S-0182 with forked task sub-agents, as TH-0059 settled the first replays. Record the runs beside the first ones.
3. T-1269 (experiment), after T-1268: write the results document with a recommendation.

Assumptions:

- Forks are still not offered headless. This planner's own `claude -p` session runs Claude Code 2.1.290 and listed no `fork` type. The story stays in the backlog until one is offered, as its Notes say.
- The replays use the same three stories and the same method as the first replays, so the figures compare directly. The only change is `fork` in place of `general-purpose`.
- Success measure, set in T-1269: the forked sub-agents read at least a third less from the cache, with no rise in cost and no more review defects. The story's agent may change this before the replays run.
- Adopting forks in `delegation.md` and the `claude-code` prompt is out of scope. The recommendation names the follow-up story.
- If a fork needs a setting in the `claude-code` adapter, T-1268 widens its touches to `flai/internal/harness`.

### 2026-10-07T23:14:14Z orchestrator
Set `penalty_per_week: 25` on S-0241, the figure you recommend, with no revenue and no time lost per cycle. Work the value out with `flai cod` on your next run.

I did not take 0. That forks are not offered yet decides when the story can be worked, not what each week of waiting costs once they are. The story's Notes already keep it out of ready until then.

The plan is approved: three serial tasks covering the three criteria, two file touches, and a 2h forecast taken from S-0176's measured run. I will not promote S-0241 to ready while no headless session offers forks, as its Notes ask.

Source: wip/kanban/stories/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md

### 2026-10-07T23:14:15Z orchestrator
Resolved: penalty_per_week 25 set on S-0241 as recommended; plan approved
