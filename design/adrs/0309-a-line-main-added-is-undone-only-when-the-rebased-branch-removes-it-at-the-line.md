---
id: ADR-0309
title: "A line main added is undone only when the rebased branch removes it at the line main added it, by its position in main's file, not when the branch removes a line of the same text elsewhere"
status: accepted
date: 2026-11-16
supersedes: []
superseded_by: []
refines: [ADR-0155, ADR-0297]
---

# ADR-0309 A line main added is undone only when the rebased branch removes it at the line main added it, by its position in main's file, not when the branch removes a line of the same text elsewhere

## Context

[ADR-0155](0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md) has `flai stream sync` compare, for each path main changed since the branch's old base and the rebased branch still changes, the lines main added with the lines the branch's net change removes, as a multiset of line text. [ADR-0297](0297-a-markdown-line-of-main-s-that-a-rebased-branch-removes-is-not-undone-when-a.md) keeps a markdown line of main's that a line the branch adds holds every word of. ADR-0155 named the cost of a compare by text: a line main added that the branch removes for another reason, such as a blank line or a brace, is reported too.

[I-0385](../issues/I-0385-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-the-branch.md) records that cost. S-0739 replaces a loop in `LoadBoard` in `flai/internal/workitem/board.go`, so its diff against main removes the loop's closing brace, a tab and a `}`. Main, through S-0790, had added four functions to the same file, each ending in the same line. The two braces matched by text, the sync recorded the path as undoing main's change, and the merge queue refused the story as `undone_main`. Neither remedy fits: `flai stream restore` would put main's file over the story's change, and `--keep-removal` lands the story by hand, outside the queue. Any line with no change of its own, such as `}`, `)`, `return nil`, `default:`, or a blank line, can match the same way in code.

## Decision

A line main added since the branch's old base is undone only when the rebased branch's net change removes that line, at the position main added it, in main's file.

- **Positions.** Both diffs are read with `-U0`, and each hunk header, `@@ -a,b +c,d @@`, gives the line numbers of what it removes in the old file and of what it adds in the new one. Main's diff, from the old base to `onto`, the commit of main the rebase took, gives the lines main added by their numbers in main's file at `onto`. The branch's diff, from `onto` to its head, gives the lines it removes by their numbers in the same file. A line main added is undone when the branch removes the line of that number. A path the branch deletes removes every line, so every line main added is undone, as before.
- **The record.** `RecordedUndone` checks the record against main now: main's diff from the recorded old base to main, and the branch's diff from main to the branch, both by position in main's file now. Of what that finds, it keeps only the record's lines, counted as a multiset of their text, as before. So a line only a commit after `onto` adds, the case of I-0338, is still not counted unless the sync recorded it, and a path recorded by the text compare, such as S-0739's, drops out at its next read once the branch removes none of main's lines where main added them. Reading the record now costs main's diff for every recorded path, not only for a markdown path.
- **Markdown.** ADR-0297's rule applies to what the position rule leaves: a markdown line of main's the branch removes where main added it is still kept when a line the branch adds holds every word of it, and not when that line is one main removed.
- **The revert.** The revert of I-0124 is still caught. The branch's revert, replayed over main's identical change, removes the very lines main added, so its removals fall at the numbers main's additions hold in main's file.

It applies wherever `FindUndone` and `RecordedUndone` are read: `flai stream sync`, `flai task done`'s sync, `flai accept`'s rebase and its preview, and the merge queue's eligibility. ADR-0155's record, blocker, restore, keep, and send back are unchanged.

## Consequences

- A branch that removes a brace, a blank line, or another common line of its own, in a file where main added a line of the same text, is no longer reported, and the merge queue batches it.
- A branch that removes a line main added and adds the same text at another place in the file removes main's line where main added it, and is still reported.
- Git chooses where a run of lines that repeats is added or removed, so a diff can name a brace one line above the one it moved. Where main's addition and the branch's removal meet such a run, a line can still be reported, or left out, by one line. The revert of a whole change of main's is still reported, since most of its lines meet main's.
- Reading the record runs one more `git diff` per recorded path that is not markdown.

## Alternatives considered

- **Never count a line that holds only brackets, braces, or white space alone** (the issue's narrower option). It settles I-0385's brace, but not `return nil`, `default:`, `err = nil`, or a comment line repeated in a file, and it would miss a revert of main's change that was only braces or blank lines. Lost to the position rule, which settles every such line.
- **Count by text within main's hunks only.** It still matches a line the branch removes in a hunk of its own to one main added elsewhere. Lost.
- **Leave it to `--keep-removal`.** Every story that removes a brace in a file main extended would land by hand, outside the queue, for a removal that undoes nothing. Lost.
