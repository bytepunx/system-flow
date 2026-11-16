---
id: I-0385
title: flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story
class: defect
status: open
count: 1
cost: 45m
first_reported: 2026-11-16T01:39:37Z
last_reported: 2026-11-16T01:39:37Z
updated: 2026-11-16T01:39:43Z
---

# I-0385 flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story

## Description
flai stream sync reports a rebased branch undoes main's line when the branch removes a closing brace elsewhere in a Go file main added a closing brace to, and the merge queue then never batches the story

## Instances

### 2026-11-16T01:39:37Z
S-0739's sync on 2026-11-16T00:42Z recorded flai/internal/workitem/board.go as undoing one line main added since 4414db3: a tab and a closing brace. Main (S-0790) had added LockBoard, UpdateBoard, NormalizeOrder, and unparsedStories to the file, whose lines include several such braces. S-0739 replaces the loop in LoadBoard that copied DefaultWIPLimits, so its diff against main removes that loop's closing brace, a line main did not add. The undone check compares the lines as a multiset of line text (ADR-0155), so the two braces matched. Every line main added to the file since the story's first base is in the branch's file as often as in main's, and git range-diff of the story's patch before and after the rebase differs only by the import main now has. The record stays while the branch removes any such line, so flai queue refuses the story as undone_main, and the two remedies offered are wrong for it: flai stream restore would commit main's file over the story's change, and --keep-removal lands it by hand outside the queue.

## Remediation

The undone check does not count a line of main's that carries no change of its own. Proposed, for the story to settle in an ADR refining ADR-0155: a line main added is undone only when the branch removes it where main added it, by comparing the hunks' positions in main's file rather than the file's lines as one multiset; or, more narrowly, a line that holds only brackets, braces, or white space is never counted alone. The revert ADR-0155 was written for (I-0124) must still be caught.

Story S-0800 remediates this issue, created from it at 2026-11-16T01:39:43Z.
