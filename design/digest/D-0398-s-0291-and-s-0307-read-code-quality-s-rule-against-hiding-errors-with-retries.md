---
id: D-0398
title: S-0291 and S-0307 read code-quality's rule against hiding errors with retries as allowing a bounded retry of a transient condition outside flai's control, keeping the last error and never retrying an HTTP error
kind: decision
sources: [S-0291, S-0307]
topics: [code, git, release, go, conventions]
paths: [install.sh, flai/internal/selfupgrade/selfupgrade.go, flai/internal/storygit/indexlock.go, design/conventions/code-quality.md]
date: 2026-10-28
status: compacted
absorbed_by: ""
asked_again: []
cited: []
---

# D-0398 S-0291 and S-0307 read code-quality's rule against hiding errors with retries as allowing a bounded retry of a transient condition outside flai's control, keeping the last error and never retrying an HTTP error
