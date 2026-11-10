---
id: D-0554
title: A story's agent works a task itself, not through a task sub-agent, when the task is a few lines in a file it has already read, since a hand-off would cost more reading than the edit
kind: decision
sources: [S-0649, S-0650, S-0653]
topics: [delegation, code, testing]
paths: [design/conventions/delegation.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0554 A story's agent works a task itself, not through a task sub-agent, when the task is a few lines in a file it has already read, since a hand-off would cost more reading than the edit

`delegation.md` § Tasks worked by sub-agents says to hand each task of the plan to a task sub-agent. Three remediations accepted on 2026-11-10 each worked a task themselves instead, and said why in their narratives' `## Decisions`:

```text
S-0649: T-2834 was worked here, not by a sub-agent: it is one file I had already read, so a hand-off would have been slower (delegation.md).
S-0650: T-2836 worked by the story's agent, not a task sub-agent: one small edit to one test file, quicker done here.
S-0653: T-2849 and T-2851 were worked by the story's agent, not task sub-agents: each is a few lines in files already read while confirming the planner's cause, so a sub-agent would have cost more reading than the edit.
```

Each was a load-dependent test fix whose cause the story's agent had already read while confirming the planner's plan. So a story's agent whose task is a small edit to files already in its context may work it itself, and records the exception in its Decisions as these did. D-0523 is the related case of a fresh sub-agent in place of a fork, and D-0480 the case of work outside the worktree that sub-agents cannot reach.

Sources: S-0649, S-0650, S-0653 (`wip/archive/agents/<story>.md`, Decisions).
