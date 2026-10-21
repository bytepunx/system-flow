---
id: T-1958
type: task
nature: feature
title: The dashboard's design and user docs describe the top bar and the nav rail in place of the site menu
status: in-progress
parent: S-0458
owner: alex
created: 2026-10-21T07:02:36Z
updated: 2026-10-21T08:14:11Z
transitions:
  - to: ready
    at: 2026-10-21T08:14:09Z
    by: agent-S-0458
  - to: in-progress
    at: 2026-10-21T08:14:11Z
    by: agent-S-0458
stream: S-0458
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-1956]
---
# T-1958 The dashboard's design and user docs describe the top bar and the nav rail in place of the site menu

## Work

- In `design/system/flaiover-dashboard.md` § Views, replace the site menu's paragraph with the rail and the top bar:
  - the five groups and their pages, from `SITE_MENU`;
  - expand and collapse, with the icon highlighted or dimmed;
  - the per-browser `localStorage` key;
  - the drawer below `md`;
  - where Overview, Updates, and License are reached;
  - `locate` and `groupBadges`.

  In the table, each route's mention of the site menu's group becomes the rail's group.
- In `docs/users/flaiover.md` § The header, describe the top bar's four parts and the rail, and how to open it on a phone. Each section heading that names a site-menu group names its rail group instead.
- Waits for T-1956, so that the docs describe what was built. It runs beside the Playwright test.

## Done when

Both documents name the rail and the top bar, nothing in them describes the two-tier site menu as current, and the markdown tier passes.

## Notes
