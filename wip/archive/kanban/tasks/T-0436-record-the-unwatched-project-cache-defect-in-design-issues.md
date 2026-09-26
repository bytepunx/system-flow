---
id: T-0436
type: task
nature: remediation
title: Record the unwatched project cache defect in design/issues
status: done
parent: S-0117
owner: arobson
created: 2026-09-26T03:29:51Z
updated: 2026-09-26T03:33:21Z
transitions:
  - to: ready
    at: 2026-09-26T03:29:54Z
    by: agent-S-0117
  - to: in-progress
    at: 2026-09-26T03:31:37Z
    by: agent-S-0117
  - to: done
    at: 2026-09-26T03:33:21Z
    by: agent-S-0117
stream: S-0117
tags: [dashboard]
touches: [design/issues]
---
# T-0436 Record the unwatched project cache defect in design/issues

## Work

- Record the defect with `flai issue new`: a named project's cached answers were kept until the container restarted because only the default project's Repo watched.

## Done when

- The issue file exists with its instance and remediation, and `design/issues/summary.md` lists it.

## Notes
