---
id: T-028
type: task
nature: feature
title: Fixture repo with expected numbers
status: done
parent: S-008
owner: alex
created: 2026-09-15T18:00:58Z
updated: 2026-09-15T18:09:35Z
transitions:
  - to: ready
    at: 2026-09-15T18:09:35Z
    by: agent
  - to: in-progress
    at: 2026-09-15T18:09:35Z
    by: agent
  - to: done
    at: 2026-09-15T18:09:35Z
    by: agent
stream: S-008
tags: [cli, test]
---

# T-028 Fixture repo with expected numbers

## Work
A fixture project under internal/metrics/testdata with a handful of stories and tasks whose transitions give hand-computed numbers; tests assert them exactly. The same fixture serves the dashboard's alignment test later.

## Done when
Expected numbers are written in the fixture README and asserted.

## Notes
