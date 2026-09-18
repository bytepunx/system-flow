---
id: S-0044
type: story
nature: improvement
title: Dashboard theme from the brand palette
status: backlog
parent: E-0003
owner: alex
created: 2026-09-18T05:51:08Z
updated: 2026-09-18T05:51:08Z
transitions: []
tags: [dashboard]
---

# S-0044 Dashboard theme from the brand palette

## Goal
flaiover wears the brand palette instead of Tailwind defaults: `#0f1108`, `#23b5d3`, `#241909`, `#645853`, `#054a91`, applied through design tokens for a light and a dark theme, with charts and status colours derived from the same system and checked for contrast and colour-vision safety.

## Acceptance criteria
- [ ] A token layer (CSS custom properties consumed by Tailwind) defines surface, ink, muted, primary, accent, border, and focus for light and dark: `#0f1108` as the dark ground and light-theme ink, `#241909` as the dark raised surface, `#645853` as muted text and borders, `#054a91` as primary (buttons, active navigation, headings), `#23b5d3` as accent (links, focus rings, live indicators); the light ground is a warm off-white derived from the palette, not pure white
- [ ] Every text and interactive colour pair meets WCAG AA (4.5:1 body, 3:1 large text and controls) in both themes; the pairs and ratios are listed in the design doc and checked by a unit test
- [ ] Pages and components use tokens only: no literal Tailwind palette colours (`zinc`, `amber`, `sky`, `red`) remain in `src/routes` and `src/lib/components` except through the status tokens
- [ ] Chart series, sequential, diverging, and status palettes in `src/lib/viz` are re-derived from the brand hues and pass the dataviz validator in light and dark; status colours stay reserved and carry an icon or label
- [ ] Dark mode is a selected theme, not an automatic inversion: tokens are defined for both, the system preference is the default, and a toggle in the navigation persists per browser
- [ ] design/system/flaiover-dashboard.md gains a Theme section with the token table and contrast ratios; docs/users/flaiover.md shows both themes

## Tasks

## Notes
- Palette supplied by the operator on 2026-09-18. Roles above are a starting assignment; the first task confirms them against contrast before any component changes.
- Follow the dataviz skill's method for the chart palettes: assign by job, validate, then apply.
