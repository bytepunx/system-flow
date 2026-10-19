---
id: T-1904
type: task
nature: feature
title: The harness frames a component with a 10px margin, measures all it draws, and keeps nothing of a prior component
status: done
parent: S-0449
owner: alex
created: 2026-10-19T05:04:04Z
updated: 2026-10-19T05:19:36Z
transitions:
  - to: ready
    at: 2026-10-19T05:05:09Z
    by: agent-S-0449
  - to: in-progress
    at: 2026-10-19T05:05:09Z
    by: agent-S-0449
  - to: done
    at: 2026-10-19T05:19:35Z
    by: agent-S-0449
stream: S-0449
tags: []
touches: [flai/internal/components/harness/harness.html, flai/internal/components/harness/harness.mjs, flai/internal/components/integration_test.go, design/system/component-library.md, docs/users/flai.md, flai/internal/components/testdata/harness/src/lib/Dialog.svelte, flai/internal/components/testdata/harness/src/lib/Leaky.svelte]
usage:
  source: log
  seconds: 866
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 81
      output: 33661
      cache_read: 4631874
      cache_write: 138938
      cost: 2.4762
---
# T-1904 The harness frames a component with a 10px margin, measures all it draws, and keeps nothing of a prior component

## Work

The harness in `flai/internal/components/harness/` frames the component, so the margin, the measurement, and what stays mounted change there (criteria 1, 2, and 4).

- `harness.html`: a 10px margin around the component, as padding on `#flai-component`, and nothing else between the component and the frame's edge.
- `harness.mjs` `measure()`: the size posted is everything the component draws, plus the margin. Count elements the component positions `fixed` or `absolute`, such as a dialog, so a dialog is framed whole (I-0160). A `fixed` element that only spans the viewport, such as a dialog's backdrop, adds nothing of its own, but its descendants count. Add the height of a horizontal scroll bar when the content is wider than the frame, so the bottom is not hidden by it. Measure again on a window resize and on a change below `#flai-component`, as well as on the ResizeObserver, so a new frame height settles.
- `harness.mjs` `close()`: once the component is unmounted, nothing it rendered stays. Empty `#flai-component`, and remove what the component added to `body` outside it, such as a portal.
- Framed (with `?origin=`), the harness renders only the component the page selects. It does not mount `?c=` itself, because the page sends its select on every ready, so a reload never shows the component the frame was first loaded with. Unframed it still renders `?c=`.
- A props message that arrives while a select is still loading its module is kept, and given to the component once it mounts.
- Tests: extend `flai/internal/components/integration_test.go`, whose browser checks run the harness in Playwright's Chromium. Each changed behaviour gets a check that fails without it: the margin, a fixed dialog framed whole, nothing of a prior component left after a select, no `?c=` mount when framed, and props kept across a load.
- Docs: `design/system/component-library.md` (the harness and its size message) and `docs/users/flai.md` (§ Render a component live), in the same change.

It waits for no task. It shares no path with the frame task, so the two run together.

## Done when

- The component sits 10px from the frame's top and left edges, and the posted size includes 10px on each side.
- A component that draws a `fixed inset-0` dialog posts a height that fits the whole dialog, plus the margin.
- After a select of B, `#flai-component` and `body` hold nothing that A rendered.
- A framed harness loaded with `?c=A` mounts nothing until the page selects.
- The integration checks pass with `flai test` on the paths changed, and fail without the change.

## Notes
