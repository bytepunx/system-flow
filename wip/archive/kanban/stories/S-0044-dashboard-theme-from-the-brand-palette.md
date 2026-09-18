---
id: S-0044
type: story
nature: improvement
title: Dashboard theme from the brand palette
status: done
parent: E-0003
owner: alex
created: 2026-09-18T05:51:08Z
updated: 2026-09-18T16:15:12Z
transitions:
  - to: ready
    at: 2026-09-18T05:59:45Z
    by: alex
  - to: in-progress
    at: 2026-09-18T05:59:46Z
    by: alex
  - to: review
    at: 2026-09-18T06:10:05Z
    by: alex
  - to: in-progress
    at: 2026-09-18T06:14:23Z
    by: alex
  - to: review
    at: 2026-09-18T06:19:23Z
    by: alex
  - to: done
    at: 2026-09-18T16:15:12Z
    by: alex
tags: [dashboard]
touches: [flaiover/src/routes, flaiover/src/lib/components, flaiover/src/lib/viz, flaiover/src/app.css, docs/users/flaiover.md, design/system/flaiover-dashboard.md]
---

# S-0044 Dashboard theme from the brand palette

## Goal
flaiover wears the brand palette instead of Tailwind defaults: `#f7f3e3`, `#23b5d3`, `#119822`, `#645853`, `#054a91` (revised 2026-09-18 from `#0f1108`, `#23b5d3`, `#241909`, `#645853`, `#054a91`), applied through design tokens for a light and a dark theme, with charts and status colours derived from the same system and checked for contrast and colour-vision safety.

## Acceptance criteria
- [x] A token layer (CSS custom properties consumed by Tailwind) defines surface, ink, muted, primary, accent, border, and focus for light and dark: `#f7f3e3` as the light ground, `#645853` as muted text and borders, `#054a91` as primary (buttons, active navigation), `#23b5d3` as accent (links, focus rings, live indicators), `#119822` as the success family; the dark ground and surface are deep steps of the warm grey hue since the palette has no dark colour
- [x] Every text and interactive colour pair meets WCAG AA (4.5:1 body, 3:1 large text and controls) in both themes; the pairs and ratios are listed in the design doc and checked by a unit test
- [x] Pages and components use tokens only: no literal Tailwind palette colours (`zinc`, `amber`, `sky`, `red`) remain in `src/routes` and `src/lib/components` except through the status tokens
- [x] Chart series, sequential, diverging, and status palettes in `src/lib/viz` are re-derived from the brand hues and pass the dataviz validator in light and dark; status colours stay reserved and carry an icon or label
- [x] Dark mode is a selected theme, not an automatic inversion: tokens are defined for both, the system preference is the default, and a toggle in the navigation persists per browser
- [x] design/system/flaiover-dashboard.md gains a Theme section with the token table and contrast ratios; docs/users/flaiover.md shows both themes

## Tasks
- T-0128 Palette roles and contrast: derive light and dark token values from the five colours, verify WCAG AA pairs with a unit test, record the table
- T-0129 Token layer: CSS custom properties and Tailwind theme, explicit dark variant with a persisted toggle in the navigation, system preference as default
- T-0130 Apply tokens across routes and components; no literal Tailwind palette colours remain outside the status tokens
- T-0131 Chart palettes in src/lib/viz re-derived from the brand hues and validated with the dataviz validator in both themes
- T-0132 Docs: Theme section with tokens and ratios in the dashboard design; docs/users/flaiover.md describes both themes and the toggle
- T-0133 Re-derive both themes and the chart palettes for the revised palette (cream ground, green good family); update the design table, test, and story text

## Notes
- Palette supplied by the operator on 2026-09-18 and revised the same day (cream ground and a green replacing the near-black and brown). Roles were confirmed against contrast before component changes.
- Follow the dataviz skill's method for the chart palettes: assign by job, validate, then apply.
- 2026-09-18T06:14:23Z: moved to in-progress: palette revised by the operator: #f7f3e3, #23b5d3, #119822, #645853, #054a91
