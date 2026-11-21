---
id: TH-1310
title: S-0809 passed review, but the acceptance stops on the D-0671 lint error that TH-1309 names
anchor:
  path: wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md
  item: S-0809
status: resolved
participants: [orchestrator]
created: 2026-11-18T09:20:48Z
updated: 2026-11-18T13:22:56Z
---

# TH-1310 S-0809 passed review, but the acceptance stops on the D-0671 lint error that TH-1309 names

On wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md.

## Entries

### 2026-11-18T09:20:48Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0671-a-workitem-repo-method-that-clones-items-itself-as-childrenof-does-must-give.md` in the main checkout, as TH-1309 says, then resolve this thread. I accept S-0809 and S-0812 on my next run. The other way: run `flai accept S-0809 --by alex` once the file is fixed.

**What stops the acceptance:** `accept` lints the main checkout's uncommitted markdown before it commits it. One finding is left that `--fix` cannot fix:

```text
design/digest/D-0671-a-workitem-repo-method-that-clones-items-itself-as-childrenof-does-must-give.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
```

**What passed at 902cd57ddb2d09c2ec7cfea76ec8c811049d5cda, the head of `story/S-0809`:**

- `flai verify S-0809 --last`: every step passed, integration and smoke included.
- The verifier found all three criteria met. Every changed file is within the story's touches.
- The accept dry run named no blocker.

S-0809 stays in review. I changed nothing.

### 2026-11-18T13:22:56Z orchestrator
Resolved: D-0671 now lints clean on main (its blockquotes are joined by a > line). S-0809 was rebased since and passed its close-out at c4d796d; the verifier found all three criteria met there.
