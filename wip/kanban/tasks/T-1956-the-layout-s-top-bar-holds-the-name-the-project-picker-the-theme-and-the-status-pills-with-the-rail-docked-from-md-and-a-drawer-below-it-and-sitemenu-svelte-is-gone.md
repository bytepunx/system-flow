---
id: T-1956
type: task
nature: feature
title: The layout's top bar holds the name, the project picker, the theme, and the status pills, with the rail docked from md and a drawer below it, and SiteMenu.svelte is gone
status: backlog
parent: S-0458
owner: alex
created: 2026-10-21T07:02:12Z
updated: 2026-10-21T07:02:12Z
transitions: []
stream: S-0458
tags: [dashboard]
touches: [flaiover/src/routes/+layout.svelte, flaiover/src/routes/layout.svelte.test.ts, flaiover/src/routes/layout.css, flaiover/src/lib/components/SiteMenu.svelte]
after: [T-1955]
---
# T-1956 The layout's top bar holds the name, the project picker, the theme, and the status pills, with the rail docked from md and a drawer below it, and SiteMenu.svelte is gone

## Work

- Rewrite the `<header>` in `flaiover/src/routes/+layout.svelte` to hold, in order:
  - below `md`, the menu button;
  - the flaiover name, linking to `/`;
  - `ProjectSwitcher`;
  - the theme button;
  - `HostFlai`.

  `HostFlaiBanner` and `UnsignedBanner` stay below the header.
- Lay out the body as the rail and the main area. From `md`, `NavRail` is docked at the left and sticky under the header. Below `md` it is a drawer that the menu button opens, with `aria-expanded` and `aria-controls` on the button. Escape or a chosen page (`onnavigate`) closes it, and focus returns to the button.
- Keep `max-w-6xl` on the main area's content by default, so pages read as before beside the rail.
- Delete `flaiover/src/lib/components/SiteMenu.svelte`, and any test only it used. Put any rail width or drawer style that the utilities cannot express in `layout.css`.
- Write `layout.svelte.test.ts`: the header holds exactly the four parts, and the menu button toggles the drawer and closes it on Escape and on navigation.
- Waits for T-1955, the rail it places.

## Done when

`layout.svelte.test.ts` and the whole vitest tier pass, nothing imports `SiteMenu.svelte`, and every page shows the rail and the new top bar.

## Notes
