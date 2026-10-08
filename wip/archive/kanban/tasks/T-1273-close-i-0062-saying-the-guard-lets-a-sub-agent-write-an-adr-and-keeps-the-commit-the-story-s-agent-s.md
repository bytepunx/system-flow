---
id: T-1273
type: task
nature: remediation
title: Close I-0062 saying the guard lets a sub-agent write an ADR and keeps the commit the story's agent's
status: done
parent: S-0287
owner: alex
created: 2026-10-07T23:19:29Z
updated: 2026-10-08T07:22:53Z
transitions:
  - to: ready
    at: 2026-10-08T07:22:47Z
    by: agent-S-0287
  - to: in-progress
    at: 2026-10-08T07:22:48Z
    by: agent-S-0287
  - to: done
    at: 2026-10-08T07:22:53Z
    by: agent-S-0287
stream: S-0287
tags: [issues]
touches: [design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md, design/issues/summary.md]
after: [T-1270, T-1271, T-1272]
usage:
  source: log
  seconds: 5
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 3
      output: 993
      cache_read: 198953
      cache_write: 7771
      cost: 0.1133
---
# T-1273 Close I-0062 saying the guard lets a sub-agent write an ADR and keeps the commit the story's agent's

## Work

In the story's worktree, run `flai issue close I-0062 --reason "<what fixed it>"`, naming S-0287, the new ADR, and the guard rule: a sub-agent's `flai adr new`, `topics`, and `accept`, and `adr_new`, pass without a commit, and the commit stays the story's agent's. It writes the issue and regenerates `design/issues/summary.md`.

It waits for the other three tasks: the issue closes on the fix built, recorded, and documented.

## Done when

- I-0062 is closed with a reason that names S-0287, the ADR, and the rule.
- `design/issues/summary.md` no longer lists I-0062.

## Notes
