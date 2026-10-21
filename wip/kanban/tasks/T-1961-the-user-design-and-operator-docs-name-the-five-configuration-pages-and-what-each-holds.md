---
id: T-1961
type: task
nature: improvement
title: The user, design, and operator docs name the five Configuration pages and what each holds
status: backlog
parent: S-0459
owner: alex
created: 2026-10-21T07:03:25Z
updated: 2026-10-21T07:03:25Z
transitions: []
stream: S-0459
tags: [dashboard]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md, docs/operators/settings.md]
after: [T-1960]
---
# T-1961 The user, design, and operator docs name the five Configuration pages and what each holds

## Work

- In `docs/users/flaiover.md` § Settings, replace the one page with the five Configuration pages. Name each page's sections and what each changes.
- In `design/system/flaiover-dashboard.md` § Views, give `/settings/:page` its row. Describe the redirect from `/settings` and its anchors, and `settingspages.ts`.
- In `docs/operators/settings.md`, each mention of the dashboard's Settings page names the Configuration page that holds the setting: Host actions under Features, Shared paths under File System, and so on.
- Waits for T-1960, so that the docs name the routes as built.

## Done when

No document in the three describes Settings as one page, and the markdown tier passes.

## Notes
