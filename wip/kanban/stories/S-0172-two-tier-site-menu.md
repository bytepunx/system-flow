---
id: S-0172
type: story
nature: improvement
title: Two-tier Site Menu
status: ready
parent: E-0013
owner: alex
created: 2026-09-30T00:46:03Z
updated: 2026-09-30T00:46:03Z
transitions:
  - to: ready
    at: 2026-09-30T00:46:03Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
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
- [ ] All parent items are visible and hovering over, clicking, or touching them should result in the display of their child menu
- [ ] When a child item is clicked/touched, the site navigates to that page and the child and parent menus show as active
- [ ] If a child item would have an indicator, the parent should receive the same indicator (i.e. Inbox has a 3 next to it, so should the parent)

## Tasks

## Notes
