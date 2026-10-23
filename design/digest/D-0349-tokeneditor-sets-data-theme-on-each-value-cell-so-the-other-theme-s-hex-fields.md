---
id: D-0349
title: TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
kind: cause
sources: [I-0159]
topics: [flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src/lib/components/design/TokenEditor.svelte, flaiover/src/lib/components/design/TokenEditor.svelte.test.ts, flaiover/src/routes/design/design.svelte.test.ts, design/system/component-library.md, design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md, design/issues/summary.md, design/issues/I-0182-criteria-tick-is-refused-on-another-open-story-s-narrative-state-finding-which-the-story-s-own-scoped-check-leaves-out.md]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0349 TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read

I-0159 was closed on 2026-10-21.

The cause:

```text
TokenEditor sets data-theme on each value cell, so the other theme's hex fields are dark text on a dark field and cannot be read
```

The fix:

```text
S-0474: TokenEditor marks each value cell's theme with data-cell-theme, which no stylesheet reads, instead of data-theme, which layout.css reads as the dark palette and Tailwind's dark variant on any element. The page's theme now styles every cell, and a test in TokenEditor.svelte.test.ts fails while any element in the editor carries data-theme.
```

Source: I-0159.
