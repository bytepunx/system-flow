---
id: T-3094
type: task
nature: improvement
title: "ADR: the section merge merges by sentences the lines of a paragraph both sides change, refining ADR-0175 and ADR-0276"
status: backlog
parent: S-0719
owner: alex
created: 2026-11-11T08:54:32Z
updated: 2026-11-11T08:54:32Z
transitions: []
stream: S-0719
tags: [adr, design]
touches: [design/adrs/0286-a-section-merge-merges-a-paragraph-both-sides-change-by-its-sentences.md, design/adrs/README.md]
---
# T-3094 ADR: the section merge merges by sentences the lines of a paragraph both sides change, refining ADR-0175 and ADR-0276

## Work

Write the ADR with `flai adr new`, refining ADR-0175 and ADR-0276, from I-0346's instance. Context: S-0688 and S-0666 each changed a different sentence of one paragraph in `design/system/agent-adapters.md` and `docs/operators/index.md`; this repository writes a paragraph on one line, so `mergeLines` saw both change the same base line and the acceptance of S-0688 stopped and sent it back (TH-1061), though the hand resolution was mechanical.

Decision: where the line merge of a section finds the two sides' changes touching the same base lines, it takes the base lines those changes span, and each side's lines for them, and merges them again by sentences with the same three-way merge on sentence tokens. A sentence ends at `.`, `?`, or `!` followed by a space or the line's end; the space and each line break stay their own tokens, so the merged text keeps both sides' spacing and lines. Both changes are kept when they touch different sentences; insertions at one place come ours' first, then theirs', as for lines. Both changing one sentence still leaves the stop. Only prose lines and list items are merged by sentences: a line inside a fence, a table row, or a heading keeps the line merge. It runs wherever `MergeMarkdown` runs, so `flai stream sync`, `flai task done`'s sync, `flai accept`, and the merge queue's batch all settle it. The trial merge still reports the pair, as ADR-0276 leaves it.

Alternatives: rebasing the second story as soon as the first is accepted (already the rule on `overlapped`, and no agent runs on a story in review); asking writers to break paragraphs into one sentence a line (changes nothing already written, and fights the lint); a word-level merge (settles more, but merges edits within one sentence that a reader should judge).

`flai adr new` gives the number; rename the touch with `flai touches` when it is not 0286. Run `flai adr index` so `design/adrs/README.md` lists it.

## Done when

- The ADR is accepted in `design/adrs/`, refines ADR-0175 and ADR-0276, cites I-0346, and names the alternatives.
- `design/adrs/README.md` lists it, and `flai check --strict` and the markdown lint pass on both files.

## Notes
