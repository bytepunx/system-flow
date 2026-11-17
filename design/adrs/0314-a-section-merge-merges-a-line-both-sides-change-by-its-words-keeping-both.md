---
id: ADR-0314
title: "A section merge merges a line both sides change by its words, keeping both changes when they touch different words"
status: accepted
date: 2026-11-17
supersedes: []
superseded_by: []
refines: [ADR-0287]
---

# ADR-0314 A section merge merges a line both sides change by its words, keeping both changes when they touch different words

## Context

[ADR-0287](0287-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences-keeping.md) merges again by its sentences a run of prose lines that both sides of a section merge change. A table row keeps the line merge, and both sides changing one sentence still leaves the stop for the agent.

[I-0387](../issues/I-0387-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md) records what is left. Since ADR-0287, every acceptance that stopped did so on `design/system/flai-cli.md`: S-0689, S-0761, and S-0797. In S-0797's stop, main and the branch each inserted words at different places in one table row, the `flai queue` row of § Commands, and in one sentence of the `--json` paragraph, and removed nothing. The table row is left to the line merge, and the sentence is changed by both, so the file was not settled and the acceptance sent the story back. Its agent settled the stop by hand, mechanically: main's line with the branch's words inserted.

ADR-0287 considered merging by words and lost it on correctness: it would merge two edits inside one sentence, which a reader should judge. The stops since show that the edits it leaves to the agent are almost always insertions at different places, settled the same mechanical way, while each costs a send-back, a second close-out, and a second acceptance.

## Decision

Where both sides of a section merge change the same lines, and those lines are prose, list items, or table rows on all three sides, the merge settles them by their words when their sentences do not settle them.

- **The region.** It is the run of base lines that the overlapping changes span, with each side's lines for it, as ADR-0287 has it.
- **Which lines.** Every line of the region, on all three sides, must be prose, a list item, or a table row (a line that starts with `|`). A heading, a fence line, or a line inside a fenced code block keeps the line merge, and the stop.
- **Sentences first.** A region of prose and list items alone is merged by its sentences first, as ADR-0287 says. Only when that does not settle it, or when the region holds a table row, is it merged by its words.
- **Tokens.** Each side's region is cut into tokens: each run of characters that are not spaces or tabs is a word, each run of spaces and tabs is a token, and each line break is a token. So the merged lines keep both sides' spacing and line breaks, and a table row's `|` is a word of its own.
- **The merge.** The tokens are merged by the same three-way merge the lines are: changes to different words are both kept, and a change both make alike is kept once.
- **What still stops.** Both sides changing or removing the same word, or both inserting words at the same point between two words, leaves the file not settled, with the stop for the agent. No order of two insertions at one point is right, so the word merge does not concatenate them, unlike the line and sentence merges, which put ours first.

It runs wherever `MergeMarkdown` runs: `flai stream sync`, `flai task done`'s sync, `flai accept`, rebase and resolve, and the merge queue's batch.

What stays as ADR-0175, ADR-0276, and ADR-0287 have it:

- The trial merge (`git merge-tree`) still reports a conflict on such a line between two open stories, and they are still told in their conversation.
- ADR-0287's sentence merge, and its concatenation of sentences inserted at one place, are unchanged.
- Every other case ADR-0175 cannot settle is unchanged.

## Consequences

- Two stories that insert or change words at different places of one table row, or of one sentence, of a design or docs document, or of a shared markdown file, no longer stop at sync or acceptance, and the second is not sent back for it.
- The merged line holds every word of main's line that the branch did not change, so [ADR-0297](0297-a-markdown-line-of-main-s-that-a-rebased-branch-removes-is-not-undone-when-a.md) does not report it as undoing main's change.
- Two edits inside one sentence can now both be kept, and the sentence may read oddly where they depend on each other, as ADR-0287 warned. The story's own review, the operator's acceptance, and the markdown lint are the check. Edits that meet at one word or one point still stop.
- A table row whose two edits each add a cell keeps both cells, and the row's cell count may then differ from the header's. The markdown lint is the check.
- The token diff of a long line costs the product of its tokens on two sides, as `lineHunks` does for lines: a paragraph of a few thousand words costs tens of megabytes for one merge. It runs only on a region both sides change.

## Alternatives considered

- **Keep ADR-0287 as it is.** Every acceptance stop since it was on a line both changed in different words, and each was settled by hand the same way. Lost on the flow.
- **Merge a table row by its cells, and a sentence by its words, apart.** A row's cells are words between `|` words, so the word merge already keeps two edits to different cells, and an edit inside one cell the other side also edits stops on the same word or point. Lost on simplicity.
- **Concatenate words both sides insert at one point, ours first, as the sentence merge does.** Two phrases put side by side mid-sentence rarely read right. Lost on correctness.
- **Merge by characters.** It would settle edits inside one word, which no stop recorded needed, and would merge a word both sides changed into a word neither wrote. Lost on correctness.
- **Make the trial merge run the section merge.** Deferred, as ADR-0276 and ADR-0287 deferred it.
