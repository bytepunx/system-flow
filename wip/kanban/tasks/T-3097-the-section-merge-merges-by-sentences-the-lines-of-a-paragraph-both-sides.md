---
id: T-3097
type: task
nature: improvement
title: The section merge merges by sentences the lines of a paragraph both sides change, and still stops when both change one sentence
status: done
parent: S-0719
owner: alex
created: 2026-11-11T08:54:43Z
updated: 2026-11-11T14:35:19Z
transitions:
  - to: ready
    at: 2026-11-11T14:24:27Z
    by: agent-S-0719
  - to: in-progress
    at: 2026-11-11T14:24:27Z
    by: agent-S-0719
  - to: done
    at: 2026-11-11T14:35:19Z
    by: agent-S-0719
stream: S-0719
tags: [flai, go]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go]
after: [T-3094]
usage:
  source: log
  seconds: 651
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 31
      output: 181
      cache_read: 1307996
      cache_write: 77335
      cost: 0.5129
---
# T-3097 The section merge merges by sentences the lines of a paragraph both sides change, and still stops when both change one sentence

## Work

In `flai/internal/storygit/mdmerge.go`, where `mergeLines` finds two hunks that `touches` reports as overlapping, gather the run of overlapping hunks from both sides into one region of base lines, and the lines each side has for it. When every line of the region on all three sides is prose or a list item (not a fence line or a line inside a fence, not a table row starting with `|`, not a heading), cut each side's region into sentence tokens as the ADR says, the spaces and line breaks their own tokens, and merge the tokens with the same three-way merge `mergeLines` runs on lines (it takes any `[]string`). Join the merged tokens back into lines in place of the region. When the token merge also finds both sides changing one sentence, or the region holds a line that is not prose, report false as today, so `mergeSections` still leaves "both sides change the same lines".

Keep `mergeLines`' own contract for everything else: changes to different lines, insertions at one place (ours' first), and changes made alike are merged as before.

In `flai/internal/storygit/mdmerge_test.go`, add cases to `TestMergeMarkdownSettles`: a one-line paragraph where ours changes its first sentence and theirs its third, as I-0346's instance did; one side adds a sentence and the other changes another; a list item changed in two sentences. Add cases to `TestMergeMarkdownLeavesWhatItCannotSettle`: both change the same sentence of one paragraph; both change one table row in different cells; both change one line inside a fence. Cite I-0346 and the ADR in the test comments.

Waits for T-3094: it builds what that ADR decides.

## Done when

- `MergeMarkdown` settles a paragraph both sides changed in different sentences, keeping both changes, and still errors with `ErrNotSettled` when both change one sentence, a table row, or a fenced line.
- `go test ./internal/storygit/` passes, and `scripts/flai-test.sh` lint is clean for the package.

## Notes
