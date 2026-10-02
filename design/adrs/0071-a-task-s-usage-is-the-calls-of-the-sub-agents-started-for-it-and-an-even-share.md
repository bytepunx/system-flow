---
id: ADR-0071
title: "A task's usage is the calls of the sub-agents started for it and an even share of the story's agent's calls while it was in progress"
status: accepted
date: 2026-10-02
supersedes: []
superseded_by: []
refines: [ADR-0051]
---

# ADR-0071 A task's usage is the calls of the sub-agents started for it and an even share of the story's agent's calls while it was in progress

## Context

ADR-0051 measures a task over the intervals it was in progress: each session's reported totals in the share of its calls' input and cache tokens that fall in them. That assumes one task at a time. S-0176 worked two layers of two tasks at once, each by its own task sub-agent, and each task of a layer got the whole window: T-0679 and T-0680 were each given 5.08 US dollars, and the six tasks summed to about 22 dollars against the run's 17.79 (I-0054).

Claude Code's stream-json says which sub-agent made a call. A sub-agent's calls are in its story's agent's session, and each carries `parent_tool_use_id`, the ID of the `Agent` (or `Task`) `tool_use` that started it, and `task_description`, that call's description. The story's agent's own calls carry none. In S-0176's log, 355 of 543 calls were its 14 sub-agents', and every task sub-agent's description named its task ("Task sub-agent: T-0677 after").

## Decision

A task is given each session's reported totals in the share of their calls' input and cache tokens that are its own: the calls of the sub-agents started for it, and an even share of the story's agent's calls made while it was in progress.

1. **Its sub-agents' calls, wherever they fall.** A sub-agent is a task's when its `Agent` call's description names one of the story's task IDs (the first it names), or, when the description names none, its prompt does (the first of the story's task IDs the prompt names). `task_description` stands in for the description when the call was not seen.
2. **An even share of every other call made while it was in progress.** The story's agent's own calls, and those of sub-agents started for no task, are split evenly among the tasks in progress when each was made.

Calls no result reported are estimated, and shared, the same way. A task's seconds are still the time it was in progress while a run went on, and a story's usage is unchanged. The story's agent names the task's ID in each task sub-agent's description, as `delegation.md` and the `claude-code` prompt say, so that the match does not rest on the prompt's wording.

## Consequences

- Tasks worked at once no longer each get the run: what they are given sums to no more than the story.
- A task worked alone with no sub-agent gets what ADR-0051 gave it.
- A sub-agent's work is counted to its task even when it starts before the task moves to in-progress or ends after it leaves, and a task never moved to in-progress whose sub-agents worked is measured, with no seconds.
- A sub-agent whose description and prompt name no task, or name another story's, is shared like the story's agent's own calls. An agent that names the wrong task moves the usage with it.
- A sub-agent started by a sub-agent is matched by its own description and prompt, not by its parent's.

## Alternatives considered

- **Split each shared window evenly, without reading sub-agents.** Simpler, and it stops the double count, but it gives a long task and a short one in the same layer the same share.
- **Count only a task's sub-agents.** It leaves out the story's agent's reviewing and committing of the task, which is part of what the task cost.
