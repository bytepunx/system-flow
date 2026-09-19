---
title: Work management
updated: 2026-09-19
audience: agent
order: 30
status: active
---

# Work management

How work is pulled, sized, tracked, and finished. The board is `wip/kanban`; the rules behind it are in `design/system/workflow.md` and `design/system/work-hierarchy.md`.

## Rules

- Work is pulled, never pushed. Take the top `ready` story from the board's `order`; do not pick by preference.
- With the `flai` MCP server connected, `inbox` is how you learn what changed: it lists threads awaiting you, the stories ready to pull in pull order, and what others did to work items since you last looked. Call it at the start of every turn or session, at every task transition, and before moving a story to review. If you stay running while idle, hold `wait_for_events`; if you end your turn, `inbox` at your next start reports everything in between. When nothing is in progress and `can_pull` is true, pull the first ready story it lists without waiting to be told.
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
- WIP limits are in `wip/kanban/board.md`: ready 5, in-progress 2, review 3.
- The operator accepts stories; the agent runs `flai accept S-nnnn --by alex` on their word, which merges the story branch, moves to done, archives, commits, releases, and pushes in one step. Done means accepted: the operator moving a story to done, from the board or with `flai move`, runs the same flow and is sufficient on its own (S-0046). The agent moves stories to review and never to done.
- E-0006 (designer's workbench) is pulled in its story order; S-0036 (authentication) before anything that writes from the dashboard.
- Declare what a story or task changes with `--touches` or `flai touches`; check `flai board` for overlaps before editing a path another in-progress item touches.
- The MCP server (`flai mcp`, registered in `.mcp.json`) runs the installed `flai` from `PATH`, not the tree's; when it is older than the tree, `inbox` may lack ready work and changes (they came with S-0058), so also run `scripts/flai.sh board` at the start of a turn until it is upgraded.
- The operator sets the pull order from the board by dragging a card within `ready` or `backlog`, or with `flai order` (S-0057). Pull the story they name, or the first ready one as `inbox` and `flai board` list it; do not reorder their list, and if an order change seems to matter, say so rather than making it.
- Check the MCP inbox at the start of every turn, at every task transition, and before moving a story to review; hold `wait_for_events` when idle if the session stays running. Answer threads with `thread_reply` and ask the designer with `thread_open` rather than stopping to ask in the conversation when the designer is not watching.
