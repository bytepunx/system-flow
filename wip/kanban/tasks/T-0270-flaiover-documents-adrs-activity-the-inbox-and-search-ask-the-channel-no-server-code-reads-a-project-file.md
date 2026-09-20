---
id: T-0270
type: task
nature: feature
title: "flaiover: documents, ADRs, activity, the inbox, and search ask the channel; no server code reads a project file"
status: done
parent: S-0074
owner: alex
created: 2026-09-20T08:30:00Z
updated: 2026-09-20T08:47:33Z
transitions:
  - to: ready
    at: 2026-09-20T08:37:40Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:37:40Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:47:33Z
    by: system-flow
stream: S-0074
tags: []
---
# T-0270 flaiover: documents, ADRs, activity, the inbox, and search ask the channel; no server code reads a project file

## Work
repo.ts, search.ts, inbox.ts, and activity.ts ask flai; the front matter parser, the mtime cache, minisearch, and yaml leave; a test fails if server code imports node:fs for project files; large answers fit the message cap or are paged.

## Done when
- Unit tests with a fake Ask and through flai hostapi
- make flaiover-test and make flaiover-build pass

## Notes
