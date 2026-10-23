---
id: D-0147
title: flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number
kind: cause
sources: [I-0065]
topics: [cli, flai, go, code, all]
paths: [flai/internal/storygit, flai/internal/issues, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md, design/system/continuous-improvement.md, design/issues/I-0065-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md, design/issues/summary.md, design/issues/I-0057-flai-check-strict-stops-a-story-s-close-out-on-wip-findings-outside-the-story.md]
date: 2026-10-05
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0147 flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number

I-0065 was closed on 2026-10-05.

The cause:

```text
flai issue new numbers from the story's worktree only, so parallel story branches take the same issue number
```

The fix:

```text
S-0252: flai issue new numbers one past the highest issue on main, in every story worktree, and on every story branch (storygit.FolderNames), with a test reproducing the collision
```

Source: I-0065.
