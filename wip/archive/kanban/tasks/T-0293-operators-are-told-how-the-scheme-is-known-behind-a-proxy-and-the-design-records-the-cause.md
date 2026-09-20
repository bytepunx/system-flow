---
id: T-0293
type: task
nature: feature
title: Operators are told how the scheme is known behind a proxy, and the design records the cause
status: done
parent: S-0083
owner: alex
created: 2026-09-20T14:00:07Z
updated: 2026-09-20T14:07:31Z
transitions:
  - to: ready
    at: 2026-09-20T14:06:17Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T14:06:17Z
    by: system-flow
  - to: done
    at: 2026-09-20T14:07:31Z
    by: system-flow
stream: S-0083
tags: []
---
# T-0293 Operators are told how the scheme is known behind a proxy, and the design records the cause

## Work
docs/operators: what the dashboard needs from a TLS-terminating proxy or tunnel (X-Forwarded-Proto), and that plain HTTP by any address now works. design/system/flaiover-dashboard.md: the cause and the fix. An issue for the defect.

## Done when
- The markdown lint and flai check --strict pass, run after the last edit, in the worktree and in the main checkout

## Notes
