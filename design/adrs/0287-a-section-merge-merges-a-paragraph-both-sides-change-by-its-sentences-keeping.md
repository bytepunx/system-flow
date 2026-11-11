---
id: ADR-0287
title: "A section merge merges a paragraph both sides change by its sentences, keeping both changes when they touch different sentences"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0175, ADR-0276]
---

# ADR-0287 A section merge merges a paragraph both sides change by its sentences, keeping both changes when they touch different sentences

## Context

[ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md) merges a markdown file a rebase stops on by sections, and [ADR-0276](0276-a-rebase-that-stops-on-a-markdown-file-under-the-manifest-s-design-or-docs.md) widened it to every markdown file under the design and docs folders. A section both sides change is merged line by line, and both sides changing the same base lines leaves the stop for the agent.

This repository writes each paragraph on one line. Two stories that change different sentences of one paragraph therefore change the same base line, and the line merge cannot keep both. [I-0346](../issues/I-0346-a-trial-merge-that-finds-two-open-stories-editing-different-sentences-of-one.md) records the cost: S-0688 and S-0666 each changed a different sentence of one paragraph in `design/system/agent-adapters.md` and in `docs/operators/index.md`. The two agents agreed that the second to land keeps both. The acceptance of S-0688 after S-0666 still stopped, sent S-0688 back on TH-1061, and cost a hand rebase, a second close-out, and a second acceptance. The resolution was mechanical: main's paragraph with S-0688's sentence applied.

## Decision

Where the line merge of a section finds the two sides' changes touching the same base lines, it merges those lines again by their sentences, and keeps both changes when they touch different sentences.

- The region is the run of base lines that the overlapping changes of both sides span, with each side's lines for it.
- It is merged by sentences only when every line of the region, on all three sides, is prose or a list item. A heading, a table row (a line that starts with `|`), a fence line, or a line inside a fenced code block keeps the line merge, and the stop.
- Each side's region is cut into tokens: a sentence ends at `.`, `?`, or `!` followed by a space or the line's end; each run of spaces and each line break is a token of its own. So the merged text keeps both sides' spacing and line breaks.
- The tokens are merged by the same three-way merge the lines are: changes to different tokens are both kept; insertions at one place take ours' tokens, then theirs'; a change both make alike is kept once.
- Both sides changing or removing the same sentence still leaves the file not settled, with the stop for the agent, as before.

It runs wherever `MergeMarkdown` runs: `flai stream sync`, `flai task done`'s sync, `flai accept`, and the merge queue's batch.

What stays as ADR-0175 and ADR-0276 have it:

- The trial merge (`git merge-tree`) still reports a conflict between two open stories on such a paragraph, and they are still told in their conversation. Making the trial merge run the section merge stays deferred, as ADR-0276 says.
- Every other case ADR-0175 cannot settle is unchanged.

## Consequences

- Two stories that change different sentences of one paragraph of a design or docs document, or of a shared markdown file, no longer stop at sync or acceptance, and the second is not sent back for it.
- A sentence one side adds next to a sentence the other changes lands beside it, ours first at one place. The merged paragraph may read oddly where the two edits depend on each other. The story's own review and the markdown lint are the check, as ADR-0175 says.
- An abbreviation such as `e.g.` ends a token, so a sentence may be cut finer than a reader would cut it. Two changes inside one such sentence can then both be kept. This is rare in this repository's prose, and is accepted.

## Alternatives considered

- Rebase the second story as soon as the first is accepted, on the `overlapped` change. That is already the rule in `work-management.md`; it did not help because S-0688 was in review, with no agent running. Lost.
- Write one sentence a line. It changes nothing already written, and it is not this repository's style. Lost.
- Merge by words. It would settle more, but it would also merge two edits inside one sentence, which a reader should judge. Lost on correctness.
- Make the trial merge run the section merge. Deferred, as ADR-0276 deferred it: it changes what holds a ready story and what the stories are messaged.
