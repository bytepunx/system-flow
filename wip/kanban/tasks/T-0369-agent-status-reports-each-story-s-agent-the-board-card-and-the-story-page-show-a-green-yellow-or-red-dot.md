---
id: T-0369
type: task
nature: feature
title: "agent.status reports each story's agent; the board card and the story page show a green, yellow, or red dot"
status: done
parent: S-0104
owner: alex
created: 2026-09-23T17:25:01Z
updated: 2026-09-23T17:39:14Z
transitions:
  - to: ready
    at: 2026-09-23T17:33:41Z
    by: system-flow
  - to: in-progress
    at: 2026-09-23T17:33:41Z
    by: system-flow
  - to: done
    at: 2026-09-23T17:39:14Z
    by: system-flow
stream: S-0104
tags: []
---

# T-0369 agent.status reports each story's agent; the board card and the story page show a green, yellow, or red dot

## Work
`agent.status` gains `stories`: for each story with an agent run, `working`, `waiting` (running, with a thread on the story whose last entry is not the designer's, or the story blocked), `failed` (the reason, the exit code, the log), or `done`. The board card shows a dot (green working, yellow waiting, red failed) with a title that says why. The story page shows the same dot with the harness, the model, the time it started, and for a failure the reason and the log's name. Both refresh on file events and at a slow poll while an agent runs.

## Done when
- flai tests for the per-story states.
- flaiover tests for the dot on a card and on the story page in each state.

## Notes
