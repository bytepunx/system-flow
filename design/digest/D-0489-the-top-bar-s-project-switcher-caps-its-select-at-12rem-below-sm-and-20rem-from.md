---
id: D-0489
title: The top bar's project switcher caps its select at 12rem below sm and 20rem from sm, with a shrinking, wrapping wrapper, so a long project label cannot make the page wider than a 390px viewport
kind: decision
sources: [I-0180, S-0557]
topics: [flaiover, dashboard, sveltekit]
paths: [flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts, flaiover/e2e/project-switcher.e2e.ts]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0489 The top bar's project switcher caps its select at 12rem below sm and 20rem from sm, with a shrinking, wrapping wrapper, so a long project label cannot make the page wider than a 390px viewport

I-0180 found that `ProjectSwitcher.svelte`'s `<select>` had no maximum width, so with two long-named projects connected the top bar, and the page, grew wider than a 390px viewport. S-0557 fixed it and recorded the widths in its narrative's Decisions (`wip/archive/agents/S-0557.md`):

```text
The select is capped at `max-w-48` (12rem) below `sm` and `max-w-xs` (20rem) from `sm`. 12rem leaves room at 390px beside the menu button, the name, and the theme button. The wrapper and label take `min-w-0 max-w-full`, and the wrapper takes `flex-wrap`, so the "n not connected: why?" link goes under the select. The current project's whole label is the first line of the select's `title`.
```

`flaiover/e2e/project-switcher.e2e.ts` reproduces I-0180 with two long-named projects at 390px and 1440px: it failed on the old switcher, the select ending at 548px, and passes with the cap. A story that adds to the top bar or changes the switcher keeps the cap, or measures again at 390px, rather than widening the select.

Sources: I-0180, S-0557.
