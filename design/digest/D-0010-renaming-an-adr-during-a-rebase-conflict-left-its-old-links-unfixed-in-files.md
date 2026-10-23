---
id: D-0010
title: Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit
kind: cause
sources: [I-0034]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/cmd, flaiover/src, design/system]
date: 2026-09-22
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0010 Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit

I-0034 was closed on 2026-09-22.

The cause:

```text
Renaming an ADR during a rebase conflict left its old links unfixed in files touched by a later commit
```

The fix:

```text
S-0081, commit 735702d: the three stale ADR-0032 references in dashboard-host-channel.md, flaiover-dashboard.md, and docs/operators/index.md are fixed
```

Source: I-0034.
