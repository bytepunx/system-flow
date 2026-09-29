---
id: T-0508
type: task
nature: feature
title: Replay the pack on S-0125 and two archived stories and attach the outputs to the story notes
status: done
parent: S-0137
owner: alex
created: 2026-09-29T00:36:08Z
updated: 2026-09-29T00:46:26Z
transitions:
  - to: ready
    at: 2026-09-29T00:36:14Z
    by: agent-S-0137
  - to: in-progress
    at: 2026-09-29T00:43:16Z
    by: agent-S-0137
  - to: done
    at: 2026-09-29T00:46:26Z
    by: agent-S-0137
stream: S-0137
tags: [cli]
touches: [wip/kanban/stories]
---
# T-0508 Replay the pack on S-0125 and two archived stories and attach the outputs to the story notes

## Work

- Run `flai prime --story` for S-0125 and two archived stories that named ADRs.
- Attach each output's size and its selection to the story notes; check every ADR the story named that existed when it was created is printed or in the catalog.

## Done when

- The notes carry the three runs and the check.

## Notes
