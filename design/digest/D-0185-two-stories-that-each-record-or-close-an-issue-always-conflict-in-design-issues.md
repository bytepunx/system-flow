---
id: D-0185
title: Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
kind: cause
sources: [I-0074]
topics: [flai, cli, go, code, all]
paths: [flai/cmd/branch.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, design/adrs, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/summary.md, design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md, design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-06
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0185 Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite

I-0074 was closed on 2026-10-06.

The cause:

```text
Two stories that each record or close an issue always conflict in design/issues/summary.md, whose updated line both rewrite
```

The fix:

```text
S-0278 (ADR-0098): flai stream sync and flai accept regenerate design/issues/summary.md and continue the rebase when it is the only conflict, and the trial merge between open story branches leaves it out, so no thread is opened over it. Two branches bumping the same issue is I-0092.
```

Source: I-0074.
