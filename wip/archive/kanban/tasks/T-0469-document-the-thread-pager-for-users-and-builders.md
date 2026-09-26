---
id: T-0469
type: task
nature: feature
title: Document the thread pager for users and builders
status: done
parent: S-0133
owner: alex
created: 2026-09-26T17:49:20Z
updated: 2026-09-26T17:52:30Z
transitions:
  - to: ready
    at: 2026-09-26T17:49:31Z
    by: agent-S-0133
  - to: in-progress
    at: 2026-09-26T17:51:56Z
    by: agent-S-0133
  - to: done
    at: 2026-09-26T17:52:30Z
    by: agent-S-0133
stream: S-0133
tags: []
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
---

# T-0469 Document the thread pager for users and builders

## Work

Say in `docs/users/flaiover.md` that threads on an item or document are shown one at a time with the indicator and arrows, and in `design/system/flaiover-dashboard.md` under the workbench's threads line.

## Done when

- Both documents describe the pager, with `updated` bumped, and `scripts/lint-md.sh` passes.

## Notes
