---
id: T-1572
type: task
nature: improvement
title: flai stream open starts a ready story as flai story start does, moving it to in-progress and refusing one the board holds
status: in-progress
parent: S-0304
owner: alex
created: 2026-10-14T20:56:19Z
updated: 2026-10-14T20:56:28Z
transitions:
  - to: ready
    at: 2026-10-14T20:56:27Z
    by: agent-S-0304
  - to: in-progress
    at: 2026-10-14T20:56:28Z
    by: agent-S-0304
stream: S-0304
tags: []
touches: [flai/internal/storystart/start.go, flai/internal/storystart/start_test.go, flai/cmd/stream.go, flai/cmd/stream_open_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
---
# T-1572 flai stream open starts a ready story as flai story start does, moving it to in-progress and refusing one the board holds

## Work

I-0097: `flai stream open` on a ready story opens its narrative and worktree but leaves it in ready, so `permission_prompt` and `flai guard` refuse the agent's protected write ("is ready, not in progress") until it moves the story by hand. `flai story start` folds the move in, but `flai stream open` is still documented and still named by many refusals ("open it with flai stream open").

Make `flai stream open` on a story in `ready` start it as `flai story start` does: refuse it, changing nothing, when the board holds it (overlap or `after`), with the same reason; otherwise move it to `in-progress` as the agent (`FLAI_AGENT`), its epic following, and then open the narrative and the branch. Share the refusal and the move with `internal/storystart` rather than copy them. Every other state is opened as before. `--json` says whether it moved the story. Record the decision in a new ADR. Update `flai-cli.md`, `docs/users/flai.md`, and the generated reference.

Waits for nothing: it is the first layer.

## Done when

- A test opens a ready story with `flai stream open` and finds it in progress afterwards, and fails against main's `flai stream open`.
- A test finds a held ready story refused with nothing written.
- The design, user docs, and reference say so, and an ADR records it.

## Notes
