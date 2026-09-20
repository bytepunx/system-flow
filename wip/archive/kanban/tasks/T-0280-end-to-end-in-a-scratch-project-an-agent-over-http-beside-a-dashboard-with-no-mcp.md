---
id: T-0280
type: task
nature: feature
title: "End to end in a scratch project: an agent over HTTP beside a dashboard with no /mcp"
status: done
parent: S-0076
owner: alex
created: 2026-09-20T12:16:19Z
updated: 2026-09-20T12:38:07Z
transitions:
  - to: ready
    at: 2026-09-20T12:30:44Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:30:44Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:38:07Z
    by: system-flow
stream: S-0076
tags: []
---
# T-0280 End to end in a scratch project: an agent over HTTP beside a dashboard with no /mcp

## Work
Scratch project, own FLAI_CONFIG and FLAI_CACHE_DIR, throwaway token: flai mcp start, initialize, inbox, wait_for_events woken by a move made from a scratch dashboard, stop. The operator's dashboard and flai serve are not touched.

## Done when
- What was tried and what was not is in the narrative
- Story to review

## Notes
