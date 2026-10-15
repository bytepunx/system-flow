---
id: T-1541
type: task
nature: feature
title: The /design page lays out the project's tokens and guidance beside the live preview, says when there is no design system, and follows changes
status: done
parent: S-0382
owner: alex
created: 2026-10-13T17:47:14Z
updated: 2026-10-15T01:53:37Z
transitions:
  - to: ready
    at: 2026-10-15T01:46:36Z
    by: agent-S-0382
  - to: in-progress
    at: 2026-10-15T01:46:39Z
    by: agent-S-0382
  - to: done
    at: 2026-10-15T01:53:37Z
    by: agent-S-0382
stream: S-0382
tags: [dashboard]
touches: [flaiover/src/lib/components/design/TokenSheet.svelte, flaiover/src/lib/components/design/TokenSheet.svelte.test.ts, flaiover/src/routes/design/+page.svelte, flaiover/src/routes/design/design.svelte.test.ts]
after: [T-1538, T-1539, T-1540]
usage:
  source: log
  seconds: 418
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 89
      output: 33808
      cache_read: 5125599
      cache_write: 148139
      cost: 2.5679
---
# T-1541 The /design page lays out the project's tokens and guidance beside the live preview, says when there is no design system, and follows changes

## Work

`TokenSheet.svelte` lays out the tokens: colour swatches with their contrast in both themes, the type scale as sample text, and spacing, radii, and shadows drawn to size. `flaiover/src/routes/design/+page.svelte` loads `/api/design`, shows the sheet, the guidance rendered as markdown, and the `PreviewFrame`; with no document it names `flai design init`; and it reloads on a `document` change through `follow`, as the ADRs page does. Component and page tests mock `$lib/api` as `docs.svelte.test.ts` does. It waits for T-1538, T-1539, and T-1540: it reads the endpoint, sits in the menu, and shows the frame.

## Done when

- `/design` shows the sheet, the guidance, and the preview for a project with a document, and the command for one without, and its tests pass.

## Notes

- Layer 2.
