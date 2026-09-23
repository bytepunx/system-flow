---
id: S-0101
type: story
nature: improvement
title: Flai serve, dashboard, and mcp commands should work in parent folders without system-flow.yaml
status: done
owner: alex
created: 2026-09-23T05:18:13Z
updated: 2026-09-23T05:46:03Z
transitions:
  - to: ready
    at: 2026-09-23T05:21:25Z
    by: alex
  - to: in-progress
    at: 2026-09-23T05:21:53Z
    by: system-flow
  - to: review
    at: 2026-09-23T05:41:56Z
    by: system-flow
  - to: done
    at: 2026-09-23T05:46:03Z
    by: alex
tags: [flai]
touches: [flai/cmd]
---

# S-0101 Flai serve, dashboard, and mcp commands should work in parent folders without system-flow.yaml

## Goal

## Acceptance criteria
- [x] flai serve, flai mcp, and flai dashboard should all successfully start, even in an empty directory
- [x] missing system-flow.yaml should not prevent flai processes from working correctly

## Tasks
- T-0356 flai dashboard runs outside a project: it starts or reports the shared dashboard and flai serve, and names the folder for import
- T-0357 flai mcp in a folder serves every system-flow project below it
- T-0358 Tried in an empty folder and a parent folder; ADR and docs

## Notes
