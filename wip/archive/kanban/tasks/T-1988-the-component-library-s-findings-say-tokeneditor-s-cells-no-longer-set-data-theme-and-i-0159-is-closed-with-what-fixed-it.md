---
id: T-1988
type: task
nature: remediation
title: The component library's findings say TokenEditor's cells no longer set data-theme, and I-0159 is closed with what fixed it
status: done
parent: S-0474
owner: alex
created: 2026-10-21T12:07:01Z
updated: 2026-10-21T12:13:43Z
transitions:
  - to: ready
    at: 2026-10-21T12:13:25Z
    by: agent-S-0474
  - to: in-progress
    at: 2026-10-21T12:13:26Z
    by: agent-S-0474
  - to: done
    at: 2026-10-21T12:13:43Z
    by: agent-S-0474
stream: S-0474
tags: [dashboard]
touches: [design/system/component-library.md, design/issues/I-0159-tokeneditor-sets-data-theme-on-each-value-cell-so-the-other-theme-s-hex-fields-are-dark-text-on-a-dark-field-and-cannot-be-read.md, design/issues/summary.md]
after: [T-1985]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 9
      output: 1369
      cache_read: 530442
      cache_write: 18844
      cost: 0.2842
---
# T-1988 The component library's findings say TokenEditor's cells no longer set data-theme, and I-0159 is closed with what fixed it

## Work

- In `design/system/component-library.md` § What the pass found, rewrite the **TokenEditor** bullet as the other bullets record their fixes: say that since S-0474 each value cell marks its theme with the attribute T-1985 chose, which no stylesheet reads, so every hex field takes the page's field and text colours. Keep the link to I-0159 and bump the document's `updated`.
- Close the issue in the story's worktree with `flai issue close I-0159 --reason "<what fixed it>"`, naming the attribute and the test T-1985 added. flai writes the issue file and regenerates `design/issues/summary.md`.

It waits for T-1985 because the doc and the close reason name the attribute and the test that task chooses and adds.

## Done when

- The TokenEditor bullet in `design/system/component-library.md` describes the fix and links I-0159.
- I-0159 reads `status: closed` with the reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the three paths.

## Notes
