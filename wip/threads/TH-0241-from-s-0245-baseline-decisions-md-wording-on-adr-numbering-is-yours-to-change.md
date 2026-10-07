---
id: TH-0241
title: "From S-0245: baseline decisions.md wording on ADR numbering is yours to change"
anchor:
  path: wip/kanban/stories/S-0275-issue-and-adr-handling-from-the-worktree-is-one-call-each-bump-close-and-adr-new-number-from-the-whole-repository-and-commit-on-the-story-branch.md
  item: S-0275
status: resolved
participants: [agent-S-0245, alex]
created: 2026-10-07T02:59:12Z
updated: 2026-10-07T06:37:48Z
---

# TH-0241 From S-0245: baseline decisions.md wording on ADR numbering is yours to change

On wip/kanban/stories/S-0275-issue-and-adr-handling-from-the-worktree-is-one-call-each-bump-close-and-adr-new-number-from-the-whole-repository-and-commit-on-the-story-branch.md.

## Entries

### 2026-10-07T02:59:12Z agent-S-0245
From S-0245 (TH-0240): S-0245 made `flai adr new` number past the ADRs in this checkout, on main, and in every story's worktree and branch (`adr.NextNumber` now reads `storygit.FolderNames` too, as `issues.NextID` does). `design/conventions/decisions.md` (baseline) and `template/root/design/conventions/decisions.md` still say it "calculates the next number from the files present". The operator left that line to S-0275, which claims both files. Proposed wording: "calculates the next number from the ADRs in this checkout, on main, and in every story's worktree and branch". No answer needed; resolve this thread when the line is changed.

### 2026-10-07T06:37:48Z alex
Resolved.
