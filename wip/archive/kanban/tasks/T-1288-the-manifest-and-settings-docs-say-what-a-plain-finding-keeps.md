---
id: T-1288
type: task
nature: improvement
title: The manifest and settings docs say what a plain finding keeps
status: done
parent: S-0313
owner: alex
created: 2026-10-07T23:39:37Z
updated: 2026-10-12T15:56:34Z
transitions:
  - to: ready
    at: 2026-10-12T15:55:57Z
    by: agent-S-0313
  - to: in-progress
    at: 2026-10-12T15:55:57Z
    by: agent-S-0313
  - to: done
    at: 2026-10-12T15:56:34Z
    by: agent-S-0313
stream: S-0313
tags: [docs]
touches: [design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md]
after: [T-1286, T-1287]
usage:
  source: log
  seconds: 37
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 11
      output: 3293
      cache_read: 721340
      cache_write: 25613
      cost: 0.4006
---
# T-1288 The manifest and settings docs say what a plain finding keeps

## Work

`design/system/project-manifest.md` (the `format` bullet under `tests`) and the `tests[].format` row of `docs/operators/settings.md` say a `plain` tier "reads only the exit status". Since T-1286 a plain finding keeps the failure lines go test printed, then the last 20 lines and the exit status.

- Rewrite both to say what a plain finding keeps, and bump each file's `updated`.
- In `docs/users/flai.md`, under "Run the tests for what changed" and "Verify a story before review", say the same where the findings are described. Change nothing there if neither describes a plain finding.

It waits for T-1286, since it describes what that task built.

## Done when

- Neither document says a plain tier reads only the exit status.
- `flai test` on the three files passes the markdown lint.

## Notes

Drafted by the planner for S-0313.
