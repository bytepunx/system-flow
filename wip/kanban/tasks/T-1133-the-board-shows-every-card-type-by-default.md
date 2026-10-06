---
id: T-1133
type: task
nature: improvement
title: The board shows every card type by default
status: backlog
parent: S-0303
owner: alex
created: 2026-10-06T23:25:59Z
updated: 2026-10-06T23:25:59Z
transitions: []
stream: S-0303
tags: [dashboard]
touches: [flaiover/src/lib/boardtypes.svelte.ts, flaiover/src/lib/boardtypes.svelte.test.ts, flaiover/src/routes/board/board.svelte.test.ts]
---
# T-1133 The board shows every card type by default

## Work

Change the defaults in `$lib/boardtypes.svelte.ts` from stories alone to every type shown (`{ epic: true, story: true, task: true }`), for a browser with no stored choice and for any type the stored choice leaves out or holds a non-boolean for. Keep the `localStorage` key `flaiover-board-types`, so a browser that already chose keeps its choice. Rewrite the file's opening comment, which still names checkboxes and stories alone, to say the type legend toggles the choice. Update `boardtypes.svelte.test.ts` for the new defaults, and fix any test in `routes/board/board.svelte.test.ts` that relied on stories alone being the default by setting the types it needs with `boardTypes.set`.

Waits for nothing: it is the first layer, beside the legend task, with which it shares no path.

## Done when

- [ ] `parseShown(null)`, unparsable text, and an object missing a type give that type shown, and a unit test says so.
- [ ] A stored choice is still read back and kept as it was.
- [ ] `routes/board/board.svelte.test.ts` and `boardtypes.svelte.test.ts` pass, with the lint and type check clean.

## Notes

Drafted by the planner for S-0303.
