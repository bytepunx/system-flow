---
id: T-1546
type: task
nature: feature
title: The Design page edits tokens with the preview following before saving, saves the token file through file.save, shows a refusal without losing the edits, and asks before leaving
status: done
parent: S-0383
owner: alex
created: 2026-10-13T17:47:41Z
updated: 2026-10-15T02:49:05Z
transitions:
  - to: ready
    at: 2026-10-15T02:39:22Z
    by: agent-S-0383
  - to: in-progress
    at: 2026-10-15T02:39:23Z
    by: agent-S-0383
  - to: done
    at: 2026-10-15T02:49:05Z
    by: agent-S-0383
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
after: [T-1545]
usage:
  source: log
  seconds: 582
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 75
      output: 51730
      cache_read: 4711694
      cache_write: 184542
      cost: 2.994
---
# T-1546 The Design page edits tokens with the preview following before saving, saves the token file through file.save, shows a refusal without losing the edits, and asks before leaving

## Work

`flaiover/src/routes/design/+page.svelte` gains an Edit mode:

- A `TokenEditor` stands beside each namespace.
- The page owns the edits. It feeds `editedThemes` to the `PreviewFrame`, so the preview follows each change before saving, and feeds `rolePairs` to the editors.
- The token file's raw text and hash come from `GET /api/design/tokens` (T-1594), so a light edit reaches dark as the resolver does.
- Save composes the new file with `composeTokenFile` and writes it with `PUT /api/design/tokens` `{path, content, hash}`.
- A refusal is shown with the edits kept:
  - a 422 with its findings, such as `design.tokens`
  - a 409 conflict, offering to load the current version or save over it, as `DocEditor.svelte` does
  - a 501 from a flai without `file.save`
- Discard goes back to the saved tokens.
- Leaving with unsaved changes asks first, with `beforeNavigate` and `onbeforeunload`, as the `/edit` route does.

Page tests cover a save, a refusal, a conflict, a 501, discard, and the leave guard, with `fetch` stubbed.

It waits for T-1545, whose editor it places. It runs beside T-1594: the two share no path, and they meet only at the endpoint's contract, which both tasks state.

## Done when

- `design.svelte.test.ts` passes for each case.

## Notes

- Layer 3, beside T-1594. Retargeted from `doc.save` at the start (TH-0454): the tokens are a DTCG file (ADR-0146), saved through `file.save` (ADR-0147).
