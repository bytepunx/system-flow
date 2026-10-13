---
id: T-1518
type: task
nature: feature
title: The document editor opens over any page at a line, on a story's branch when given one, and shows the lint findings as the operator types
status: backlog
parent: S-0378
owner: alex
created: 2026-10-13T16:57:08Z
updated: 2026-10-13T16:57:08Z
transitions: []
stream: S-0378
tags: [dashboard]
touches: [flaiover/src/lib/components/DocEditor.svelte, flaiover/src/lib/components/DocEditor.svelte.test.ts, flaiover/src/lib/components/EditorDialog.svelte, flaiover/src/lib/components/EditorDialog.svelte.test.ts, flaiover/src/lib/editor.svelte.ts, flaiover/src/lib/edit.ts, flaiover/src/lib/edit.test.ts, flaiover/src/routes/+layout.svelte, "flaiover/src/routes/edit/[...path]/+page.svelte"]
after: [T-1515]
---
# T-1518 The document editor opens over any page at a line, on a story's branch when given one, and shows the lint findings as the operator types

## Work

- Give `DocEditor.svelte` (S-0040, a `<textarea>`) a `line` prop that puts the caret at the start of that line and scrolls it into view, and a `story` prop it passes to `/api/docs/edit` and to the save. Put the line-to-offset arithmetic in `flaiover/src/lib/edit.ts`.
- Show the lint findings below the text: post the content to `/api/docs/lint` after a pause in typing, list each finding's line, rule, and message, and move the caret to its line on a click. Say that these are flai's subset of the rules and that the markdown tier decides at acceptance. Show nothing when the endpoint answers 501.
- Add `EditorDialog.svelte`, which holds `DocEditor` in a dialog over the current page. Open it with `openEditor({path, line, story, onsaved})` from a store in `flaiover/src/lib/editor.svelte.ts`. Mount the dialog once in `flaiover/src/routes/+layout.svelte`, so a link on any page opens it without leaving the page. Closing with unsaved changes asks first, as the `/edit` route does.
- Make `/edit/<path>` read `?line=` and `?story=`, for links opened in a new tab.
- Waits for T-1515, whose findings it lists. It uses T-1517's endpoints by the shapes T-1517 names, and runs beside it.

## Done when

- `DocEditor.svelte.test.ts` covers the caret at a given line, findings listed from a mocked `/api/docs/lint`, a click on a finding moving the caret, and a save with `story`.
- `EditorDialog.svelte.test.ts` covers opening from the store, the confirmation on unsaved changes, and `onsaved` called after a save.
- `edit.test.ts` covers the line-to-offset arithmetic, and `scripts/flai.sh test` on these paths passes.

## Notes

- Criteria 1 and 2 of S-0378. The editor opens over the page rather than on a page of its own, so that after a save the operator is still where the acceptance was refused (criterion 4).
