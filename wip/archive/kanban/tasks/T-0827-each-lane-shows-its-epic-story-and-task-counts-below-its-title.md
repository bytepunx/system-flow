---
id: T-0827
type: task
nature: improvement
title: Each lane shows its epic, story, and task counts below its title
status: done
parent: S-0256
owner: alex
created: 2026-10-04T21:24:02Z
updated: 2026-10-04T21:27:06Z
transitions:
  - to: ready
    at: 2026-10-04T21:24:17Z
    by: agent-S-0256
  - to: in-progress
    at: 2026-10-04T21:24:18Z
    by: agent-S-0256
  - to: done
    at: 2026-10-04T21:27:06Z
    by: agent-S-0256
stream: S-0256
tags: []
touches: [flaiover/src/lib/lanes.ts, flaiover/src/lib/lanes.test.ts, flaiover/src/routes/board/+page.svelte, flaiover/src/routes/board/board.svelte.test.ts]
---
# T-0827 Each lane shows its epic, story, and task counts below its title

## Work

Count the epics, stories, and tasks each lane holds in a pure function in `flaiover/src/lib/lanes.ts`, with the short form `1 | 2 | 4` and the long form `1 epic, 2 stories, 4 tasks`, each type pluralized. Show the short form below each lane title on the board, with the long form as its tooltip and accessible name. Count every card the lane holds whichever types are ticked above the board; the done lane counts what it shows when release tags lag (`doneLane`). Waits for nothing.

## Done when

- [x] `lanes.test.ts` covers the counts and both forms, including singular, plural, and zero.
- [x] The board test shows each lane's count below its title with the long form as its title.
- [x] The flaiover tests for these files pass.

## Notes
