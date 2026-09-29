---
id: T-0586
type: task
nature: feature
title: Each file's summary line has a toggle control, and its panel shows signs in a margin, dim backgrounds, and an outline round each run
status: done
parent: S-0164
owner: alex
created: 2026-09-29T21:12:38Z
updated: 2026-09-29T21:19:43Z
transitions:
  - to: ready
    at: 2026-09-29T21:13:08Z
    by: agent-S-0164
  - to: in-progress
    at: 2026-09-29T21:16:02Z
    by: agent-S-0164
  - to: done
    at: 2026-09-29T21:19:43Z
    by: agent-S-0164
stream: S-0164
tags: []
touches: [flaiover/src/lib/components]
usage:
  source: log
  seconds: 221
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 25
      output: 10254
      cache_read: 1444004
      cache_write: 44947
      cost: 1.7729
---

# T-0586 Each file's summary line has a toggle control, and its panel shows signs in a margin, dim backgrounds, and an outline round each run

## Work

In `flaiover/src/lib/components/DiffView.svelte`: the summary line of each file gains a toggle control that says whether its panel is open, and the line names the panel it opens (`aria-expanded`, `aria-controls`). The panel draws the runs of T-0585: the sign of every added and removed line in a margin of its own, a dim green or red background behind each line, and a brighter green or red outline round each run. The colours are the theme's `good` and `danger`, which have a light and a dark value. A component test in `flaiover/src/lib/components/DiffView.svelte.test.ts`.

## Done when

- The toggle opens and closes a file's panel, and the panel is absent while closed.
- Every removed line has `-` in its margin and every added line `+`, and neither is in the line's text.
- Consecutive removed lines share one outlined block, as do consecutive added lines; lines apart are in blocks apart.
- The review page's tests still pass.

## Notes
