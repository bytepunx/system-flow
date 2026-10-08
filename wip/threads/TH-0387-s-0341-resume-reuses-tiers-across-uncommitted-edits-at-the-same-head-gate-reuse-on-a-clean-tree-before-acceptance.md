---
id: TH-0387
title: "S-0341: resume reuses tiers across uncommitted edits at the same head; gate reuse on a clean tree before acceptance"
anchor:
  path: wip/kanban/stories/S-0341-flai-verify-resumes-at-the-tier-that-failed-when-the-branch-head-and-its-base-are-unchanged-so-a-retry-re-runs-the-failure-and-not-the-tiers-that-passed.md
  item: S-0341
status: open
participants: [orchestrator]
created: 2026-10-08T09:43:22Z
updated: 2026-10-08T09:43:22Z
---

# TH-0387 S-0341: resume reuses tiers across uncommitted edits at the same head; gate reuse on a clean tree before acceptance

On wip/kanban/stories/S-0341-flai-verify-resumes-at-the-tier-that-failed-when-the-branch-head-and-its-base-are-unchanged-so-a-retry-re-runs-the-failure-and-not-the-tiers-that-passed.md.

## Entries

### 2026-10-08T09:43:22Z orchestrator
Recommendation for agent-S-0341: before I accept, make `flai verify` reuse tiers only when the worktree is clean (`git status --porcelain` empty). Add a test with a dirty tree and one docs sentence. Alternatively, put a hash of the working-tree diff and untracked files into the reuse key.

Where S-0341 stands:

- Verify passed every step at head c7c37c90.
- The verifier matched all five criteria to the diff: `reusable()` in `story.go`, the fingerprint, `--fresh`, `close-out.sh`'s last line, the tests, docs, and the review page.

Why I hold it: this changes the gate I accept on. Under ADR-0093 I accept a story whose verify passed at its branch head, and the reuse key is the commit, the base commit and the tier fingerprint. `ChangedPaths` counts uncommitted files, and a close-out verifies a dirty tree before it commits. So the following passes:

1. Run 1 at head H with uncommitted content A passes gofmt, vet and lint, then fails go-test.
2. The agent changes A to B, still uncommitted at H.
3. Run 2 reuses gofmt, vet and lint, which were checked against A, runs go-test on B, and records "passed at H".

The same happens if the dirty changes are reverted: the record then says H passed tiers that never ran on H's content. No test covers a dirty tree.

Two minor points to fix at the same time:

- `flai/internal/mcpserver/verify.go` line 17 still describes states as "passed, failed, or not-reached". Add `reused` and `reused_from`.
- `Review.svelte.test.ts` has no reused case, though `review.ts` is tested.

Operator: if you judge the edge case acceptable as is, accept S-0341 yourself and file the dirty-tree gate as a follow-up.
