---
id: T-0156
type: task
nature: remediation
title: Acceptance that cannot run git for a worktree says what to do, not the raw git error
status: done
parent: S-0050
owner: alex
created: 2026-09-18T19:47:24Z
updated: 2026-09-18T19:56:47Z
transitions:
  - to: ready
    at: 2026-09-18T19:55:20Z
    by: alex
  - to: in-progress
    at: 2026-09-18T19:55:20Z
    by: alex
  - to: done
    at: 2026-09-18T19:56:47Z
    by: alex
stream: S-0050
tags: []
touches: [flai/cmd, flaiover]
---

# T-0156 Acceptance that cannot run git for a worktree says what to do, not the raw git error

## Work
In `flai/cmd/accept.go`, when a git command for the story worktree fails because git cannot open it, return an error that says the worktree cannot be read from here, and to accept from a host shell with `flai accept <id>`, in place of the bare git output; keep the git text as the wrapped cause for logs. The dashboard shows flai's message, so check `AcceptConfirm.svelte` and the move endpoint render it without truncation. Test with a worktree whose `.git` file points at a missing path.

## Done when
- The test reproduces I-0017's failure and asserts the new message
- The dashboard component test, if one pins error rendering, passes

## Notes
