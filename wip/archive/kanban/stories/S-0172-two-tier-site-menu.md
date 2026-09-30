---
id: S-0172
type: story
nature: improvement
title: Two-tier Site Menu
status: done
parent: E-0013
owner: alex
created: 2026-09-30T00:46:03Z
updated: 2026-09-30T01:25:31Z
transitions:
  - to: ready
    at: 2026-09-30T00:46:03Z
    by: alex
  - to: in-progress
    at: 2026-09-30T01:17:42Z
    by: agent-S-0172
  - to: review
    at: 2026-09-30T01:25:12Z
    by: agent-S-0172
  - to: done
    at: 2026-09-30T01:25:31Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 471
  models:
    - model: claude-opus-5-5
      input: 118
      output: 34824
      cache_read: 5999221
      cache_write: 132214
      cost: 2.9545
---
# S-0172 Two-tier Site Menu

## Goal

The current site menu (Overview, Board, Inbox, Activity, Charts, Docs, ADRs, Search, Host, and Settings) would benefit from a two-tier system with the following top level and child items:

- Workflow
  - Overview
  - Board
  - Inbox
  - Activity
- Status
  - Charts
  - ADRs
  - Docs
  - Search
- Host
  - Updates
  - Settings

When a parent item is hovered over, it becomes "active" (bold lettering, emphasis coloring) and a low of child menu options appear below it. When a child menu item is clicked, the site navigates to that page and then both the parent item and child item are "active" and the child menu stays visible.

If a child page (i.e. Inbox) has an indicator associated with it, the parent item should show the same indicator to make it clear there is an indicator on one of the child menu items.

## Acceptance criteria
- [x] All parent items are visible and hovering over, clicking, or touching them should result in the display of their child menu
- [x] When a child item is clicked/touched, the site navigates to that page and the child and parent menus show as active
- [x] If a child item would have an indicator, the parent should receive the same indicator (i.e. Inbox has a 3 next to it, so should the parent)

## Tasks
- T-0608 The header's navigation is a two-tier site menu
- T-0609 The user guide and the dashboard design describe the two-tier menu

## Notes

- Updates is the existing `/host` page, where the dashboard and flai are upgraded; no new page was made.
- Verified by `SiteMenu.svelte.test.ts` and `sitemenu.test.ts`, and in a browser against a scratch dev server: hover and click at 1100 px, a tap at 390 px with no sideways scroll. The parent's Inbox count was verified by the component test only, since the scratch server had no host flai to give an inbox.
