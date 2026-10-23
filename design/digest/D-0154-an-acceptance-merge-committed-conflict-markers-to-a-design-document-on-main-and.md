---
id: D-0154
title: An acceptance merge committed conflict markers to a design document on main, and nothing caught it
kind: cause
sources: [I-0066]
topics: [flai, cli, go, code, all]
paths: [flai/internal/conflictmark, flai/internal/check, flai/cmd/branch.go, flai/cmd/accept_conflict_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0154 An acceptance merge committed conflict markers to a design document on main, and nothing caught it

I-0066 was closed on 2026-10-05.

The cause:

```text
An acceptance merge committed conflict markers to a design document on main, and nothing caught it
```

The fix:

```text
S-0253: flai check reports markdown.conflict-marker, an error, on each conflict marker line in the markdown under design, docs, wip, and the root, and flai accept refuses, merging nothing, a story branch any of whose added or changed files carries a marker, naming each path and line
```

Source: I-0066.
