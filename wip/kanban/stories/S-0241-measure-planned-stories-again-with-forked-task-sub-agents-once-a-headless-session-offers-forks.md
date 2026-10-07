---
id: S-0241
type: story
nature: experiment
title: Measure planned stories again with forked task sub-agents once a headless session offers forks
status: backlog
owner: arobson
created: 2026-10-02T17:14:07Z
updated: 2026-10-07T23:14:44Z
transitions: []
tags: [template]
topics: [conventions]
touches: [design/system/agent-context.md, design/experiments/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 348
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 696
          cache_read: 7455294
          cache_write: 25044
          cost: 1.8431
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: orchestrator
    at: 2026-10-07T23:14:09Z
  value: 25
  by: planner-S-0241
  at: 2026-10-07T23:14:44Z
forecast:
  duration: 2h
  delivery: 2026-10-08T04:25:00Z
  basis: "S-0176 took 1h45m from in-progress to review and 6355 s of agent time for the same three replays, their blind reviews, and a results document, so 2h replaces flai's 8m from size 5; delivery is flai's pull-order date moved by that difference, and holds only if a headless session offers forks by then, which Claude Code 2.1.290 does not."
  by: planner-S-0241
  at: 2026-10-07T23:12:30Z
---
# S-0241 Measure planned stories again with forked task sub-agents once a headless session offers forks

## Goal

S-0176 measured a story's agent handing its tasks to task sub-agents, and found that each sub-agent started fresh and read the code again, because the headless `claude -p` session flai serve starts refused `subagent_type: fork`. A fork inherits the parent's conversation and prompt cache, so it would start primed. S-0230 probed again on 2026-10-02 (Claude Code 2.1.286) and forks were still refused (TH-0069). When a headless session offers forks, repeat S-0176's measurement with forked task sub-agents and record whether the sub-agents' cost falls.

## Acceptance criteria

- [ ] A headless session started as flai serve starts one is shown to offer `subagent_type: fork`, with the Claude Code version that does
- [ ] The measurement S-0176 made (`design/system/agent-context.md` § Tasks in parallel) is repeated with forked task sub-agents, and `design/system/agent-context.md` records the result beside the first
- [ ] The experiment's results document is written under `design/experiments`, with a recommendation

## Tasks

- T-1267 Show a headless session flai serve starts offering subagent_type fork, and record its Claude Code version
- T-1268 Replay S-0186, S-0185, and S-0182 with forked task sub-agents and record the runs beside S-0176's in agent-context.md
- T-1269 Write S-0241's results document under design/experiments with a recommendation

## Notes

Split from S-0230's fourth criterion on TH-0069. Do not move it to ready until a headless session offers forks: it cannot be worked before then.

### Planning

Touches, file by file, no folder touch kept:

- `design/system/agent-context.md`: from the design. Criteria 1 and 2 are recorded in its § Tasks in parallel, in § Forks are not offered headless and § Measured against runs without the plan.
- `design/experiments/S-0241-measure-planned-stories-again-with-forked-task-sub-agents-once-a-headless-session-offers-forks.md`: from the layout. Criterion 3 asks for the results document, and `design/experiments/README.md` names it `<S-nnnn>-<slug>.md` after the story's file. The README has no list of documents, so it is not touched.
- `flai touches suggest` declared none to start from. Started from the two paths above, it listed co-changes of `agent-context.md` from code stories (`design/system/flai-cli.md`, `docs/users/flai.md`, `flai/internal/harness/harness_test.go`, `design/conventions/delegation.md`, each 8 to 43%). None is added: this story changes no code and no convention. Adopting forks in `delegation.md` and the `claude-code` prompt is a later story that the recommendation names. If forks need a setting in the `claude-code` adapter, T-1268 widens its touches to `flai/internal/harness`.

Forecast: `flai forecast` gave 8m (86 s per unit of size times size 5). Adjusted to 2h, because S-0176 took 1h45m from in-progress to review, and 6355 s of agent time, for the same three replays, their blind reviews, and a results document. Delivery is flai's pull-order date, 2026-10-08T02:33Z, moved later by that difference to 2026-10-08T04:25Z. It holds only if a headless session offers forks by then. The planner's own headless session, Claude Code 2.1.290 on 2026-10-07, offered none.

Cost of delay: 25 USD a week, as `flai cod` works it out from the orchestrator's input on TH-0332, `penalty_per_week: 25`. Kept as given: it is the forks' expected saving on sub-agent cache reads across the week's planned stories.

Tasks, in three layers, one task each, since each needs the one before:

1. T-1267: show the fork offered headless, with its version.
2. T-1268: replay S-0186, S-0185, and S-0182 with forked task sub-agents, after T-1267.
3. T-1269: the results document, after T-1268.
