---
title: Icons
updated: 2026-10-21
status: active
topics: [dashboard]
---

# Icons

| | |
|-|-|
| Library | `@lucide/svelte` 1.61, ISC licence, one Svelte component per icon, imported by name so the build keeps only the icons used |
| Used in | `flaiover`: the nav rail's groups and pages, defined in `flaiover/src/lib/sitemenu.ts` (S-0458) |

## Why

- The nav rail needs one icon per group and per page, and flaiover had no icon library.
- Each icon is a Svelte 5 component drawing an inline SVG with `stroke="currentColor"`, so it takes its colour from the text utility of the element around it and its size from a `size-*` utility, both from the tokens.
- Chosen on TH-0553 question 2: the operator left the planner's recommendation standing when they finalized E-0027's stories.

## Rules

- Give every icon its size with a `size-*` utility class, never with the `size`, `color`, or `strokeWidth` props, so that no size or colour comes from outside the tokens.
- An icon beside a visible label is decorative: leave it `aria-hidden`, which the components set by default. An icon alone in a control needs the control's own accessible name.

## Considered

- Hand-written inline SVG components under `flaiover/src/lib/icons/`: no dependency, but every new page would need its icon drawn by hand.
- `unplugin-icons` with an Iconify set: more sets to choose from, at the cost of a Vite plugin and a build step for each icon.
