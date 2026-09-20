---
id: T-0252
type: task
nature: improvement
title: "flaiover: cancelling an epic or a story shows what goes with it before it happens"
status: done
parent: S-0070
owner: alex
created: 2026-09-20T06:16:44Z
updated: 2026-09-20T06:30:42Z
transitions:
  - to: ready
    at: 2026-09-20T06:30:42Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:30:42Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:30:42Z
    by: system-flow
stream: S-0070
tags: []
---
# T-0252 flaiover: cancelling an epic or a story shows what goes with it before it happens

## Work
The move endpoint takes `dry_run` and passes `--dry-run`; the board and the item page ask it first when the target is `cancelled`, and show a confirmation with the children listed by type and state, stories in review called out, and the reason field. Confirming posts the move; cancelling the dialog changes nothing.

## Done when
- Unit tests for the arguments and the dialog
- Checked in a browser against a scratch project: an epic with stories and tasks cancelled from the board and from the item page

## Notes
