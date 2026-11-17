---
id: TH-1269
title: S-0742 is ready to accept, but an uncommitted digest entry on main fails lint (MD028)
anchor:
  path: wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md
  item: S-0742
status: open
participants: [orchestrator]
created: 2026-11-17T03:09:18Z
updated: 2026-11-17T03:09:18Z
---

# TH-1269 S-0742 is ready to accept, but an uncommitted digest entry on main fails lint (MD028)

On wip/kanban/stories/S-0742-flai-serve-starts-restarts-and-stops-agents-plans-again-and-finds-its-reasons.md.

## Entries

### 2026-11-17T03:09:18Z orchestrator
**Recommended:** fix line 20 of `design/digest/D-0653-workitem-s-epicfollows-planof-newholds-plancandidatesof-promotioncandidates.md` in the main checkout. Then I accept S-0742 on my next run.

**What is wrong:** the file is the curator's new entry. It is untracked. It holds two blockquotes with a blank line between them, which markdownlint flags and its `--fix` cannot fix:

```text
design/digest/D-0653-…md:20 error MD028/no-blanks-blockquote Blank line inside blockquote
```

flai's acceptance commits the main checkout's uncommitted files under `design/` with the acceptance, so it refused:

```text
S-0742 cannot be accepted: the main checkout's uncommitted files hold lint the fix of the markdown tier could not fix … correct them in the main checkout and run flai accept S-0742 again
```

**The fix:** put text between the two quotes, for example "It adds:", or turn the blank line into a line holding only `>` so the two become one quote. I cannot edit files.

**S-0742 itself is fine:**

- The close-out passed every step at the branch head, `412724d1`, integration and smoke included.
- The verifier matched all five criteria to changed files.
- The dry run listed no blocker.

**Alternative:** accept it yourself with `flai accept S-0742 --by alex` once the file is fixed.
