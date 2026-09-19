---
id: T-0175
type: task
nature: feature
title: "Editor view: body and front matter, preview, save, findings, and conflict handling"
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:30Z
updated: 2026-09-19T02:23:19Z
transitions:
  - to: ready
    at: 2026-09-19T02:19:34Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:19:34Z
    by: alex
  - to: done
    at: 2026-09-19T02:23:19Z
    by: alex
stream: S-0040
tags: []
touches: [flaiover/src/routes, flaiover/src/lib/components]
---

# T-0175 Editor view: body and front matter, preview, save, findings, and conflict handling

## Work
Add an editor route, `/edit/[...path]`, reached from an Edit link on the document page and the item page, shown only when the board is writable. Load through `/api/docs/edit`. Mode `body`: the front matter is displayed read-only with a line saying flai owns it, and the body is a monospace textarea. Mode `full`: front matter YAML and body are both editable. Mode `none`: no editor, the reason shown. A preview beside or under the textarea renders the body with the same `render` and `enhance` the explorer uses. Save sends the hash; on success show the commit (or that it was left uncommitted) and adopt the new hash; on 422 list the findings and keep the text; on 409 show the diff and offer to load the current version or to save over it, which resends with the current hash. Warn before leaving with unsaved changes. No new dependency: a textarea, not an editor library. Component or page-level tests for the three modes, a refused save, and a conflict.

## Done when
- The tests pass; lint and svelte-check are clean
- A production build succeeds

## Notes
