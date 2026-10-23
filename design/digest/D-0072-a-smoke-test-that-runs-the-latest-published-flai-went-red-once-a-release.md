---
id: D-0072
title: A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted
kind: cause
sources: [I-0038]
topics: [flai, cli, go, code, all]
paths: [flai/internal/serve, scripts/install-test.sh, design, docs]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0072 A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted

I-0038 was closed on 2026-10-01.

The cause:

```text
A smoke test that runs the latest published flai went red once a release changed what it pinned, after that release's story was accepted
```

The fix:

```text
Fixed by S-0112 (T-0391, b1fdfe9): install-test checks self-upgrade on a copy outside any project. Its macOS path mismatch is I-0045 (S-0180).
```

Source: I-0038.
