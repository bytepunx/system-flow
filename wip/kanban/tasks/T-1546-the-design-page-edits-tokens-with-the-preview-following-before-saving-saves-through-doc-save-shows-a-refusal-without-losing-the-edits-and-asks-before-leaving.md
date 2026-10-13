---
id: T-1546
type: task
nature: feature
title: The Design page edits tokens with the preview following before saving, saves through doc.save, shows a refusal without losing the edits, and asks before leaving
status: backlog
parent: S-0383
owner: alex
created: 2026-10-13T17:47:41Z
updated: 2026-10-13T17:47:41Z
transitions: []
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
after: [T-1545]
---
# T-1546 The Design page edits tokens with the preview following before saving, saves through doc.save, shows a refusal without losing the edits, and asks before leaving

## Work

On `flaiover/src/routes/design/+page.svelte`, an Edit mode puts a `TokenEditor` beside each namespace and feeds the edited tokens to the `PreviewFrame`, so it follows each change before saving. Save reads the document's hash with `GET /api/docs/edit`, composes the raw text with the edit model, and writes it with `PUT /api/docs/file`, as `DocEditor.svelte` does; a refusal, a `design.tokens` finding or a conflict, is shown with the edits kept. Discard goes back to the saved tokens. Leaving with unsaved changes asks first, with `beforeNavigate` and `onbeforeunload` as the `/edit` route does. Page tests cover save, a refusal, a conflict, discard, and the leave guard. It waits for T-1545, whose editor it places.

## Done when

- `design.svelte.test.ts` passes for each case.

## Notes

- Layer 3. S-0378, in progress, changes the document endpoints this task calls; read them as they are when the story starts.
