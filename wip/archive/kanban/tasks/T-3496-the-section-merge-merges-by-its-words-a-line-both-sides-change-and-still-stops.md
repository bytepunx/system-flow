---
id: T-3496
type: task
nature: improvement
title: The section merge merges by its words a line both sides change, and still stops on the same word or the same point
status: done
parent: S-0817
owner: alex
created: 2026-11-17T02:31:40Z
updated: 2026-11-17T02:56:42Z
transitions:
  - to: ready
    at: 2026-11-17T02:36:04Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T02:36:05Z
    by: agent-S-0817
  - to: done
    at: 2026-11-17T02:56:42Z
    by: agent-S-0817
stream: S-0817
tags: [flai, merge]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go]
after: [T-3495]
usage:
  source: log
  seconds: 1237
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 10661
      cache_read: 1538608
      cache_write: 80987
      cost: 1.028
---
# T-3496 The section merge merges by its words a line both sides change, and still stops on the same word or the same point

## Work

In `flai/internal/storygit/mdmerge.go`, extend the settle step `mergeLines` passes to `merge3`:

- Prose and list items: keep `mergeSentences`. Where the sentence tokens' own merge finds both sides changing one sentence, settle that run of sentences by words instead of reporting false.
- Table rows: a region whose lines are table rows on all three sides (`|` first), or a mix of table rows and prose, merges by words. Headings, fence lines, and lines inside a fence still report false.
- Words: cut lines into tokens that join back into them, as `sentences` does: each run of characters that are not spaces or tabs, each run of spaces and tabs, and the line break ending each line.
- Run `merge3` on the word tokens with no settle. Insertions both sides make at the same point must stop at this level, not concatenate as line insertions do: give `merge3` a way to refuse a same-point insertion, used for words only, so the line and sentence merges keep their contract (ours' lines first).

In `flai/internal/storygit/mdmerge_test.go`, add cases to `TestMergeMarkdownSettles`: S-0797's table row, where each side inserts words in a different place of one cell; a row where each side changes a different cell; one sentence where each side inserts words at a different place, as in `docs/users/flai.md`. Add cases to `TestMergeMarkdownLeavesWhatItCannotSettle`: both change the same word of one sentence; both insert at the same point of one table row; a heading both change in different words; a fenced line both change in different words. Add a test that the word tokens join back into their lines, beside `TestSentencesCutsTokensThatJoinBack`. Cite I-0387 and the ADR in the test comments.

Waits for T-3495: it builds what that ADR decides.

## Done when

- `MergeMarkdown` settles a table row and a sentence both sides changed in different words, keeping both changes.
- It still errors with `ErrNotSettled` on the same word, a same-point insertion, a heading, and a fenced line.
- The existing cases of both tests pass unchanged.
- `go test ./internal/storygit/` passes, and `flai test` on the two files is clean.

## Notes
