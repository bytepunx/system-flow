---
id: T-1955
type: task
nature: feature
title: NavRail.svelte draws the five groups, expanding and collapsing each with its icon highlighted or dimmed, and keeps the choice per browser
status: done
parent: S-0458
owner: alex
created: 2026-10-21T07:01:57Z
updated: 2026-10-21T08:08:15Z
transitions:
  - to: ready
    at: 2026-10-21T08:03:30Z
    by: agent-S-0458
  - to: in-progress
    at: 2026-10-21T08:03:32Z
    by: agent-S-0458
  - to: done
    at: 2026-10-21T08:08:15Z
    by: agent-S-0458
stream: S-0458
tags: [dashboard]
touches: [flaiover/src/lib/components/NavRail.svelte, flaiover/src/lib/components/NavRail.svelte.test.ts, flaiover/src/lib/navrail.svelte.ts, flaiover/src/lib/navrail.svelte.test.ts, design/components/flaiover/src/lib/components/NavRail.examples.json]
after: [T-1953, T-1954]
usage:
  source: log
  seconds: 283
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 52
      output: 367
      cache_read: 2152328
      cache_write: 80172
      cost: 0.8297
---
# T-1955 NavRail.svelte draws the five groups, expanding and collapsing each with its icon highlighted or dimmed, and keeps the choice per browser

## Work

- Write `flaiover/src/lib/navrail.svelte.ts`, the expanded groups as rune state. It reads and writes `localStorage` under `flaiover-rail-groups`, as `$lib/boardtypes.svelte.ts` keeps the board's types, and always expands the group that holds the current page.
- Write `flaiover/src/lib/components/NavRail.svelte` from `SITE_MENU`:
  - Each group is a button with `aria-expanded`. Its icon is `text-primary` when expanded and `text-muted` when collapsed, with `groupBadges` beside it.
  - Each page is a link with its icon, its label, and `aria-current="page"` when current.
  - Workflow is drawn as coming soon, with `aria-disabled` and no link.
  - The component takes an `onnavigate` callback, so the layout's drawer can close when a page is chosen.
  - Build every colour and size from the token utilities.
- Test the expand and collapse, the persistence, the current page and its group, the badges on a collapsed group, and the coming-soon entry.
- Write `NavRail.examples.json` with a set for each of these: all groups collapsed, one group expanded, and a page with a badge.
- Waits for T-1953 (the icons it draws) and T-1954 (the data it reads).

## Done when

`NavRail.svelte.test.ts` and `navrail.svelte.test.ts` pass, and `flai check` reports no `components.examples` finding for it.

## Notes
