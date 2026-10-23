---
id: D-0007
title: flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI
kind: cause
sources: [I-0011]
topics: [flai, cli, go, code, all]
paths: [flai/internal/issues, flai/internal/release]
date: 2026-09-19
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0007 flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI

I-0011 was closed on 2026-09-19.

The cause:

```text
flai stream log wrote duplicate same-second headings that markdownlint MD024 rejects in CI
```

The fix:

```text
S-0047: issue instances recorded in the same second share one heading, as narrative log entries have since the first fix; a test bumps twice at one instant.
```

Source: I-0011.
