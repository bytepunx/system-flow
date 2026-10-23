---
id: D-0075
title: A story created without a component tag cannot be accepted when it touches two components
kind: cause
sources: [I-0024]
topics: [release, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/release, flai/cmd/release.go, flai/internal/hostapi, flai/internal/check, flaiover/src, flai/cmd/hostapi_reads_test.go, flai/cmd/push_test.go, flai/cmd/release_pending_test.go, flai/internal/pending/pending.go, flai/internal/preview/push.go, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users, design/issues]
date: 2026-10-01
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0075 A story created without a component tag cannot be accepted when it touches two components

I-0024 was closed on 2026-10-01.

The cause:

```text
A story created without a component tag cannot be accepted when it touches two components
```

The fix:

```text
S-0174: flai check warns on the untagged story (story.component-tag), and the pending plan names what it cannot plan
```

Source: I-0024.
