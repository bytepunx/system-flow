---
id: ADR-0297
title: "A markdown line of main's that a rebased branch removes is not undone when a line the branch adds holds every word of it"
status: accepted
date: 2026-11-13
supersedes: []
superseded_by: []
refines: [ADR-0155]
---

# ADR-0297 A markdown line of main's that a rebased branch removes is not undone when a line the branch adds holds every word of it

## Context

[ADR-0155](0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md) has `flai stream sync` compare, for each path main changed since the branch's old base, the lines main added with the lines the rebased branch's net change removes, as a multiset of line text. A path where the branch removes any of them is recorded, and it holds the story's acceptance until the operator restores the path, keeps the removal with a reason, or sends the story back.

This repository writes each paragraph and each table row on one line. A branch that adds a sentence to a line main changed removes main's line and adds a longer one. The whole-line compare counts that as undoing main's change. [I-0360](../issues/I-0360-flai-stream-sync-reports-a-rebased-branch-undoes-main-s-line-when-both-changed.md) records the cost. S-0735 and S-0761 both changed the `flai check` row of § Commands in `design/system/flai-cli.md`. S-0735 was accepted first, and S-0761's acceptance stopped on that row, which the section merge of [ADR-0287](0287-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences-keeping.md) does not settle, and was sent back on TH-1147. Its agent resolved the stop with main's row plus S-0761's sentence and list entry. The sync still warned that the branch removes 1 line main added to the file, and recorded it. The acceptance would then wait for `--keep-removal`, and restoring the path would drop the story's own edit.

## Decision

In a markdown path, a line of main's that the rebased branch removes is kept, not undone, when a line the branch adds holds every word of it.

- **Markdown.** A path is markdown when its extension is `.md`, in any case.
- **Words.** A line is split into words on white space, and the punctuation at each word's ends is left out, so `ADR-0256),` and `ADR-0256,` are one word. A word of punctuation alone is none. Case counts.
- **Every word.** The branch's line holds main's words as a multiset: each word at least as often as main's line has it, in any order. A sentence added at the end, a list entry added in the middle, a word moved, or a bracket dropped keep main's line. A word dropped or changed does not.
- **Once.** Each line the branch adds keeps at most one of main's lines. Main's lines are taken in order, and each is kept by the first line the branch adds that holds it and has kept none yet.
- **Not a revert.** A line the branch adds that is one of the lines main removed since the branch's old base keeps none: the branch restores main's old line, which is the revert ADR-0155 catches, even where main's change only took words out.
- **Blank lines.** A line of main's with no words is never kept by the rule.
- **Other paths.** Every other path keeps the multiset of whole lines, because in code a line can keep every word and still undo main's change, as a line commented out does.
- **The record.** `RecordedUndone` applies the same rule, reading the lines main removed since the recorded old base. So a path recorded before this rule, such as S-0761's, drops out once the branch's line holds main's words.

It applies wherever `FindUndone` and `RecordedUndone` are read: `flai stream sync`, `flai task done`'s sync, `flai accept`'s rebase and its preview, and the merge queue's eligibility. ADR-0155's record, blocker, restore, keep, and send back are unchanged, and so is its case of a revert replayed over main's change.

## Consequences

- A story that extends a paragraph or a table row main changed is no longer reported as undoing main's change, and its acceptance no longer waits for `--keep-removal`. S-0761's row is kept: it moved one `and` of main's to after its new list entry and turned `ADR-0256),` into `ADR-0256,`, and every word of main's row is still there.
- A branch that inserts words which reverse main's meaning, such as a `not`, or that reorders main's words, keeps every word and is not reported. The story's own review is the check.
- A branch that changes only punctuation main added, such as a comma, is not reported.
- A branch that joins two of main's lines into one keeps only the first. The second is reported.
- Reading the record costs one more `git diff` per recorded markdown path whose branch adds a line.

## Alternatives considered

- **Diff main's line against the branch's by words, as ADR-0287 merges by sentences.** It would name the words removed. But the net diff does not pair each of main's lines with the line that replaced it, and the test of main's words against each added line answers the same question without the pairing. Lost on simplicity.
- **Compare by sentences.** ADR-0287 does not cut a table row into sentences, and the line in I-0360 was a table row. Lost.
- **Require main's line as one unbroken run of the branch's line.** S-0761 added a list entry inside main's row, which splits main's words into two runs. It would still be reported. Lost.
- **Require main's words in order, as a subsequence of the branch's line.** It was the first plan. S-0761's row moved one of main's words and dropped a bracket, so 3,555 of main's 3,970 words were in order and the row would still be reported. Lost.
- **Apply the rule to every path.** In code, `// x := 1` holds every word of `x := 1` and undoes it. Lost on correctness.
- **Leave it to `--keep-removal`.** Every story that extends main's line would wait for the operator at acceptance, for a removal that removes nothing, and restoring the path would drop the story's edit. Lost.
