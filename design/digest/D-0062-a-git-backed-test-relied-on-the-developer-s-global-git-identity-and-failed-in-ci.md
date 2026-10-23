---
id: D-0062
title: A git-backed test relied on the developer's global git identity and failed in CI
kind: cause
sources: [I-0013]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/workitem, flai/cmd, flai/internal/check, flaiover/src/routes/docs, docs/users/flai.md]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0062 A git-backed test relied on the developer's global git identity and failed in CI

I-0013 was closed on 2026-10-01.

The cause:

```text
A git-backed test relied on the developer's global git identity and failed in CI
```

The fix:

```text
Fixed by S-0037 (411ca5a): git-backed tests set their own identity. The inverse, the host's identity leaking into a test, is I-0045 (S-0180).
```

Source: I-0013.
