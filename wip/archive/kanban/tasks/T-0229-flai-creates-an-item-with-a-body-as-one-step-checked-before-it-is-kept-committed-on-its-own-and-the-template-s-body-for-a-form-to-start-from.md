---
id: T-0229
type: task
nature: feature
title: "flai creates an item with a body as one step: checked before it is kept, committed on its own, and the template's body for a form to start from"
status: done
parent: S-0059
owner: alex
created: 2026-09-19T09:21:36Z
updated: 2026-09-19T09:25:01Z
transitions:
  - to: ready
    at: 2026-09-19T09:21:37Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:21:37Z
    by: system-flow
  - to: done
    at: 2026-09-19T09:25:01Z
    by: system-flow
stream: S-0059
tags: []
---

# T-0229 flai creates an item with a body as one step: checked before it is kept, committed on its own, and the template's body for a form to start from

## Work
`flai epic new` and `flai story new` gain `--body-stdin`: the item's body below its heading comes from standard input instead of the template's empty sections, so the same is possible from a script. With it the creation is one step that happens or does not: `flai check` runs before and after, anything the new item introduces refuses it, the new file is removed and the parent restored, and the findings are returned (exit code 4 and JSON, as `flai doc save` does). `--autocommit` commits the new item and its parent on their own, path limited, unless the manifest sets `dashboard.autocommit: false`, with `--trailer` lines; nothing is pushed. `--print-body` prints the body the project's template gives that type, without creating anything, for a form to start from. Task creation is unchanged. Tests with real git: body, parent link, next ID, refusal leaving nothing behind, commit holding only the two paths, autocommit off, the template body.

## Done when
- The flags are tested end to end with real git
- `make flai-test` passes

## Notes
