---
id: T-0406
type: task
nature: feature
title: "Users guide: a paragraph for every dashboard view"
status: done
parent: S-0019
owner: alex
created: 2026-09-24T08:28:56Z
updated: 2026-09-24T08:30:40Z
transitions:
  - to: ready
    at: 2026-09-24T08:29:00Z
    by: agent-S-0019
  - to: in-progress
    at: 2026-09-24T08:29:00Z
    by: agent-S-0019
  - to: done
    at: 2026-09-24T08:30:40Z
    by: agent-S-0019
stream: S-0019
tags: []
touches: [docs/users]
---
# T-0406 Users guide: a paragraph for every dashboard view

## Work
Rewrite `docs/users/flaiover.md` as the dashboard guide: how to start it and log in, the header (project switcher, host flai badge, theme), then one section per view in the order the navigation lists them (Overview, Board, Inbox, Activity, Charts, Docs, ADRs, Search, Host, Settings) and the pages reached from them (item page, review page, editor, new item, new ADR, login). Keep the detail already written; reorganise it under the views and remove the "full guide arrives with S-0019" line. Link the operators guide for running and securing it.

## Done when
- Every page under `flaiover/src/routes` has at least a paragraph in the guide.
- `status` is `active`, `updated` is today, and `scripts/lint-md.sh` passes.

## Notes
