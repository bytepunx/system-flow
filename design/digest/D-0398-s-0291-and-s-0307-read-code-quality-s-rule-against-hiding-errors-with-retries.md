---
id: D-0398
title: S-0291 and S-0307 read code-quality's rule against hiding errors with retries as allowing a bounded retry of a transient condition outside flai's control, keeping the last error and never retrying an HTTP error
kind: decision
sources: [S-0291, S-0307]
topics: [code, git, release, go, conventions]
paths: [install.sh, flai/internal/selfupgrade/selfupgrade.go, flai/internal/storygit/indexlock.go, design/conventions/code-quality.md]
date: 2026-10-28
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0398 S-0291 and S-0307 read code-quality's rule against hiding errors with retries as allowing a bounded retry of a transient condition outside flai's control, keeping the last error and never retrying an HTTP error

`design/conventions/code-quality.md` says "Do not hide errors with retries". Two stories retried anyway, each a bounded number of times, for a condition that passes on its own:

- S-0291 gave `install.sh` and `flai self-upgrade` three attempts at a GitHub call a dropped connection cut short (I-0086), and "never retries an HTTP error". Its plan was approved on TH-0367 (D-0251).
- S-0307 gave the acceptance commit `storygit.RunPastIndexLock`, about nine seconds of retries while git reports `index.lock': File exists`, because "a lock held for a moment is a passing condition rather than a defect" (I-0100, D-0208).

Both keep the last error and report it when the attempts run out. No operator ruling settles this reading; it is the stories' own.

Sources: S-0291, S-0307.
