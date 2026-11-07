---
id: T-2517
type: task
nature: remediation
title: I-0262 is closed with a reason naming the per-story schedule starts and the check
status: done
parent: S-0571
owner: alex
created: 2026-11-07T03:47:08Z
updated: 2026-11-07T04:07:16Z
transitions:
  - to: ready
    at: 2026-11-07T03:51:01Z
    by: agent-S-0571
  - to: in-progress
    at: 2026-11-07T04:06:56Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:07:16Z
    by: agent-S-0571
stream: S-0571
tags: [issues]
touches: [design/issues/I-0262-flai-serve-does-not-start-a-defined-role-that-takes-items-on-a-schedule-entry.md, design/issues/summary.md]
after: [T-2511, T-2512]
usage:
  source: log
  seconds: 20
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1075
      cache_read: 618526
      cache_write: 11497
      cost: 0.2343
---
# T-2517 I-0262 is closed with a reason naming the per-story schedule starts and the check

## Work

Run `flai issue close I-0262 --reason "<reason>"` in the story's worktree. The reason says that `flai serve`'s role scheduler now starts a defined role that takes stories once for each ready story under `each-ready-story`, with a held start per story. It also says that `flai check` reports the condition on a role that takes no stories. Name the test T-2511 added. flai rewrites `design/issues/summary.md`; commit both files.

It waits for T-2511 and T-2512 because the issue is closed only once the fix it names is in. It touches no path of the docs task's, so the two run together.

## Done when

- `flai issue show I-0262` reports it closed with that reason.
- `design/issues/summary.md` no longer lists it.
- `flai check --strict` is clean for the story.

## Notes
