---
id: D-0115
title: "Accept ADR-0067 as it stands, with S-0195's criteria settling its open consequence: 1 - if the remote has changed, our tags are invalid so we need to delete them and either rebase on the remote OR…"
kind: answer
sources: [TH-0070]
topics: [release, flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flai/internal/mcpserver, flai/cmd, flai/internal/harness, flai/internal/hostapi, flaiover/src, flai/internal/workitem/boardview.go, flai/internal/workitem/boardview_test.go, flai/internal/release/remote.go, flai/internal/release/remote_test.go, flai/internal/preview, design/system, docs, design/conventions/git.md]
date: 2026-10-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0115 Accept ADR-0067 as it stands, with S-0195's criteria settling its open consequence: 1 - if the remote has changed, our tags are invalid so we need to delete them and either rebase on the remote OR…

The operator, alex, answered TH-0070 at 2026-10-02T23:46:41Z.

What was asked:

```text
S-0195's first criterion needs you to accept ADR-0067 (design/adrs/0067-accepted-work-reaches-the-remote-only-when-it-is-published-and-agents-publish.md), or ask for changes, before anything is built.
```

The answer:

```text
1 - if the remote has changed, our tags are invalid so we need to delete them and either rebase on the remote OR merge and then we can re-tag after the verification completes.
```

Source: TH-0070.
