---
id: T-1519
type: task
nature: feature
title: Repository file links in rendered markdown, threads, reports, and plain refusal text open the editor at the line they name
status: done
parent: S-0378
owner: alex
created: 2026-10-13T16:57:18Z
updated: 2026-10-13T18:19:39Z
transitions:
  - to: ready
    at: 2026-10-13T18:11:22Z
    by: agent-S-0378
  - to: in-progress
    at: 2026-10-13T18:11:22Z
    by: agent-S-0378
  - to: done
    at: 2026-10-13T18:19:39Z
    by: agent-S-0378
stream: S-0378
tags: [dashboard]
touches: [flaiover/src/lib/filelinks.ts, flaiover/src/lib/filelinks.test.ts, flaiover/src/lib/markdown.ts, flaiover/src/lib/markdown.test.ts, flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/AnalyzerPanel.svelte, flaiover/src/lib/editor.svelte.ts, flaiover/src/lib/components/EditorDialog.svelte, flaiover/src/lib/components/EditorDialog.svelte.test.ts]
after: [T-1518]
usage:
  source: log
  seconds: 497
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 86
      output: 405
      cache_read: 4387501
      cache_write: 110461
      cost: 2.0588
---
# T-1519 Repository file links in rendered markdown, threads, reports, and plain refusal text open the editor at the line they name

## Work

- Add `flaiover/src/lib/filelinks.ts`. It finds repository markdown paths in plain text, `path`, `path:line`, and `path:line:column`, under the manifest's design, docs, and wip folders. It turns them into links that call `openEditor` (T-1518) with the path, the line, and a story when the caller gives one. Refusal text is plain, shown in a `<pre>` today.
- In `flaiover/src/lib/markdown.ts`, `enhance` gives each relative link to a markdown file, and each code span that holds such a path, an edit control beside it that opens the editor at the line it names. Following the link still shows the rendered document, as today.
- Use them in `Threads.svelte` (its entries and anchors) and in `AnalyzerPanel.svelte` (the report links), so threads and analysis reports offer the editor. Pages that render through `render` and `enhance`, a document, an item, or a narrative, get it from `markdown.ts`.
- Leave out paths the editor cannot edit: generated files, threads, issue files, and the archive (ADR-0023 mode `none`).
- Waits for T-1518, whose `openEditor` it calls.

## Done when

- `filelinks.test.ts` covers `path`, `path:line`, `path:line:column`, a path outside the folders, and a path the editor cannot edit.
- `markdown.test.ts` covers the edit control on a relative link and on a code-span path with a line.
- `scripts/flai.sh test` on these paths passes.

## Notes

- Criterion 1 of S-0378. ADR-0023's modes decide what is editable. `documentation.md`'s rule against linking to line numbers is about documents, not the dashboard's own links.
