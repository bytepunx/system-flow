---
id: T-1130
type: task
nature: improvement
title: The legend's nature tags are toggle buttons, lit when shown
status: backlog
parent: S-0302
owner: alex
created: 2026-10-06T23:20:08Z
updated: 2026-10-06T23:20:08Z
transitions: []
stream: S-0302
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardLegend.svelte, flaiover/src/lib/components/BoardLegend.svelte.test.ts]
after: [T-1128, T-1129]
---
# T-1130 The legend's nature tags are toggle buttons, lit when shown

## Work

In `flaiover/src/lib/components/BoardLegend.svelte`, make each nature tag a `button` with `aria-pressed` bound to `boardNatures.shown[nature]`: clicking it calls `boardNatures.toggle(nature)`. A shown nature wears its `natureShownTint` class, a hidden one its default `natureTint` class, so the default colouring means filtered out and the brighter one shown. The state is not told by colour alone: `aria-pressed` carries it, and a title such as "click to hide feature" or "click to show feature" says what a click does. Keep `data-nature` on each tag, and the type swatches as they are.

Update `BoardLegend.svelte.test.ts`: every tag starts pressed on its shown tint, a click turns it unpressed on its default tint and stores the choice, and a second click turns it back. Waits for the tints task, whose `natureShownTint` it reads, and the store task, whose `boardNatures` it binds.

## Done when

- [ ] Each nature tag toggles on click, shown on its brighter tint with `aria-pressed="true"`, hidden on its default tint with `aria-pressed="false"`
- [ ] The legend tests cover the default, a toggle off, and a toggle back on, and pass with `scripts/flaiover-test.sh`

## Notes
