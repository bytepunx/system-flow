---
title: Session start
updated: 2026-09-18
audience: agent
order: 10
status: active
---

# Session start

What to load before touching anything, how to resume after a crash, and how to leave the repository so the next session can continue.

## Rules

- Read, in order, before any change: this folder as listed in `README.md`, then `CLAUDE.md`, then `wip/agents/index.md`, then `wip/kanban/board.md`.
- If `wip/agents/index.md` lists an active stream, read that narrative's `## Current state` and `## Next steps` before anything else. Assume it is true until `git status` says otherwise.
- Run `git status` and compare with the narrative. Uncommitted changes the narrative does not mention are the first thing to reconcile, and the reconciliation goes in the log.
- Append a log entry to the narrative stating that a session started or recovered, and what state you found.
- If no stream is active, pull the top story from the board's `order` that is `ready`, respecting the WIP limit, and open its narrative with `flai stream open`. Do not start a `backlog` story; refine it to `ready` first and say so. A ready story may have no tasks: writing them is the first thing you do once it is `in-progress`, as `work-management.md` describes.
- Do not re-derive facts already recorded in the narrative, the story, or `design/system`. Read them.
- Do not read the whole repository to orient yourself. The four documents above plus the active story and its design links are enough; go wider only when a task needs it.
- Before any long-running or risky operation, rewrite `## Current state` and `## Next steps` so a crash mid-operation loses nothing.
- End every session with a closing log entry, a true `## Current state`, an ordered `## Next steps` whose first item is the very next action, and `flai check --strict` passing. Leave nothing half-edited that the narrative does not describe.

## When in doubt

- If the narrative and the repository disagree, the repository is the fact and the narrative is corrected, with a log entry saying what was wrong.
- If you cannot tell whether to continue a stream or start a new one, continue the stream and ask in `## Open questions`.

<!-- system-flow:end-of-baseline -->

## Project additions
