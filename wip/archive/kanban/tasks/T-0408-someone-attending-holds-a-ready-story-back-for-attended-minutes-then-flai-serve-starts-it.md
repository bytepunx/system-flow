---
id: T-0408
type: task
nature: remediation
title: Someone attending holds a ready story back for attended minutes, then flai serve starts it
status: done
parent: S-0114
owner: alex
created: 2026-09-24T08:34:03Z
updated: 2026-09-24T08:37:57Z
transitions:
  - to: ready
    at: 2026-09-24T08:34:44Z
    by: agent-S-0114
  - to: in-progress
    at: 2026-09-24T08:34:44Z
    by: agent-S-0114
  - to: done
    at: 2026-09-24T08:37:57Z
    by: agent-S-0114
stream: S-0114
tags: []
touches: [flai/internal/serve]
---
# T-0408 Someone attending holds a ready story back for attended minutes, then flai serve starts it

## Work

In `flai/internal/serve/agents.go`, the launcher's look already runs every minute, on every change under `kanban/` and `threads/`, and when an agent it started ends. What has no bound is the attended rule: while a sign of an agent stays fresh, a ready story is never started, however long it waits. Give the hold a timer: a story that only someone attending holds back is started once it has been held for the attended window (`attended_minutes`, 6 by default), the same window that says how recent a sign must be. The `waiting` reason says until when the attending agent has to pull it. Report the in-progress limit ahead of attendance when the limit is what is full, so the reason is the true one. Log `agent not started` when a story's reason changes in kind, not when only the sign's age in it changes, as the design already says. Tests in `agents_test.go` with a clock the test moves.

## Done when

- [x] A ready story someone attending holds back is started once the hold has lasted the attended window, and the reason before that names the time
- [x] A full in-progress limit is reported as the reason even when someone is attending
- [x] `agent not started` is logged once per reason, not once per look
- [x] `make flai-test` passes

## Notes
