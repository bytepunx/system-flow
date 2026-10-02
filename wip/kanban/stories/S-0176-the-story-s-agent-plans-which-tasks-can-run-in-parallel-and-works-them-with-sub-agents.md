---
id: S-0176
type: story
nature: experiment
title: The story's agent plans which tasks can run in parallel and works them with sub-agents
status: review
owner: alex
created: 2026-10-01T07:17:38Z
updated: 2026-10-01T13:22:12Z
transitions:
  - to: ready
    at: 2026-10-01T07:38:54Z
    by: alex
  - to: in-progress
    at: 2026-10-01T11:37:09Z
    by: agent-S-0176
  - to: review
    at: 2026-10-01T13:22:12Z
    by: agent-S-0176
tags: [flai, template]
topics: [conventions]
touches: [flai/internal/workitem, flai/internal/check, flai/internal/harness, flai/internal/mcpserver, flai/cmd, flaiover/src, template, design/conventions, design/system/work-hierarchy.md, design/system/agent-context.md, design/system/agent-coordination.md, design/system/conventions.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/issues, docs/operators/settings.md, docs/users/flai-reference.md, docs/users/flai.md, flai/internal/itemedit, flai/internal/hostapi]
after: [S-0175]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 6355
  models:
    - model: claude-haiku-4-5-20251001
      input: 678
      output: 20397
      cache_read: 3751049
      cache_write: 194388
      cost: 0.7208
    - model: claude-opus-5-5
      input: 790
      output: 254128
      cache_read: 56159779
      cache_write: 801367
      cost: 21.2798
    - model: claude-sonnet-5-5
      input: 132
      output: 28884
      cache_read: 2503172
      cache_write: 260960
      cost: 1.4421
---
# S-0176 The story's agent plans which tasks can run in parallel and works them with sub-agents

## Goal

Within a story, one agent works the tasks one after another, even when they touch disjoint paths and do not depend on each other. Parallelism exists only across stories (the in-progress limit and claims, ADR-0046). This experiment tests whether a story finishes faster, at an acceptable cost, when its agent plans the tasks for parallel work and hands independent ones to sub-agents.

The story's agent owns the plan. When it writes a story's tasks it decides, for each, whether it can start at once or must wait for a named task to be done, and records that in the task's front matter. It then runs the tasks whose dependencies are done, in parallel where their `touches` do not overlap, each with a sub-agent, and integrates, commits, and keeps the narrative itself. Forked sub-agents, which inherit the parent's conversation and prompt cache, are the first candidate, so a task's sub-agent starts already primed.

The hypothesis: for stories with three or more tasks of which at least two are independent, wall-clock time to review falls by a third or more, cost rises by no more than half, and defects found at review do not rise.

## Acceptance criteria
- [x] Tasks carry `after: [T-nnnn]`, the tasks of the same story that must be done before it starts, set with `flai task new --after` and `flai edit --after`, and the MCP `item_new` and `item_edit` equivalents. `flai check` reports an entry that names no task, a task of another story, the task itself, and every cycle. Front matter stays strict: the design says which flai the host must run first
- [x] `flai show`, the board, and the story's page in flaiover show each task as ready to start, waiting with the tasks it waits for, or in progress, and the plan as layers of tasks that can run at once
- [x] The conventions and `harness.Prompt` for `claude-code` tell the story's agent to plan when it writes the tasks: what each task touches, which must wait for which and why, and which can run together (an `after` and no overlap in `touches`), and to record the plan's reasoning in the narrative's `## Decisions`
- [x] The story's agent runs ready, non-overlapping tasks in parallel with sub-agents, waits for them, reviews their work, and moves each task; only it commits, syncs the stream, moves items, and talks to the designer. A sub-agent's question for the designer reaches the designer by the mechanism S-0175 chose
- [x] Parallel sub-agents do not corrupt each other's work in the shared worktree: either each task gets its own worktree from the story's branch and the parent merges it back, or the experiment shows non-overlapping `touches` are enough in one worktree. Record which and the evidence
- [x] At least four stories run with the plan and four comparable ones without it, comparing wall-clock time to review, cost and tokens (ADR-0051, sub-agents included), the number of parallel layers, conflicts, and review defects. The results and a recommendation (adopt, adapt, or drop) go in `design/system/agent-context.md`, with an ADR if adopted
- [x] `design/system/work-hierarchy.md` documents a task's `after`

## Tasks
- T-0677 Tasks carry after, set by flai and MCP and checked by flai check
- T-0678 The conventions and the claude-code prompt tell the story's agent to plan its tasks and work independent ones with sub-agents
- T-0679 flai show and the board give each task's state and the plan's layers
- T-0680 The story's page and its board card in flaiover show the task plan
- T-0681 Record whether parallel sub-agents need a worktree each, with this story's evidence
- T-0682 Measure planned runs against runs without the plan and recommend adopt, adapt, or drop

## Notes

- Depends on S-0175: sub-agents that cannot move items or consume the inbox, and the decision on how they raise questions.
- A task's `after` mirrors a story's (ADR-0046, S-0130), within a story only. Whether flai should also hold a task (refuse or warn on `flai move T-nnnn in-progress` while its `after` is not done) is part of the experiment; start with a warning, as for stories.
- Risks to watch: parallel edits to one worktree, the narrative and `flai stream sync` assuming one writer, and Anthropic's estimate that multi-agent runs use about fifteen times a chat's tokens, which ADR-0049 cited in rejecting a scout sub-agent. Forks that share the prompt cache may change that figure; measure it.
- Worked as its own first planned run (narrative `## Decisions`): layer 1 T-0677 and T-0678 in the story's worktree, layer 2 T-0679 and T-0680 each in its own worktree, then T-0681 and T-0682 by the story's agent. The tasks carry no `after` in main's `wip/`, because the flai installed on the host refuses it on a task until a release with T-0677 is installed; the plan is in the narrative.
- Forks were refused in the headless session ("Agent type 'fork' not found"), so every task sub-agent was `general-purpose` with a self-contained prompt.
- Measurement (TH-0059, option a): S-0186, S-0185, and S-0182 replayed with the plan in scratch clones under `/tmp/s0176` (no remote) against their original runs, and this run against S-0189. Recommendation: adapt. Keep a task's `after`, the plan's display, and handing tasks to sub-agents; drop the parallel layers. No story reached review sooner, so there is no ADR. See `design/system/agent-context.md § Tasks in parallel`.
- A task's state counts a cancelled `after` task as satisfied, unlike a story's hold, and a done task keeps its layer.
- I-0054 recorded: flai's per-task usage gives tasks worked at once the whole window each.
- After the pre-review verifier: the conventions and `harness.Prompt` on this branch ask for each task to go to a task sub-agent and allow a layer's tasks at once only when they are long, as the recommendation says, instead of requiring it.
- `flai check --strict` warns on TH-0032 in main's `wip/` (answered, its story archived), which is not this story's; `scripts/close-out.sh` stops on it, so the story's commits were made by hand after the verifier's run (TH-0056: a finding outside the story is a note).
