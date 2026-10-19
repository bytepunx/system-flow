---
id: T-1905
type: task
nature: feature
title: The Components page's frame takes only the selected component's answers and shows all the harness draws
status: done
parent: S-0449
owner: alex
created: 2026-10-19T05:04:20Z
updated: 2026-10-19T05:20:19Z
transitions:
  - to: ready
    at: 2026-10-19T05:05:10Z
    by: agent-S-0449
  - to: in-progress
    at: 2026-10-19T05:05:11Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T05:20:19Z
    by: agent-S-0449
stream: S-0449
tags: []
touches: [flaiover/src/lib/components/explorer/ComponentFrame.svelte, flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
usage:
  source: log
  seconds: 908
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 52
      output: 21665
      cache_read: 2981123
      cache_write: 89422
      cost: 1.5937
---
# T-1905 The Components page's frame takes only the selected component's answers and shows all the harness draws

## Work

The Components page's frame, `flaiover/src/lib/components/explorer/ComponentFrame.svelte`, shows what the harness draws, and only the selected component's answers (criteria 2 and 4).

- `rendered` and `error` answers about a component other than the one the frame now shows are ignored, as `event` already is. A late answer about the previous pick must not set or clear the error of the current one.
- The pane: the frame is as tall as the harness's size message says, which from T-1904 includes the 10px margin. Content wider than the frame scrolls inside the frame, so the wrapper must not clip it, and the frame's own scroll bars stay reachable. Keep the frame's width as the size option sets it.
- Tests: add `flaiover/src/lib/components/explorer/ComponentFrame.svelte.test.ts`, in the style of `flaiover/src/routes/design/components/components.svelte.test.ts` (fake harness messages from the frame's window). Each changed behaviour gets a test that fails without it: an error about a prior component is ignored, a rendered about a prior component does not clear the current error, and the frame takes the size message's height.
- Keep `design/components/flaiover/src/lib/components/explorer/ComponentFrame.examples.json` in step if the props change.
- Docs: `design/system/flaiover-dashboard.md` (§ The Components page) and `docs/users/flaiover.md` (§ Components), in the same change.

It waits for no task. It shares no path with T-1904, so the two run together. It relies only on the size message's shape, which T-1904 keeps.

## Done when

- A stale `rendered` or `error` from a prior pick changes nothing the frame shows.
- The frame's height is the harness's reported height, and wide content can be scrolled to.
- `flai test` on the changed paths passes, and the new tests fail without the change.

## Notes
