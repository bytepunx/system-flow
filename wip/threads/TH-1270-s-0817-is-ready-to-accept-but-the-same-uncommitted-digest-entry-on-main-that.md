---
id: TH-1270
title: S-0817 is ready to accept, but the same uncommitted digest entry on main that holds S-0742 fails lint (MD028)
anchor:
  path: wip/kanban/stories/S-0817-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md
  item: S-0817
status: resolved
participants: [orchestrator, alex]
created: 2026-11-17T03:26:10Z
updated: 2026-11-17T04:40:48Z
---

# TH-1270 S-0817 is ready to accept, but the same uncommitted digest entry on main that holds S-0742 fails lint (MD028)

On wip/kanban/stories/S-0817-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md.

## Entries

### 2026-11-17T03:26:10Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0653-workitem-s-epicfollows-planof-newholds-plancandidatesof-promotioncandidates.md` in the main checkout, as TH-1269 asks for S-0742. Then I accept S-0817 and S-0742 on my next run.

**What stopped it:** flai commits the main checkout's uncommitted `design/` files with the acceptance, and that file fails lint:

```text
design/digest/D-0653-…md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
```

**The fix:** put text between the two quotes, or make the blank line hold only `>`. I cannot edit files.

**S-0817 itself is fine:**

- The close-out passed every step at the branch head, `9b03cee4`, integration and smoke included.
- The verifier matched both criteria to changed files.
- The dry run listed no blocker.

**Alternative:** accept it yourself with `flai accept S-0817 --by alex` once the file is fixed.

### 2026-11-17T04:40:43Z alex
I fixed it

### 2026-11-17T04:40:48Z alex
Resolved.
