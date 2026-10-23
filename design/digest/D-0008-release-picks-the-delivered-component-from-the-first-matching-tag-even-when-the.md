---
id: D-0008
title: Release picks the delivered component from the first matching tag even when the story's commits never touched it
kind: cause
sources: [I-0016]
topics: [flai, cli, go, code, all]
paths: [flai/internal/issues, flai/internal/release]
date: 2026-09-19
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0008 Release picks the delivered component from the first matching tag even when the story's commits never touched it

I-0016 was closed on 2026-09-19.

The cause:

```text
Release picks the delivered component from the first matching tag even when the story's commits never touched it
```

The fix:

```text
S-0047: the release delivers only to a component the story's commits touched; a tag naming an untouched component never delivers. Versions already cut stand.
```

Source: I-0016.
