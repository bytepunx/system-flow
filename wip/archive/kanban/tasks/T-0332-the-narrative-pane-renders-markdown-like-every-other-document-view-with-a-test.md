---
id: T-0332
type: task
nature: feature
title: The narrative pane renders markdown like every other document view, with a test
status: done
parent: S-0088
owner: alex
created: 2026-09-22T23:54:59Z
updated: 2026-09-22T23:55:29Z
transitions:
  - to: ready
    at: 2026-09-22T23:55:08Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T23:55:10Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:55:29Z
    by: system-flow
stream: S-0088
tags: []
---

# T-0332 The narrative pane renders markdown like every other document view, with a test

## Work
`Review.svelte`'s narrative pane rendered `Current state` and `Next steps` with `whitespace-pre-wrap` on raw text, the only place in the dashboard that shows a document body without going through `$lib/markdown`'s `render`/`enhance` (the docs explorer, the item page, and the editors' previews all use it). Render both sections through `render()` (docPath `wip/agents/<id>.md`, for its relative-link rewriting) into a `.prose` block, and `enhance()` it after mount, the same shape `src/routes/items/[id]/+page.svelte` already uses.

## Done when
The pane shows real markdown (a numbered `Next steps` list as an actual `<ol>`, bold as `<strong>`, not literal `**text**`), a test proves it (not just incidentally passes), and the existing narrative test asserting literal `1. Accept.` text is updated for the new rendering. `pnpm run check`/`test:unit`/`lint` clean.

## Notes
