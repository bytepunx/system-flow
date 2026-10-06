---
id: T-1128
type: task
nature: improvement
title: Brighter nature tints for the shown state, in both themes
status: backlog
parent: S-0302
owner: alex
created: 2026-10-06T23:19:53Z
updated: 2026-10-06T23:19:53Z
transitions: []
stream: S-0302
tags: [dashboard]
touches: [flaiover/src/routes/layout.css, flaiover/src/lib/cardcolour.ts, flaiover/src/lib/theme.test.ts]
---
# T-1128 Brighter nature tints for the shown state, in both themes

## Work

Add a brighter companion to each nature tint, for a legend tag whose nature is shown: tokens `--t-nature-<nature>-shown` for feature, improvement, remediation, research, and experiment in the light and the dark theme of `flaiover/src/routes/layout.css`, mapped as `--color-nature-<nature>-shown` beside the existing nature colours. Each is a more saturated (light theme) or lighter (dark theme) version of the nature's pastel, so the pair reads as one hue, dim and lit. Add a map `natureShownTint` in `flaiover/src/lib/cardcolour.ts`, written out in full as `natureTint` is, so Tailwind generates the classes, and say in its comment that the legend reads it (S-0302). The card backgrounds keep the existing pastels.

Extend `flaiover/src/lib/theme.test.ts`: ink text on each shown tint meets the contrast the nature tints are held to, and the five shown tints are distinct from each other and from their default tints. Waits for nothing: it runs with the filter store task.

## Done when

- [ ] Both themes define a shown tint for every nature, and `natureShownTint` names a class for each
- [ ] `theme.test.ts` checks ink contrast on every shown tint and that each differs from its default tint, and passes
- [ ] `scripts/flaiover-test.sh` passes

## Notes
