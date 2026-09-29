---
id: S-0167
type: story
nature: feature
title: Right click menu for the board
status: done
parent: E-0013
owner: alex
created: 2026-09-29T23:30:41Z
updated: 2026-09-29T23:51:37Z
transitions:
  - to: ready
    at: 2026-09-29T23:30:47Z
    by: alex
  - to: in-progress
    at: 2026-09-29T23:31:10Z
    by: agent-S-0167
  - to: review
    at: 2026-09-29T23:49:31Z
    by: agent-S-0167
  - to: done
    at: 2026-09-29T23:51:37Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, flai/internal/workitem, flai/internal/metrics, flai/internal/hostapi, flai/cmd, design/system/workflow.md, design/system/metrics.md, design/system/flaiover-dashboard.md, design/system/flai-cli.md, design/system/work-hierarchy.md, design/adrs, docs/users, design/conventions/work-management.md, design/issues]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1132
  models:
    - model: claude-opus-5-5
      input: 260
      output: 83989
      cache_read: 22401274
      cache_write: 260714
      cost: 8.2468
---
# S-0167 Right click menu for the board

## Goal

A right click menu for lanes on the board gives the operator convenient access to the following actions:

- create item
- move stories forward
- move stories backward
- change WIP limit

## Acceptance criteria
- [x] when taking the operator to the create screen, the lane the menu was created for is the starting lane for the card (only works for backlog, ready, and in-progress, all other lanes get changed to backlog as the default)
- [x] moving stories forward should only be an option for the backlog lane
- [x] moving stories backward should only be an option in ready, in progress, review, and cancelled (stories cannot be moved from done). cancelled stories move back to backlog
- [x] changing the wip limit should change the correct file(s) needed to make the change permanent

## Tasks
- T-0592 A story moves back a column: ready to backlog, in-progress to ready, cancelled to backlog
- T-0593 flai board limit sets a column's WIP limit, and the dashboard can run it through flai on the host
- T-0594 The create screen starts an item in the lane it was asked from
- T-0595 A right click on a board lane opens its menu: create, move stories forward or back, change the WIP limit

## Notes
- Verified by component tests of the board page and the create form (`flaiover/src/routes/board/lanemenu.svelte.test.ts`, `NewItemForm.svelte.test.ts`) and by writes through the tree's flai (`flaiover/src/lib/server/writes.test.ts`, `flai/cmd/board_limit_test.go`, `flai/internal/workitem/workitem_test.go`); not clicked through in a live browser.
- flai allowed no move back but review to in-progress, so the backward moves needed a rule change: ADR-0055. `completed` is now the last closing transition, so a reopened story is not counted as completed.
- The WIP limit lives only in `wip/kanban/board.md`; `flai board limit` and the `board.limit` host write set it. The convention's project addition stopped restating the numbers.
- Moving a story back to ready makes it a story entering ready: flai serve may start an agent for it.
