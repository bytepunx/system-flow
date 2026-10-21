---
id: T-1978
type: task
nature: remediation
title: Close I-0173 saying the stubborn child now signals that it ignores SIGTERM
status: done
parent: S-0463
owner: alex
created: 2026-10-21T07:32:12Z
updated: 2026-10-21T07:36:51Z
transitions:
  - to: ready
    at: 2026-10-21T07:36:40Z
    by: agent-S-0463
  - to: in-progress
    at: 2026-10-21T07:36:42Z
    by: agent-S-0463
  - to: done
    at: 2026-10-21T07:36:51Z
    by: agent-S-0463
stream: S-0463
tags: [issues]
touches: [design/issues/I-0173-testachildisstoppedonlyonceitsprocesshasended-fails-under-load-when-the-stubborn-child-dies-on-sigterm-before-it-ignores-it.md, design/issues/summary.md]
after: [T-1977]
usage:
  source: log
  seconds: 9
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 765
      cache_read: 288474
      cache_write: 4280
      cost: 0.1072
---
# T-1978 Close I-0173 saying the stubborn child now signals that it ignores SIGTERM

## Work

- In the story's worktree, run `flai issue close I-0173 --reason "<reason>"`. The reason names the fix: the `stubborn` test child writes a file once it ignores SIGTERM, and the two host tests wait for it instead of sleeping 100 ms.
- It rewrites the issue file and `design/issues/summary.md`.
- Waits for the first task: the issue is closed only once its fix is in and its tests pass.

## Done when

- I-0173 has status `closed` with that reason.
- `design/issues/summary.md` no longer lists it as open.
- `flai check --strict` is clean.

## Notes

Drafted by the planner.
