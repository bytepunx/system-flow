---
id: D-0001
title: Path-filtered flai workflow did not run on the initial push to main
kind: cause
sources: [I-0008]
topics: []
paths: []
date: 2026-09-16
status: active
absorbed_by: ""
asked_again: []
cited: [S-0429, S-0432]
---

# D-0001 Path-filtered flai workflow did not run on the initial push to main

I-0008 was closed on 2026-09-16.

The cause:

```text
`.github/workflows/flai.yml` filters on `flai/**`. The first push of the repository to GitHub triggered `system-flow check` but not `flai`, although every file under `flai/` was new. It did run for Dependabot pull requests minutes later. Suspected: path filters need a previous commit to diff against on `push`.
```

The fix:

```text
the next push touching `flai/` triggered the workflow. First-push quirk of path filters; no change needed.
```

Source: I-0008.
