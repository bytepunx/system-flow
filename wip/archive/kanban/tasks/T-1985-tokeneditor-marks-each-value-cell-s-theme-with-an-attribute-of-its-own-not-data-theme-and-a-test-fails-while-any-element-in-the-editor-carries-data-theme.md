---
id: T-1985
type: task
nature: remediation
title: TokenEditor marks each value cell's theme with an attribute of its own, not data-theme, and a test fails while any element in the editor carries data-theme
status: done
parent: S-0474
owner: alex
created: 2026-10-21T12:06:42Z
updated: 2026-10-21T12:13:01Z
transitions:
  - to: ready
    at: 2026-10-21T12:11:11Z
    by: agent-S-0474
  - to: in-progress
    at: 2026-10-21T12:11:12Z
    by: agent-S-0474
  - to: done
    at: 2026-10-21T12:13:01Z
    by: agent-S-0474
stream: S-0474
tags: [dashboard]
touches: [flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/design.svelte.test.ts]
usage:
  source: log
  seconds: 109
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 2167
      cache_read: 839696
      cache_write: 29830
      cost: 0.45
---
# T-1985 TokenEditor marks each value cell's theme with an attribute of its own, not data-theme, and a test fails while any element in the editor carries data-theme

## Work

The cause of I-0159 is line 371 of `flaiover/src/lib/components/design/TokenEditor.svelte`: each editable cell is a `<div data-theme={c.label.replaceAll(' ', '-')}>`. `flaiover/src/routes/layout.css` scopes the dark palette (`--t-*`, `color-scheme: dark`) and Tailwind's `dark:` variant to `[data-theme='dark']` on any element, so the dark cell's inputs take the dark `bg-surface` while their text colour is inherited from the light page. The swatch needs no scope: it paints `c.token.css` directly. The attribute serves only as the tests' hook for finding a cell by theme.

- Replace `data-theme` on the cell with an attribute that no stylesheet reads, such as `data-cell-theme`, keeping the same value (`light`, `dark`, or `light-and-dark`).
- Change the two selectors that find a cell by `[data-theme="…"]` to the new attribute: `cell()` in `flaiover/src/lib/components/design/TokenEditor.svelte.test.ts` and `muted()` in `flaiover/src/routes/design/design.svelte.test.ts`.
- Add a test to `TokenEditor.svelte.test.ts` that renders the editor with both themes and fails while any element inside `[data-testid="token-editor"]` carries `data-theme`. It reproduces the cause, since jsdom applies no stylesheet and cannot show the colours themselves.
- Where a dev server can run, render the "Colours, two edited" example of `design/components/flaiover/src/lib/components/design/TokenEditor.examples.json` in the components harness on a light page and on a dark page, and check that every hex field is readable. Record the result, or why it could not be done, in the narrative.

This task waits for nothing: it is the first layer.

## Done when

- No element `TokenEditor.svelte` renders carries `data-theme`, and each value cell carries the new attribute.
- The new test fails against the old line 371 and passes with the change.
- `flai test` on the three paths passes: vitest, prettier, and eslint.

## Notes
