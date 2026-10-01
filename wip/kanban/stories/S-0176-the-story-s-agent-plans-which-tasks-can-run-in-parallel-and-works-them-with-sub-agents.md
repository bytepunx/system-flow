---
id: S-0176
type: story
nature: experiment
title: The story's agent plans which tasks can run in parallel and works them with sub-agents
status: ready
owner: alex
created: 2026-10-01T07:17:38Z
updated: 2026-10-01T07:38:54Z
transitions:
  - to: ready
    at: 2026-10-01T07:38:54Z
    by: alex
tags: [flai, template]
topics: [conventions]
touches: [flai/internal/workitem, flai/internal/check, flai/internal/harness, flai/internal/mcpserver, flai/cmd, flaiover/src, template/, design/conventions/, design/system/work-hierarchy.md, design/system/agent-context.md, design/system/agent-coordination.md]
after: [S-0175]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0176 The story's agent plans which tasks can run in parallel and works them with sub-agents

## Goal

Within a story, one agent works the tasks one after another, even when they touch disjoint paths and do not depend on each other. Parallelism exists only across stories (the in-progress limit and claims, ADR-0046). This experiment tests whether a story finishes faster, at an acceptable cost, when its agent plans the tasks for parallel work and hands independent ones to sub-agents.

The story's agent owns the plan. When it writes a story's tasks it decides, for each, whether it can start at once or must wait for a named task to be done, and records that in the task's front matter. It then runs the tasks whose dependencies are done, in parallel where their `touches` do not overlap, each with a sub-agent, and integrates, commits, and keeps the narrative itself. Forked sub-agents, which inherit the parent's conversation and prompt cache, are the first candidate, so a task's sub-agent starts already primed.

The hypothesis: for stories with three or more tasks of which at least two are independent, wall-clock time to review falls by a third or more, cost rises by no more than half, and defects found at review do not rise.

## Acceptance criteria
- [ ] Tasks carry `after: [T-nnnn]`, the tasks of the same story that must be done before it starts, set with `flai task new --after` and `flai edit --after`, and the MCP `item_new` and `item_edit` equivalents. `flai check` reports an entry that names no task, a task of another story, the task itself, and every cycle. Front matter stays strict: the design says which flai the host must run first
- [ ] `flai show`, the board, and the story's page in flaiover show each task as ready to start, waiting with the tasks it waits for, or in progress, and the plan as layers of tasks that can run at once
- [ ] The conventions and `harness.Prompt` for `claude-code` tell the story's agent to plan when it writes the tasks: what each task touches, which must wait for which and why, and which can run together (an `after` and no overlap in `touches`), and to record the plan's reasoning in the narrative's `## Decisions`
- [ ] The story's agent runs ready, non-overlapping tasks in parallel with sub-agents, waits for them, reviews their work, and moves each task; only it commits, syncs the stream, moves items, and talks to the designer. A sub-agent's question for the designer reaches the designer by the mechanism S-0175 chose
- [ ] Parallel sub-agents do not corrupt each other's work in the shared worktree: either each task gets its own worktree from the story's branch and the parent merges it back, or the experiment shows non-overlapping `touches` are enough in one worktree. Record which and the evidence
- [ ] At least four stories run with the plan and four comparable ones without it, comparing wall-clock time to review, cost and tokens (ADR-0051, sub-agents included), the number of parallel layers, conflicts, and review defects. The results and a recommendation (adopt, adapt, or drop) go in `design/system/agent-context.md`, with an ADR if adopted
- [ ] `design/system/work-hierarchy.md` documents a task's `after`

## Tasks

## Notes

- Depends on S-0175: sub-agents that cannot move items or consume the inbox, and the decision on how they raise questions.
- A task's `after` mirrors a story's (ADR-0046, S-0130), within a story only. Whether flai should also hold a task (refuse or warn on `flai move T-nnnn in-progress` while its `after` is not done) is part of the experiment; start with a warning, as for stories.
- Risks to watch: parallel edits to one worktree, the narrative and `flai stream sync` assuming one writer, and Anthropic's estimate that multi-agent runs use about fifteen times a chat's tokens, which ADR-0049 cited in rejecting a scout sub-agent. Forks that share the prompt cache may change that figure; measure it.
