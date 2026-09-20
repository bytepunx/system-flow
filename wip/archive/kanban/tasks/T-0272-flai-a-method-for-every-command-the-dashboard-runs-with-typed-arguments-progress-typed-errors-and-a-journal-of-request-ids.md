---
id: T-0272
type: task
nature: feature
title: "flai: a method for every command the dashboard runs, with typed arguments, progress, typed errors, and a journal of request IDs"
status: done
parent: S-0075
owner: alex
created: 2026-09-20T08:49:14Z
updated: 2026-09-20T08:53:12Z
transitions:
  - to: ready
    at: 2026-09-20T08:49:15Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T08:49:15Z
    by: system-flow
  - to: done
    at: 2026-09-20T08:53:12Z
    by: system-flow
stream: S-0075
tags: []
---
# T-0272 flai: a method for every command the dashboard runs, with typed arguments, progress, typed errors, and a journal of request IDs

## Work
In internal/hostapi: one method per command the dashboard uses today, each validating its arguments as data and building the command line itself (the builders move from TypeScript with their tests), run as flai in the project's folder. The channel gains progress notifications for a request and data on an error. Exit codes that carry meaning arrive as typed errors (conflict, refused, a rule). A write with a request ID is answered from a journal when it repeats. A test enumerates the methods and fails when one has no validation test.

## Done when
- Tests for every builder, the refusals, progress, typed errors, and a repeat of each kind of write
- No method accepts a flag, a command line, or a path outside the manifest's folders

## Notes
