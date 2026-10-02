---
title: "Planning a story's tasks in layers and working each layer with sub-agents at once"
updated: 2026-10-02
status: active
story: S-0176
---

# Planning a story's tasks in layers and working each layer with sub-agents at once

The results of S-0176, an `experiment` story run on 2026-10-01. The full measurement, with its method and the per-run figures, is in [agent-context.md § Tasks in parallel](../system/agent-context.md). This document records the experiment as [ADR-0066](../adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md) asks.

## Hypothesis

Within a story, one agent works the tasks one after another, even when they touch disjoint paths and do not depend on each other. If the story's agent planned its tasks, recording in each task's front matter (`after`) which must wait for which, and handed the tasks that can run together to sub-agents at once, a story would finish sooner at an acceptable cost. Forked sub-agents, which inherit the parent's conversation and prompt cache, were expected to make each task's sub-agent start already primed.

## Success measure

For stories with three or more tasks of which at least two are independent: wall-clock time from `in-progress` to `review`, less time waiting on the designer, falls by a third or more; cost (ADR-0051's usage, sub-agents included) rises by no more than half; and defects found at review do not rise. Measured over at least four stories run with the plan against four comparable runs without it.

## What was done

Built on `story/S-0176`, in six tasks worked as the plan itself prescribed (two layers of two, then one, then one):

- Tasks carry `after: [T-nnnn]`, the tasks of their story they wait for, set by `flai task new --after`, `flai edit --after`, and the MCP write tools, and checked by `flai check` (a task of another story, the task itself, cycles). `flai show`, the board, and the story's page and card in flaiover show each task's state and the plan as layers of tasks that can run at once.
- The conventions (`delegation.md`, `work-management.md`, here and in the template) and `harness.Prompt` for `claude-code` tell the story's agent to plan its tasks and to hand each to a task sub-agent, and tell a task sub-agent what it may touch and that only the story's agent commits, syncs, moves items, and talks to the designer. The guard (ADR-0060) enforces the last.
- Two worktree modes were tried: layer 1's tasks shared the story's worktree; layer 2's each had a worktree of their own made from the story branch and merged back by the story's agent. Sub-agents ran as `general-purpose` agents with self-contained prompts, because the Agent tool refused `fork` in the headless session.
- Measurement, as TH-0059 settled: three archived stories with three or more tasks of which at least two were independent (S-0186, S-0185, S-0182) were replayed with the plan in scratch clones at their original base commits, with this branch's flai, prompt, and conventions, and set against their own original runs; S-0176's own run was set against S-0189's (six tasks, worked in turn). Review defects were a fresh verifier's blind findings on each pair of diffs. The replays also carried S-0175's and S-0189's delegation prompt, which the originals did not; that confound is recorded, not removed.

## Results

| Run | Plan | Minutes to review | Cost (USD) | Agent cache reads | Review defects (high / medium / low) |
|-----|------|-------------------|------------|-------------------|--------------------------------------|
| S-0186 original / replay | none / 2 layers | 11.4 / 10.9 | 6.04 / 5.13 | 12.56M / 8.05M | 0/2/2 / 0/1/2 |
| S-0185 original / replay | none / 3 layers | 22.4 / 22.8 | 7.83 / 5.62 | 15.48M / 7.52M | 0/0/2 / 0/3/3 |
| S-0182 original / replay | none / 2 layers | 23.1 / 42.9 | 11.13 / 12.10 | 29.53M / 22.43M | 3/1/0 / 0/1/2 |
| S-0189 / S-0176 (build tasks) | none / 4 layers | 53.5 / 24.0 | 25.82 / 17.79 (both est.) | 48.07M / 24.75M | not compared |

Against the success measure:

- **Time: not met.** No replay reached review a third sooner. S-0186 saved half a minute, S-0185 lost a fraction of one, and S-0182 took nearly twice as long, because the replay clone's own older flai refused the `after` the branch's flai wrote on its tasks and the agent undid it by hand. A layer did run in the time of its slowest task (S-0176's two layers of two took 16 sub-agent minutes where four tasks in turn would have taken 24), but the tasks were three to nine minutes each, and most of a run is the story's agent's own fixed work: priming, planning, reviewing each sub-agent's work, the verifier before review, and the close-out, none of which shrinks with layers.
- **Cost: met.** Two replays cost 15% and 28% less, one 9% more. The saving is the story's agent's smaller context: with the tasks' reading and editing in sub-agents it made 14 to 38% fewer model calls and read 24 to 51% less from the cache. Where the sub-agents read a great deal (S-0182, five times the original's sub-agents; S-0176, almost as much as its agent), the saving went. Each task sub-agent started fresh and read the code again; forks, which would start from the parent's cache, were not offered.
- **Defects: inconclusive.** One replay had fewer medium findings, one more, and S-0182's replay found the cause its original had missed. Variance between two runs of the same story was larger than any effect of the plan.
- **Two defects found on the way.** An older flai refuses a task's `after`, so adopting it needs the flai that carries it on the host and `flai.minimum` raised (S-0181). And ADR-0051's per-task usage apportions a run by the task's time window, which assumes one task at a time: tasks worked at once each got the whole window (I-0054).

## Recommendation

**Adapt.** Keep a task's `after` and the plan's display, which make a story's order of work visible to the agent, the board, and the designer at no cost, and keep handing a story's tasks to sub-agents, which keeps the story's agent's context small and is where the cost fell. Drop the expectation of a faster story: running a layer's tasks at once is worth it only where tasks are long beside the story's fixed costs, which no story the project runs today is, so the conventions make it optional rather than required. Measure again once forked sub-agents are offered to a headless session.

Nothing is adopted by this experiment itself, so there is no ADR. [S-0230](../../wip/kanban/stories/S-0230-tasks-carry-after-the-plan-shows-on-the-board-and-a-story-s-agent-hands-tasks-to-sub-agents-without-parallel-layers.md) carries it on: the host's flai with `after` and a raised minimum, I-0054 fixed, the conventions and prompt saying the plan without the parallel layers, and the repeat measurement with forks.
