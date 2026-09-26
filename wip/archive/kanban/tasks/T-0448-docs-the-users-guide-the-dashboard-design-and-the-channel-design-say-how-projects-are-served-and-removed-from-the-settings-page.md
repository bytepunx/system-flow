---
id: T-0448
type: task
nature: improvement
title: "Docs: the users guide, the dashboard design, and the channel design say how projects are served and removed from the settings page"
status: done
parent: S-0122
owner: alex
created: 2026-09-26T06:05:09Z
updated: 2026-09-26T06:33:34Z
transitions:
  - to: ready
    at: 2026-09-26T06:05:20Z
    by: agent-S-0119
  - to: in-progress
    at: 2026-09-26T06:32:37Z
    by: agent-S-0122
  - to: done
    at: 2026-09-26T06:33:34Z
    by: agent-S-0122
stream: S-0122
tags: []
touches: [docs/users, design/system/flaiover-dashboard.md, design/system/dashboard-host-channel.md, docs/operators]
---
# T-0448 Docs: the users guide, the dashboard design, and the channel design say how projects are served and removed from the settings page

## Work

`docs/users/flaiover.md` (Settings, More than one project), `design/system/flaiover-dashboard.md` (the settings route, the switcher), `design/system/dashboard-host-channel.md` (the two new methods), and the operators' guide where a project is added or removed.

## Done when

- every file above says what the page and the switcher do
- `flai check --strict` and `scripts/lint-md.sh` pass

## Notes
