---
id: T-0413
type: task
nature: remediation
title: A story moved from backlog to ready with room in progress is started, per the designer's answer on TH-0010
status: done
parent: S-0116
owner: alex
created: 2026-09-24T08:46:28Z
updated: 2026-09-24T09:00:42Z
transitions:
  - to: ready
    at: 2026-09-24T08:52:13Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T08:52:14Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:00:42Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flai/internal/serve]
---
# T-0413 A story moved from backlog to ready with room in progress is started, per the designer's answer on TH-0010

## Work

Pin criterion 1 with a behavior test: a story moved from backlog to ready while the in-progress limit has room is started at the next look, and is not while the limit is full. How long someone attending may still hold it back is the designer's answer on TH-0010: with (a), ADR-0042 as S-0114 delivers it; with (b), a new ADR superseding that hold and the launcher changed to match.

## Done when

- The test passes on the story branch with S-0114's launcher in it.
- TH-0010 is answered and what it decided is in the narrative's `## Decisions`.

## Notes
