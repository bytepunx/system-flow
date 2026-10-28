---
id: D-0228
title: design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile
kind: cause
sources: [I-0089]
topics: [flai, cli, go, code, all]
paths: [flai/cmd/branch.go, flai/cmd/stream_sync.go, flai/cmd/stream_sync_test.go, flai/cmd/stream.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, design/adrs, design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0074-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md, design/issues/summary.md, design/issues/I-0092-two-story-branches-that-each-bump-the-same-issue-conflict-in-its-front-matter-whose-count-last-reported-and-updated-lines-both-rewrite.md, design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md, design/issues/I-0078-flai-check-finds-item-archive-outside-the-story-at-close-out.md]
date: 2026-10-07
status: absorbed
absorbed_by: "design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md#decision"
asked_again: []
cited: []
---

# D-0228 design/issues/summary.md is generated and committed, so a story branch that records an issue conflicts with any issue recorded on main meanwhile

I-0089 was closed on 2026-10-07.

The cause:

```text
`flai issue new`, `bump`, and `close` regenerate `design/issues/summary.md`, and the file is committed. A story records its issues on its branch. When an issue is also recorded on main while the story works, both sides have rewritten the same table, and the story's next `flai stream sync` stops on a conflict that the agent clears by regenerating the file. The same would stop an acceptance merge. It grows with the number of stories in progress at once.
```

The fix:

```text
Fixed by S-0278's regeneration of design/issues/summary.md when a rebase stops on it alone (ADR-0098), which flai stream sync and flai accept both run. S-0315 adds TestSyncAndAcceptRegenerateTheIssueSummaryChangedOnMainItself in flai/cmd/stream_sync_test.go, which records an issue on a story's branch while main records one and closes another itself, then syncs and accepts the story cleanly; it fails with the regeneration left out.
```

Source: I-0089.
