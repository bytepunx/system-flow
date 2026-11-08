---
id: D-0527
title: To measure a story's code before and after, its sub-agent copies the worktree under /tmp and puts the base commit's files into one copy, since flai guard refuses a sub-agent git archive
kind: decision
sources: [S-0607, T-2663]
topics: [server-side, code]
paths: [design/system/server-performance.md]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0527 To measure a story's code before and after, its sub-agent copies the worktree under /tmp and puts the base commit's files into one copy, since flai guard refuses a sub-agent git archive

S-0607's T-2663 measured `stats.get` on `main` against the story's branch for `design/system/server-performance.md` § Cause 17. From S-0607's `## Decisions` (`wip/archive/agents/S-0607.md`):

```text
T-2663 measured with copies of the worktree under `/tmp`, because the guard refused a sub-agent `git archive`. The base copy took `statsread.go`, `metrics.go`, and `reads.go` from 3087d9cf4, the only non-test files the story changes.
```

A story whose task measures base against head builds the base tree this way, from the non-test files the story changes, rather than having its sub-agent export the base commit with `git archive`.

Sources: S-0607, T-2663.
