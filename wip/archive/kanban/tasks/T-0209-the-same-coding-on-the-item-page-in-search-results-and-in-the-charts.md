---
id: T-0209
type: task
nature: feature
title: The same coding on the item page, in search results, and in the charts
status: done
parent: S-0055
owner: alex
created: 2026-09-19T06:26:47Z
updated: 2026-09-19T06:32:53Z
transitions:
  - to: ready
    at: 2026-09-19T06:26:47Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T06:26:47Z
    by: system-flow
  - to: done
    at: 2026-09-19T06:32:53Z
    by: system-flow
stream: S-0055
tags: []
---

# T-0209 The same coding on the item page, in search results, and in the charts

## Work
Asked for by the operator at review. A small component shows an item's nature as a chip on its tint and its type beside a swatch of its stripe colour, from the maps in `$lib/cardcolour.ts`. The item page's header line uses it in place of the plain words. Search: the hit carries the item's `nature` (the index already reads it), and an item hit is tinted and striped like a board card; a document hit stays plain. Charts: the series already take a fixed colour per nature from `NATURE_SLOT`, in the same hue family as each tint, and marks need saturated colours that a pastel cannot give, so the series colours stay and a test holds each one within 25 degrees of hue of its tint in both themes; the data table under the cycle time charts shows the nature as the chip. Design document and user docs updated; browser check of all three in both themes.

## Done when
- The item page, item hits in search, and the charts' item table show the coding, from the same maps as the board
- A test ties each nature's chart colour to its tint's hue in both themes
- `make flaiover-test` and `make flaiover-build` pass, and the three pages were looked at in a browser in both themes

## Notes
