---
id: T-1135
type: task
nature: improvement
title: The board drops its type checkboxes for the legend's toggles
status: backlog
parent: S-0303
owner: alex
created: 2026-10-06T23:26:13Z
updated: 2026-10-06T23:26:13Z
transitions: []
stream: S-0303
tags: [dashboard]
touches: [flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts, flaiover/src/lib/components/BoardTypes.svelte, flaiover/src/lib/components/BoardTypes.svelte.test.ts]
after: [T-1133, T-1134]
---
# T-1135 The board drops its type checkboxes for the legend's toggles

## Work

Remove `<BoardTypes />` and its import from `routes/board/+page.svelte`, and delete `BoardTypes.svelte` and `BoardTypes.svelte.test.ts`, moving any check of theirs that still matters, such as the choice lasting across mounts, to the store's or the legend's tests. Keep the legend where the board renders it now, or move it into the header row where the checkboxes stood if that reads better, and say which in the narrative. Add a test to `routes/board/board.svelte.test.ts` that clicking a type in the legend hides that type's cards and clicking it again shows them, and that the board opens with every type shown.

Waits for T-1133, which also changes `board.svelte.test.ts` and sets the default the new test asserts, and for T-1134, since the checkboxes can go only once the legend filters.

## Done when

- [ ] The board has no type checkboxes, and no file imports `BoardTypes.svelte`.
- [ ] Clicking a type in the legend hides and shows its cards on the board, covered by a test in `board.svelte.test.ts`.
- [ ] The flaiover tests, lint, type check, and build pass.

## Notes

Drafted by the planner for S-0303.
