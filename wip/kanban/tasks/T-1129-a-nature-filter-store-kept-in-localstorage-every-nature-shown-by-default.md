---
id: T-1129
type: task
nature: improvement
title: A nature filter store kept in localStorage, every nature shown by default
status: backlog
parent: S-0302
owner: alex
created: 2026-10-06T23:20:00Z
updated: 2026-10-06T23:20:00Z
transitions: []
stream: S-0302
tags: [dashboard]
touches: [flaiover/src/lib/boardnatures.svelte.ts, flaiover/src/lib/boardnatures.svelte.test.ts]
---
# T-1129 A nature filter store kept in localStorage, every nature shown by default

## Work

Add `flaiover/src/lib/boardnatures.svelte.ts`, modelled on `boardtypes.svelte.ts` (S-0141): the natures feature, improvement, remediation, research, and experiment, a localStorage key `flaiover-board-natures`, defaults with every nature shown, `parseShownNatures(raw)` that keeps the default for anything missing, malformed, or not a boolean, a `BoardNatures` class with reactive `shown` state, `toggle(nature)` and `set(nature, on)` that persist the whole choice and survive a private window that refuses localStorage, and `isShown(nature)` that answers true for a nature the schema does not know, so such a card is never hidden. Export one instance, `boardNatures`.

Unit tests in `flaiover/src/lib/boardnatures.svelte.test.ts`: the defaults with nothing stored, a stored choice read back, malformed or partial storage falling back per nature, a toggle written to localStorage, and an unknown nature always shown. Waits for nothing: no path in common with the tints task, so the two run together.

## Done when

- [ ] `boardNatures` shows every nature by default and remembers a toggle in localStorage across a reload
- [ ] An unknown nature is always shown
- [ ] The new unit tests pass, and `scripts/flaiover-test.sh` passes

## Notes
