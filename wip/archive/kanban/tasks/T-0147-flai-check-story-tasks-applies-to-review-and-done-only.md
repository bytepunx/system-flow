---
id: T-0147
type: task
nature: improvement
title: "flai check: story.tasks applies to review and done only"
status: done
parent: S-0049
owner: alex
created: 2026-09-18T18:08:30Z
updated: 2026-09-18T18:13:17Z
transitions:
  - to: ready
    at: 2026-09-18T18:12:00Z
    by: alex
  - to: in-progress
    at: 2026-09-18T18:12:00Z
    by: alex
  - to: done
    at: 2026-09-18T18:13:17Z
    by: alex
stream: S-0049
tags: []
touches: [flai/internal/check]
---

# T-0147 flai check: story.tasks applies to review and done only

## Work
In `flai/internal/check/check.go`, raise `story.tasks` only when the story is `review` or `done`; keep `story.criteria` for every status except backlog and cancelled. Add fixture cases: a ready story with no tasks and an in-progress story with no tasks produce no finding; a review story with no tasks produces `story.tasks`.

## Done when
- The check tests cover the three cases and pass
- `scripts/flai.sh check --strict` is clean on this repository

## Notes
