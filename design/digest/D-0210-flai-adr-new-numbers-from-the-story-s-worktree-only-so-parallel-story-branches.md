---
id: D-0210
title: flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
kind: cause
sources: [I-0063]
topics: [cli, flai, go, code, all]
paths: [flai/internal/adr, flai/cmd/adr.go, flai/cmd/adr_test.go, docs/users/flai-reference.md, docs/users/flai.md, docs/users/flaiover.md, design/system/flai-cli.md, design/system/documentation-standard.md, design/issues/I-0063-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0210 flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number

I-0063 was closed on 2026-10-07.

The cause:

```text
flai adr new numbers from the story's worktree only, so parallel story branches take the same ADR number
```

The fix:

```text
S-0245: flai adr new numbers one past the highest ADR in the checkout, the main checkout, every story worktree (uncommitted included), the main branch, and every story branch, through storygit.FolderNames; TestAdrNewNumbersPastEveryStorysAdrs reproduces the collision. The issue-number instance was fixed by S-0252.
```

Source: I-0063.
