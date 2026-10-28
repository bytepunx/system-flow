---
id: D-0267
title: Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
kind: cause
sources: [I-0092]
topics: [cli, git, continuous-improvement, flai, go, code, all]
paths: [design/adrs/README.md, flai/internal/issues/fold.go, flai/internal/issues/fold_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/internal/taskdone/taskdone.go, flai/cmd/stream_sync_test.go, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/summary.md, design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md, flai/internal/issues/merge.go, flai/internal/issues/merge_test.go, flai/internal/issues/generated.go, flai/cmd/stream_sync.go, docs/users/flai-reference.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md#decision"
asked_again: []
cited: []
---

# D-0267 Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite

I-0092 was closed on 2026-10-08.

The cause:

```text
Two story branches that each bump the same issue conflict in its front matter, whose count, last_reported, and updated lines both rewrite
```

The fix:

```text
S-0326 (ADR-0126): the rebase in flai stream sync, flai task done, and flai accept merges an issue file both sides changed by its instances, adding up the count and taking the later last_reported and updated, and the trial merge leaves issue files out. S-0297 confirmed it covers every instance; tests in flai/cmd/stream_sync_test.go reproduce each: TestSyncMergesAnIssueBothSidesBumped and TestSyncTrialMergeLeavesOutIssueFiles (S-0278, S-0221), TestAcceptMergesAnIssueBothSidesBumped (S-0254, S-0228), and TestSyncMergesABumpIntoTheIssueMainClosed (S-0275).
```

Source: I-0092.
