---
id: ADR-0153
title: "flai stream open starts a ready story as flai story start does, moving it to in-progress and refusing one the board holds"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0019, ADR-0086]
---

# ADR-0153 flai stream open starts a ready story as flai story start does, moving it to in-progress and refusing one the board holds

## Context

`flai stream open` opens a story's narrative and checks its branch out in a worktree (ADR-0019). It never moved the story. An agent that opened a ready story with it, as `flai serve`'s start prompt once told it to, was left working a story still in ready. `permission_prompt` (ADR-0086) and `flai guard` ask the operator for a protected write only in an in-progress story's worktree, so the agent's first write under `.claude/` was refused with "S-nnnn is ready, not in progress" until it moved the story itself (I-0097).

`flai story start` (S-0274) folds the move, the stream, the prime, and the inbox into one call, and `flai serve`'s prompt now names it. `flai stream open` stays documented, and many refusals still say "open it with flai stream open", so an agent can still take the old path.

## Decision

`flai stream open` on a story in `ready` starts it as `flai story start` does, with the same code: it refuses a story the board holds, by an overlap with a story in progress or in review or by `after`, with the same reason and nothing changed; otherwise it moves the story to `in-progress` as the agent (`FLAI_AGENT`), its epic following, and then opens the narrative and the branch. A story in any other state is opened as before, unmoved. `--json` gives `moved`, true when it moved the story.

## Consequences

- Opening a ready story's stream is working it: no agent is left with a story in ready and its worktree open, so a protected write is asked for as in any in-progress story.
- `flai stream open` on a ready story is held as `flai story start` is, so it no longer bypasses the board's holds.
- A ready story can still be opened without a move only by moving it back; nothing in the workflow needs that.

## Alternatives considered

- **Refuse a ready story, naming `flai story start`.** Fixes the cause too, but turns an opening that works today into an error and needs a second call.
- **Let `permission_prompt` ask for a ready story's worktree.** Leaves the story in ready while it is worked, so the board, the WIP limit, and the metrics stay wrong.
- **Only reword the refusal.** Saves minutes per occurrence but leaves the cause.
