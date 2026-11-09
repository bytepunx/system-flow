---
id: T-2752
type: task
nature: remediation
title: Render repository markdown with raw HTML escaped, and test that script in a document does not run
status: in-progress
parent: S-0627
owner: alex
created: 2026-11-09T17:11:02Z
updated: 2026-11-09T17:41:34Z
transitions:
  - to: ready
    at: 2026-11-09T17:41:34Z
    by: agent-S-0627
  - to: in-progress
    at: 2026-11-09T17:41:34Z
    by: agent-S-0627
stream: S-0627
tags: [dashboard]
touches: [flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts]
---
# T-2752 Render repository markdown with raw HTML escaped, and test that script in a document does not run

## Work

In `flaiover/src/lib/markdown.ts`, build the renderer with `new MarkdownIt({ html: false, linkify: true, typographer: false })`. Raw HTML in a document, a thread entry, a message, or an editor preview then renders as text. Every `{@html}` site and `markdownInto` take their HTML from `render()`, so this one change covers all 13 sinks I-0294 and the explorer listed, `Threads.svelte` included.

No document under `design/` or `docs/` relies on raw HTML outside code spans and fences. The task-list checkboxes and the heading permalinks are tokens the plugins make, which markdown-it prints whatever `html` says, so they keep rendering. markdown-it's own `validateLink` already refuses `javascript:`, `vbscript:`, `file:`, and non-image `data:` links; keep `rewriteLink` from undoing that.

While in `render()`, install the `link_open` rule once at module level. Today each call wraps the previous rule again, so the rule chain grows by one closure per render.

Add cases to `flaiover/src/lib/markdown.test.ts`, which runs under jsdom:

- `<script>alert(1)</script>` and `<img src=x onerror=alert(1)>` in a document produce no `script` or `img` element and no `onerror` attribute; the text shows escaped. These fail on the current renderer.
- `[x](javascript:alert(1))` produces no `href` starting with `javascript:`.
- A `- [ ]` criterion still renders a checkbox, and a heading still gets its anchor and permalink.
- Rendering many times leaves link rewriting working once per link.

The first task of the story: it waits for none.

## Done when

- `render()` escapes raw HTML, and the script and `onerror` tests fail before the change and pass after it.
- `flai test flaiover/src/lib/markdown.ts flaiover/src/lib/markdown.test.ts` passes, including `vitest`, `flaiover-lint`, and `svelte-check`.
- The component tests that render markdown (`NewItemForm.svelte.test.ts`, `NewAdrForm.svelte.test.ts`, `Messages.svelte.test.ts`, `Review.svelte.test.ts`, `DocEditor.svelte.test.ts`) still pass.

## Notes
