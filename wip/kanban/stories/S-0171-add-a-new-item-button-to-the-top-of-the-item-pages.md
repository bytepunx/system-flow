---
id: S-0171
type: story
nature: improvement
title: "Add a \"New {item}\" button to the top of the item pages"
status: ready
parent: E-0013
owner: alex
created: 2026-09-30T00:22:14Z
updated: 2026-09-30T00:23:01Z
transitions:
  - to: ready
    at: 2026-09-30T00:22:24Z
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
# S-0171 Add a "New {item}" button to the top of the item pages

## Goal

At the top of the item page, display a "New {item}" button to allow the operator to make another of the same type of work item they're currently viewing.

This prevents the operator from having to go back to the board to create a new item.

## Acceptance criteria
- [ ] The button is displayed at the top of the page
- [ ] The item type named in the button matches the item type being viewed
- [ ] The item type being created matches the type that was being viewed

## Tasks

## Notes
