---
id: D-0235
title: "Plan for S-0326: five tasks in four layers, folding a branch's duplicate issue after each rebase: Recommendation: take the planner's remedy and approve the plan as it stands"
kind: answer
sources: [TH-0353]
topics: [cli, git, continuous-improvement, flai, go, code, all]
paths: [design/adrs/README.md, flai/internal/issues/fold.go, flai/internal/issues/fold_test.go, flai/internal/storygit/sync.go, flai/internal/storygit/sync_test.go, flai/cmd/stream.go, flai/cmd/branch.go, flai/internal/taskdone/taskdone.go, flai/cmd/stream_sync_test.go, design/system/continuous-improvement.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0112-flai-verify-record-issues-opens-an-issue-on-the-story-branch-that-another-branch-opened-under-the-same-title-meanwhile.md, design/issues/summary.md, design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md, flai/internal/issues/merge.go, flai/internal/issues/merge_test.go, flai/internal/issues/generated.go, flai/cmd/stream_sync.go, docs/users/flai-reference.md, design/issues/I-0086-the-close-out-s-install-smoke-test-failed-once-and-passed-when-run-alone-with-no-cause-in-its-output.md]
date: 2026-10-08
status: absorbed
absorbed_by: "design/adrs/0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md#decision"
asked_again: []
cited: []
---

# D-0235 Plan for S-0326: five tasks in four layers, folding a branch's duplicate issue after each rebase: Recommendation: take the planner's remedy and approve the plan as it stands

The operator, alex, answered TH-0353 at 2026-10-08T04:20:59Z.

What was asked:

```text
S-0326 is planned: 15 touches, a forecast of 45m delivered by 2026-10-08T08:44:00Z, and a cost of delay of 25 USD a week. It stays a draft for you to finalize.
```

The answer:

```text
Recommendation: take the planner's remedy and approve the plan as it stands.

The remedy is that once a rebase completes, flai folds an open issue the branch added into the main branch's open issue of the same title, and commits the fold on the branch.

- **It removes I-0112's cause.** The record-time fix cannot work: a close-out cannot bump a file its branch does not have. Folding after the rebase is the first point where both files are visible.
- **It goes beyond ADR-0098.** ADR-0098 lets a sync regenerate one derived file, `summary.md`, and only when it is the sole conflict. Folding merges two source issue files and deletes one, after every completed rebase in sync, task done, and accept. That is a new automatic change to recorded data, so it refines ADR-0085 and ADR-0098 in a way only you can confirm. T-1319 records it.
- **The alternative is to leave folding to people.** flai would only warn on a duplicate title after a sync, and the story's agent or you would fold it by hand. That is safer, but it keeps the duplicates that I-0109/I-0111 and S-0323/S-0325 show.

The plan fits the story:

- T-1321 reproduces I-0112 across two branches (criterion 1).
- T-1323 closes I-0112 (criterion 2).
- T-1320 and T-1322 share no path.
- The 45m forecast and the 25 USD a week value stand.

S-0326 stays a draft until this thread settles.

Source: design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md § Decision
```

Source: TH-0353.
