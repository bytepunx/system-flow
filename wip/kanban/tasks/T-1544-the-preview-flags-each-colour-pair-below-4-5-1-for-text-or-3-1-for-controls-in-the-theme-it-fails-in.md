---
id: T-1544
type: task
nature: feature
title: "The preview flags each colour pair below 4.5:1 for text or 3:1 for controls, in the theme it fails in"
status: backlog
parent: S-0383
owner: alex
created: 2026-10-13T17:47:30Z
updated: 2026-10-13T17:47:30Z
transitions: []
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/lib/design/preview.ts, flaiover/src/lib/design/preview.test.ts]
---
# T-1544 The preview flags each colour pair below 4.5:1 for text or 3:1 for controls, in the theme it fails in

## Work

Extend `flaiover/src/lib/design/preview.ts` with the pairs the gallery draws, text on ground and surface and controls against their background, and a check that returns each pair under 4.5:1 for text or 3:1 for controls with its ratio and theme, as `theme.test.ts` holds the dashboard's own palette. Tests cover a passing and a failing pair in each theme. It waits for nothing.

## Done when

- `preview.test.ts` passes with the new cases.

## Notes

- Layer 1, beside T-1543.
