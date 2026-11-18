---
id: TH-1309
title: S-0812 passed review, but the acceptance stops on a lint error in D-0671, an uncommitted digest entry on main
anchor:
  path: wip/kanban/stories/S-0812-flai-test-at-ci-runs-the-suites-verify-at-ci-lists-and-the-template-s-and-this.md
  item: S-0812
status: open
participants: [orchestrator]
created: 2026-11-18T09:17:45Z
updated: 2026-11-18T09:17:45Z
---

# TH-1309 S-0812 passed review, but the acceptance stops on a lint error in D-0671, an uncommitted digest entry on main

On wip/kanban/stories/S-0812-flai-test-at-ci-runs-the-suites-verify-at-ci-lists-and-the-template-s-and-this.md.

## Entries

### 2026-11-18T09:17:45Z orchestrator
**Recommended:** fix line 19 of `design/digest/D-0671-a-workitem-repo-method-that-clones-items-itself-as-childrenof-does-must-give.md` in the main checkout. Either write `>` on the blank line between its two blockquotes, or put a line of prose between them. Then resolve this thread, and I accept S-0812 on my next run. The other way: run `flai accept S-0812 --by alex` yourself once the file is fixed.

**What stops the acceptance:** `accept` lints the main checkout's uncommitted markdown before it commits it with the acceptance. One finding is left that `--fix` cannot fix:

```text
design/digest/D-0671-a-workitem-repo-method-that-clones-items-itself-as-childrenof-does-must-give.md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
```

D-0671 is untracked (`??`) on main. It looks like the curator wrote it with `flai digest`. If so, `flai digest` writes two quoted sources separated by a blank line, which breaks MD028. That is worth an issue.

**What passed at cca8f937d0f7d23e7b910ece06e59c1a9912e029, the head of `story/S-0812`:**

- `flai verify S-0812 --last`: every step passed, integration and smoke included.
- The verifier found all four criteria met, and every changed file is within the story's touches.
- The accept dry run named no blocker.

S-0812 stays in review. I changed nothing.
