---
title: Tailwind CSS
updated: 2026-09-15
status: active
---

# Tailwind CSS

| | |
|-|-|
| Version | 4.3 |
| Used in | `flaiover` |
| Integration | `@tailwindcss/vite` plugin, no PostCSS config |
| Plugins | `@tailwindcss/typography` for rendered markdown |

## Why

Required by the brief. Version 4 uses CSS-first configuration, which keeps theme tokens in one `app.css` file and drops `tailwind.config.js`. Dark mode via `prefers-color-scheme` with a manual toggle stored in `localStorage`.

## Component approach

No component library initially. Small hand-built components with Tailwind classes keep the bundle small and the design distinct. If forms and dialogs grow, `bits-ui` (headless) is the first candidate; revisit with an ADR.
