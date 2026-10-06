---
id: T-1134
type: task
nature: improvement
title: The legend's type entries are bordered toggles that show or hide their cards
status: done
parent: S-0303
owner: alex
created: 2026-10-06T23:26:05Z
updated: 2026-10-06T23:50:44Z
transitions:
  - to: ready
    at: 2026-10-06T23:48:02Z
    by: agent-S-0303
  - to: in-progress
    at: 2026-10-06T23:48:03Z
    by: agent-S-0303
  - to: done
    at: 2026-10-06T23:50:44Z
    by: agent-S-0303
stream: S-0303
tags: [dashboard]
touches: [flaiover/src/lib/components/BoardLegend.svelte, flaiover/src/lib/components/BoardLegend.svelte.test.ts, flaiover/src/lib/cardcolour.ts]
usage:
  source: log
  seconds: 161
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 18
      output: 6287
      cache_read: 602263
      cache_write: 35128
      cost: 0.4721
---
# T-1134 The legend's type entries are bordered toggles that show or hide their cards

## Work

Make each type entry of `BoardLegend.svelte` a button (`aria-pressed`, keeping `data-type`) that toggles the type through the `boardTypes` store, `types.set(type, !types.shown[type])`, taking the store as an optional prop the way `BoardTypes.svelte` does so tests can pass their own. Give every type entry the border the nature entries have (`rounded border border-line px-1.5 py-0.5`), keeping its swatch. A type that is on gets a highlighted colouring, a tint of its own type colour from a new `typeHighlight` map in `$lib/cardcolour.ts` (such as `bg-type-epic/25` with a stronger border), written out in full for Tailwind; a type that is off keeps the plain colouring it has today. Give each button an accessible name that says what it does, such as `Show epics` with its pressed state, and a visible focus ring. The nature entries stay as they are.

Waits for nothing: the store's `shown` and `set` exist already, so it runs in the first layer beside the defaults task, with which it shares no path.

## Done when

- [ ] Each type entry is a bordered button, `aria-pressed` true when its type is shown and false when hidden.
- [ ] Clicking an entry flips its type in the store and in `localStorage`, and its colouring changes between the highlighted and the plain one.
- [ ] `BoardLegend.svelte.test.ts` covers the border, the toggle, the pressed state, and the colouring of both states, and passes with the lint and type check clean.

## Notes

Drafted by the planner for S-0303.
