---
id: T-1926
type: task
nature: improvement
title: I-0157 is closed with a reason naming the task-done rules number that fixed it
status: done
parent: S-0445
owner: alex
created: 2026-10-20T20:42:30Z
updated: 2026-10-20T21:06:12Z
transitions:
  - to: ready
    at: 2026-10-20T21:06:02Z
    by: agent-S-0445
  - to: in-progress
    at: 2026-10-20T21:06:03Z
    by: agent-S-0445
  - to: done
    at: 2026-10-20T21:06:12Z
    by: agent-S-0445
stream: S-0445
tags: [issues]
touches: [design/issues/I-0157-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md, design/issues/summary.md]
after: [T-1924, T-1925]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 753
      cache_read: 208996
      cache_write: 5488
      cost: 0.097
---
# T-1926 I-0157 is closed with a reason naming the task-done rules number that fixed it

## Work

Run `flai issue close I-0157 --reason "<text>"` in the story's worktree. The reason names the ADR of T-1920, the rules number, the tag annotation, and the tests of T-1923 and T-1924 that reproduce the instances.

It waits for T-1924 and T-1925 so that the fix is whole and documented before the issue says so.

## Done when

- I-0157's status is closed, with the reason under its Remediation.
- `design/issues/summary.md` no longer lists it.

## Notes

The command writes under the checkout it runs in, so it runs in the worktree, not the main checkout.
