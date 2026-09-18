---
id: T-0163
type: task
nature: improvement
title: "Dashboard: the preview lists uncommitted paths and the confirmation offers to include them"
status: done
parent: S-0051
owner: alex
created: 2026-09-18T21:00:23Z
updated: 2026-09-18T21:04:13Z
transitions:
  - to: ready
    at: 2026-09-18T21:01:33Z
    by: alex
  - to: in-progress
    at: 2026-09-18T21:01:33Z
    by: alex
  - to: done
    at: 2026-09-18T21:04:13Z
    by: alex
stream: S-0051
tags: []
touches: [flaiover]
---

# T-0163 Dashboard: the preview lists uncommitted paths and the confirmation offers to include them

## Work
`src/routes/api/items/[id]/acceptance/+server.ts`: drop `--yes` from the dry run so the preview checks what the move checks. `src/routes/api/items/[id]/move/+server.ts`: accept `include_uncommitted: true` in the body and add `--yes` only then, only when `to` is `done`. `AcceptConfirm.svelte`: when the preview has `uncommitted`, list the paths, say in one line that including them means the acceptance commit holds more than acceptance, and show a checkbox, off by default; accept is disabled until it is ticked; `onconfirm` receives the choice. The board page and the item page pass it to the move. Tests: the component with no paths, with paths and no choice, with the choice made; the move endpoint or its helper passing `--yes` only when asked.

## Done when
- The component and endpoint tests pass
- flaiover lint and svelte-check are clean

## Notes
