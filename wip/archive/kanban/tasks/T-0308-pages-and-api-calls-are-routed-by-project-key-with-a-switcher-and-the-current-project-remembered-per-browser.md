---
id: T-0308
type: task
nature: feature
title: Pages and API calls are routed by project key, with a switcher and the current project remembered per browser
status: done
parent: S-0080
owner: alex
created: 2026-09-20T21:09:24Z
updated: 2026-09-20T21:49:37Z
transitions:
  - to: ready
    at: 2026-09-20T21:30:43Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T21:30:43Z
    by: system-flow
  - to: done
    at: 2026-09-20T21:49:37Z
    by: system-flow
stream: S-0080
tags: []
---
# T-0308 Pages and API calls are routed by project key, with a switcher and the current project remembered per browser

## Work
Page routes move under a project key segment; a top-level page lists the projects flai serves and lets the designer choose, remembered in this browser (localStorage) and used when they return. The client's api() helper and the event stream carry the project key. Internal links updated throughout.

## Done when
- Component and route tests for the switcher, the remembered choice, and navigation between two projects
- The flaiover tests and build pass

## Notes
