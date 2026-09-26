---
id: T-0447
type: task
nature: improvement
title: The switcher tells a served project that is not connected apart from a connected one, and links to the settings page
status: done
parent: S-0122
owner: alex
created: 2026-09-26T06:05:09Z
updated: 2026-09-26T06:28:52Z
transitions:
  - to: ready
    at: 2026-09-26T06:05:19Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T06:12:43Z
    by: agent-S-0119
  - to: done
    at: 2026-09-26T06:28:52Z
    by: agent-S-0122
stream: S-0122
tags: []
touches: [flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/routes/api/projects, flaiover/src/lib/project.svelte.ts]
---
# T-0447 The switcher tells a served project that is not connected apart from a connected one, and links to the settings page

## Work

`/api/projects` adds the projects flai serve serves that the dashboard has no connection to, from `settings.get`, marked served and not connected. The switcher labels a served project that is not connected differently from a connected one, and a link beside it goes to `/settings`, where the project's last error says why.

## Done when

- route and component tests cover a connected, a served but disconnected, and a candidate project
- `make flaiover-test` passes

## Notes
