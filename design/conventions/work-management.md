---
title: Work management
updated: 2026-09-15
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
- Do not write tasks for backlog stories. Refine when the story is next; detail written early is detail rewritten.
- Definition of ready, story: goal, acceptance criteria as checkboxes, at least one task, parent epic open. `flai move` enforces it.
- Definition of done, story: every criterion checked, every task done or cancelled, decisions recorded, docs updated, `flai check --strict` clean, narrative closed. Move the story to `review`; only the operator moves it to `done`.
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
- WIP limits are in `wip/kanban/board.md`: ready 5, in-progress 2, review 3.
- The operator accepts stories; the agent runs `flai accept S-nnnn --by alex` on their word, which moves to done, archives, commits, releases, and pushes in one step. The agent moves stories to review and never to done by hand.
- E-0006 (designer's workbench) is pulled in its story order; S-0036 (authentication) before anything that writes from the dashboard.
- Declare what a story or task changes with `--touches` or `flai touches`; check `flai board` for overlaps before editing a path another in-progress item touches.
- Once S-0039 lands, check the MCP inbox at session start, at every task transition, and before moving a story to review; hold `wait_for_events` when idle.
