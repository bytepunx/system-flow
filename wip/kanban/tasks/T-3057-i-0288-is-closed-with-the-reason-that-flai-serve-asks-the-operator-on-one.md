---
id: T-3057
type: task
nature: improvement
title: I-0288 is closed with the reason that flai serve asks the operator on one thread when nothing can be made ready
status: backlog
parent: S-0691
owner: alex
created: 2026-11-10T10:35:20Z
updated: 2026-11-10T10:35:20Z
transitions: []
stream: S-0691
tags: [issues]
touches: [design/issues/I-0288-the-board-stood-empty-for-three-days-while-the-backlog-held-only-drafts-waiting.md, design/issues/summary.md]
after: [T-3050]
---
# T-3057 I-0288 is closed with the reason that flai serve asks the operator on one thread when nothing can be made ready

## Work

Run `flai issue close I-0288 --reason "<reason>"`. The reason names T-3043's ADR, and says that flai serve now opens one thread to the operator on `wip/kanban/board.md` when nothing can be made ready. That thread names each draft and the operator's thread it waits on, and each `after` chain's root, and it is reminded like any other thread. The command rewrites `design/issues/summary.md`.

It waits on T-3050, the change that removes the cause. It shares no path with the documentation task, so the two run together.

## Done when

- I-0288's status is `closed`, with that reason, and `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` passes.

## Notes
