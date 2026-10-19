---
id: S-0449
type: story
nature: feature
title: Improvements to Design Components Screen (Part 1)
status: backlog
owner: alex
created: 2026-10-19T04:51:01Z
updated: 2026-10-19T04:51:01Z
transitions: []
tags: []
touches: [flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  inputs:
    penalty_per_week: 10
    by: alex
    at: 2026-10-19T04:51:01Z
---
# S-0449 Improvements to Design Components Screen (Part 1)

## Goal

Currently the Components page needs a few usability improvements for it to be truly useful. As it stands, the preview iFrame centers components, causing them to render outside the visible pane. The JSON required doesn't provide any kind of schema - an empty schema or default schema would be much more useful. Sometimes the incorrect component is rendering in the iframe panel. If there is a component render issue, there are times that changing the ligh/dark/both option or size option can trigger at least a partial re-render.

## Acceptance criteria
- [ ] Only put a 10px margin around the rendered component
- [ ] Automatically size the pane in which the iFrame renders so that the entire component is visible on screen OR is visible via scroll bars to see over-flow content.
- [ ] Always pre-populated the JSON schema for json prop fields even if it's all default/empty values.
- [ ] Make sure that the selected component is all that is rendering, never render a prior component
- [ ] None of the component interactions should result in 500 errors - when 500 errors occur, file these as issues

## Tasks

## Notes
