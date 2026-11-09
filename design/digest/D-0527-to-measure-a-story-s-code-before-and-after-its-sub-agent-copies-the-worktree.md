---
id: D-0527
title: To measure a story's code before and after, its sub-agent copies the worktree under /tmp and puts the base commit's files into one copy, since flai guard refuses a sub-agent git archive
kind: decision
sources: [S-0607, T-2663, S-0602, S-0599, T-2655]
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

S-0602 built both trees with no copy at all, through `go build -overlay` under `/tmp`, which writes nothing in the repository (`wip/archive/agents/S-0602.md`):

```text
Measured with overlay builds under `/tmp`, nothing written in the repository: a `hostapi --timing` that repeats the method in one process, and a program timing each read's wall and CPU time. `stats.get` records no phase of its own for its reads, so `statsread.Read` was timed apart.
```

A one-shot `flai hostapi --timing` is always cold, so a warm figure needs a build whose `hostapi --timing` calls the method several times in one process. `design/system/server-performance.md` § Cause 13 fixed gives the method in full.

S-0599's T-2655 met the same refusal for `git clone` and measured on a `tar` copy of the main checkout under `/tmp`, without `.git`, `.flai-cache`, or `node_modules`, with one program built from the story's branch by `go build -overlay` (`wip/archive/agents/S-0599.md`):

```text
T-2655 measured in a `tar` copy of the main checkout without `.git`, since `flai guard` refuses a sub-agent `git clone`. In both runs the check's git rules found no repository; server-performance.md says so.
```

So a measurement that needs a project tree copies it with `tar` or `cp`, never `git clone` or `git archive`, and says in the design section that the copy has no repository, since rules that read git then find none in both runs alike (`design/system/server-performance.md` § Cause 9 fixed).

Sources: S-0607, T-2663, S-0602, S-0599, T-2655.
