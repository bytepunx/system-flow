---
id: TH-1316
title: S-0809 passed review, but the acceptance stops on the same MD028 lint error, now in D-0672, an untracked digest entry on main
anchor:
  path: wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md
  item: S-0809
status: resolved
participants: [orchestrator, alex]
created: 2026-11-18T13:23:16Z
updated: 2026-11-21T17:20:21Z
---

# TH-1316 S-0809 passed review, but the acceptance stops on the same MD028 lint error, now in D-0672, an untracked digest entry on main

On wip/kanban/stories/S-0809-the-merge-queue-s-batch-runs-the-suites-verify-at-batch-lists-over-the-merge-of.md.

## Entries

### 2026-11-18T13:23:16Z orchestrator
**Recommended:** in the main checkout, fix line 20 of `design/digest/D-0672-flai-internal-workitem-cannot-import-flai-internal-verify-since-verify-imports.md`. Write `>` on the blank line between its two blockquotes. Then resolve this thread, and I accept S-0809 on my next run. The other way: run `flai accept S-0809 --by alex` once the file is fixed.

**What stops the acceptance:** `accept` lints the main checkout's uncommitted markdown before it commits it with the acceptance. One finding is left that `--fix` cannot fix:

```text
design/digest/D-0672-flai-internal-workitem-cannot-import-flai-internal-verify-since-verify-imports.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
```

D-0671 is fixed now, so I resolved TH-1310. The curator wrote D-0672 after that with the same pattern: two quoted sources separated by a blank line. Unless `flai digest` stops writing that pattern, every new entry with two quotes will stop acceptances again. That is worth an issue.

**What passed at c4d796d9a17d58cd7d377a6ea8bdb94caeccdf2c, the head of `story/S-0809`:**

- `flai verify S-0809 --last`: every step passed, integration and smoke included.
- The verifier found all three criteria met. All 18 changed files are within the story's touches.
- The accept dry run named only my own TH-1310, which is now resolved.

S-0809 stays in review. The acceptance stopped before it changed anything.

### 2026-11-21T17:20:21Z alex
Resolved: S-0809 was accepted
