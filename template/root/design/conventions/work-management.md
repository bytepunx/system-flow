---
title: Work management
updated: 2026-09-18
audience: agent
order: 30
status: active
---

# Work management

How work is pulled, sized, tracked, and finished. The board is `wip/kanban`; the rules behind it are in `design/system/workflow.md` and `design/system/work-hierarchy.md`.

## Rules

- Work is pulled, never pushed. Take the top `ready` story from the board's `order`; do not pick by preference.
- Respect WIP limits. If pulling would exceed the limit, finish or hand off something first. A limit breach that cannot be avoided is logged in the narrative.
- One story in progress per agent at a time. Tasks within it are worked in order.
- Every change belongs to a story. If there is no story for what you are about to do, create one under the right epic and refine it before starting; if it is a five-minute fix, make it a task on the current story.
- A story is one demonstrable increment sized for one to a few sessions. If it will not fit, split it before starting, not after.
- A task is done or not done. If a task needs sub-steps that could fail independently, it is several tasks.
- Tasks are written by the agent that starts the story, when it starts it. Do not write tasks for a backlog or ready story you are not about to work; detail written early is detail rewritten.
- Definition of ready, story: goal, acceptance criteria as checkboxes, parent epic open. Tasks are not part of ready. `flai move` enforces it.
- When you pull a story, in this order: move it to `in-progress`, open the narrative, read the goal, criteria, and notes, then write its tasks with `flai task new` if it has none, each with `## Work` and `## Done when`. Then work them in order. A story cannot go to `review` without at least one task.
- If the story does not say enough to write the tasks, do not invent scope. Record it with `flai block --reason`, open a thread on the story saying what is missing (`thread_open`, or `flai thread new`), and pull the next story.
- Definition of done, story: every criterion checked, at least one task and every task done or cancelled, decisions recorded, docs updated, `flai check --strict` clean, narrative closed. Move the story to `review`; only the operator moves it to `done`.
- Transition items as you go with `flai move`, at the moment the state changes, so timestamps are true. Never backfill a history.
- Keep the narrative current: rewrite `## Current state` and `## Next steps` after every task transition and before any long operation; append a log entry at every transition, decision, and blocker. `flai stream log` does the log.
- When blocked, record it with `flai block --reason`, put the question in `## Open questions`, and move to other work. Do not wait idle and do not guess past a blocker that only the operator can clear.
- Never mark a criterion checked that you did not verify. If verification is impossible in this environment, leave it unchecked and say why in the story notes.
- Acceptance triggers the mechanical follow-ups without further asking: archive the item with `flai archive`, commit, and release per `git.md`.

## When in doubt

- If a story is growing, stop and split it. If a task is growing, it was a story.
- If the operator's request does not fit the board, ask where it belongs before starting, and do the parts that need no answer.

<!-- system-flow:end-of-baseline -->

## Project additions
