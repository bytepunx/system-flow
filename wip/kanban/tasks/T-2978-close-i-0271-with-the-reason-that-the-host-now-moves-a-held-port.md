---
id: T-2978
type: task
nature: remediation
title: Close I-0271 with the reason that the host now moves a held port
status: done
parent: S-0697
owner: alex
created: 2026-11-10T10:29:13Z
updated: 2026-11-11T08:08:39Z
transitions:
  - to: ready
    at: 2026-11-11T08:08:10Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T08:08:10Z
    by: system-flow
  - to: done
    at: 2026-11-11T08:08:39Z
    by: agent-S-0697
stream: S-0697
tags: [issues]
touches: [design/issues/I-0271-flai-host-restarts-a-project-s-mcp-server-on-the-same-port-while-another.md, design/issues/summary.md]
after: [T-2974]
usage:
  source: log
  seconds: 29
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 8
      output: 1789
      cache_read: 837584
      cache_write: 25200
      cost: 0.4049
---
# T-2978 Close I-0271 with the reason that the host now moves a held port

## Work

Waits for T-2974, because the reason names the fix and its test, which T-2974 writes.

Run `flai issue close I-0271 --reason "<reason>"` in the story's worktree. The reason says that `hostLauncher.mcpAddr` and `mcpAddrFor` now look at a project's port at each start, and move the project to a free port when another process holds it. It names the test that reproduces the issue. `flai issue close` rewrites `design/issues/summary.md` with it.

## Done when

- I-0271 is closed, with a reason naming the fix and the test.
- `design/issues/summary.md` no longer lists I-0271 as open.
- `flai check --strict` reports nothing on either file.

## Notes
