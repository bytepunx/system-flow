---
id: D-0117
title: "Sync refuses uncommitted changes, so a task is committed before it is synced: may git.md say so: a"
kind: answer
sources: [TH-0072]
topics: [flai, cli, go, template, conventions, code, all]
paths: [flai/internal/storygit, flai/cmd/stream.go, flai/cmd/branch.go, flai/cmd/stream_sync_test.go, flai/internal/harness, scripts/close-out.sh, template, design/conventions/git.md, design/conventions/work-management.md, design/conventions/delegation.md, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md, docs/users/flai-reference.md]
date: 2026-10-03
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0117 Sync refuses uncommitted changes, so a task is committed before it is synced: may git.md say so: a

The operator, alex, answered TH-0072 at 2026-10-03T00:52:11Z.

What was asked:

```text
S-0197's first criterion has `flai stream sync` refuse to run over uncommitted changes (today it uses `git rebase --autostash`). But the git convention, here and in the template, says to commit each task by running `flai stream sync` *first*, resolving, testing, and *then* committing the task's changes. Once sync refuses a dirty worktree, that order can't work: when a task finishes, its changes are uncommitted, so sync refuses.
```

The answer:

```text
a
```

Source: TH-0072.
