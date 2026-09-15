---
id: T-027
type: task
nature: feature
title: "Stats: burn-up and cumulative flow series"
status: done
parent: S-008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:34Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:34Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:34Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:34Z
    by: agent
stream: S-008
tags: [cli, stats]
---

# T-027 Stats: burn-up and cumulative flow series

## Work
Daily burn-up (scope and done, whole set and per parent) and cumulative flow (count per state per day) series so the dashboard can chart without re-deriving.

## Done when
Series are in the --json output and tested on the fixture.

## Notes
