---
id: D-0208
title: The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
kind: cause
sources: [I-0100, S-0307]
topics: [cli, git, flai, go, code, all]
paths: [flai/internal/storygit/indexlock.go, flai/internal/storygit/indexlock_test.go, flai/internal/preview/accept.go, flai/cmd/accept.go, flai/cmd/accept_lock_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md, design/issues/summary.md, design/issues/I-0073-flai-check-finds-threads-archived-outside-the-story-at-close-out.md, design/issues/I-0076-flai-check-finds-wip-overlap-outside-the-story-at-close-out.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0208 The acceptance's commit step fails when another process holds git's index lock, and the acceptance cannot then be finished by flai
