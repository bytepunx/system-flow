---
id: T-0541
type: task
nature: feature
title: flai serve records usage from its agents' logs, and flai serve agent usage records it by hand
status: done
parent: S-0143
owner: alex
created: 2026-09-29T06:08:15Z
updated: 2026-09-29T06:20:48Z
transitions:
  - to: ready
    at: 2026-09-29T06:15:48Z
    by: agent-S-0143
  - to: in-progress
    at: 2026-09-29T06:15:48Z
    by: agent-S-0143
  - to: done
    at: 2026-09-29T06:20:48Z
    by: agent-S-0143
stream: S-0143
tags: []
touches: [flai/internal/serve, flai/cmd, docs/users, docs/operators, design/system/flai-cli.md]
---
# T-0541 flai serve records usage from its agents' logs, and flai serve agent usage records it by hand

## Work

- When an agent flai serve started ends, flai serve measures its story from every log it keeps for that story, and each of the story's tasks by the windows it was in progress, writes their usage, and cascades to the epic.
- While the agent runs, a task of its story that enters done is measured at the next look, so the board shows it before the session ends.
- `flai serve agent usage [S-nnnn...]` prints what the logs say and, with `--write`, records it; `--all` covers every story with a log, to fill in stories worked before this.
- Agents flai serve starts are the only ones measured; the docs say so.

## Done when

- Tests: a run's end writes the story's and its tasks' usage and the epic's sum; a task done mid-run is measured once.
- `docs/users/flai.md`, `docs/users/flai-reference.md`, `design/system/flai-cli.md` describe it.
- `make test` and lint pass.

## Notes
