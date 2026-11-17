---
id: S-0817
type: story
nature: improvement
title: Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle
status: review
owner: alex
created: 2026-11-17T02:26:33Z
updated: 2026-11-17T03:23:24Z
transitions:
  - to: ready
    at: 2026-11-17T02:33:25Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-17T02:33:33Z
    by: agent-S-0817
  - to: review
    at: 2026-11-17T03:23:24Z
    by: agent-S-0817
tags: []
topics: [merge, sync, acceptance]
touches: [flai/internal/storygit/mdmerge.go, flai/internal/storygit/mdmerge_test.go, flai/cmd/accept_conflict_test.go, flai/cmd/stream_sync_test.go, design/adrs/0314-a-section-merge-merges-a-line-both-sides-change-by-its-words-keeping-both.md, design/adrs/README.md, design/system/workflow.md, design/system/agent-coordination.md, design/system/flai-cli.md, docs/users/flai.md, design/issues/I-0387-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md, design/issues/summary.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 3006
  priced_by: harness
  turns:
    - day: 2026-11-17
      ceremony: 1
      work: 41
  models:
    - model: claude-opus-5-5
      input: 222
      output: 85192
      cache_read: 12294950
      cache_write: 647167
      cost: 8.215
  strategic:
    - kind: planner
      seconds: 283
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 52
          output: 19991
          cache_read: 3454967
          cache_write: 172382
          cost: 2.4701
    - kind: orchestrator
      seconds: 181
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 30
          output: 217
          cache_read: 1234059
          cache_write: 192211
          cost: 0.4391
        - model: claude-sonnet-5-5
          input: 10
          output: 79
          cache_read: 133628
          cache_write: 42516
          cost: 0.141
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
cost_of_delay:
  inputs:
    time_lost_per_cycle: 4h
    by: flai
    at: 2026-11-17T02:26:33Z
  value: 600
  by: planner-S-0817
  at: 2026-11-17T02:31:15Z
forecast:
  duration: 1h
  delivery: 2026-11-17T12:19:00Z
  basis: "Its own forecast of 1h; 17th in the pull order with an in-progress limit of 3, behind S-0750, S-0805, S-0241, S-0752, S-0753, S-0754, S-0756, S-0758, S-0806, S-0807, S-0808, S-0809, S-0810, S-0811, S-0812, S-0813, S-0814 and S-0815."
  by: flai
  at: 2026-11-17T02:32:48Z
finalized:
  by: orchestrator
  at: 2026-11-17T02:33:17Z
---
# S-0817 Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle

## Goal

This story remediates [I-0387](../../../design/issues/I-0387-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md), "Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle", which the analysis in [2026-11-16-all.md](../../../design/analysis/2026-11-16-all.md) found. The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0387 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0387 is closed with `flai issue close I-0387 --reason` saying what fixed it

## Tasks
- T-3495 ADR: the section merge merges by its words a line both sides change, refining ADR-0287
- T-3496 The section merge merges by its words a line both sides change, and still stops on the same word or the same point
- T-3497 Say in the design and the user docs that the section merge settles a line both sides change in different words
- T-3498 Reproduce I-0387 in real git: two stories insert words at different places of one table row and one sentence, and neither sync nor acceptance stops
- T-3499 Close I-0387, saying the section merge now settles a line two stories change in different words

## Notes

The analyzer found I-0387: the analysis the goal links gives its evidence.

Cost of delay inputs set by flai from I-0387. time_lost_per_cycle 4h carried over from I-0387's Impact section. Its Impact time_lost_per_cycle was taken rather than the 36m derived from its cost and count.

Strategic usage carried over from I-0387: analyzer 1m39s, 1.2588 USD.

### Planning

The proposed fix, from I-0387's two instances and `.flai-cache/conflicts.jsonl`:

- The first instance, E-0033's stories on `design/system/server-performance.md`, stopped acceptances last on 2026-11-09. S-0719 (ADR-0287, accepted 2026-11-11) made the section merge settle a paragraph both sides change in different sentences. No acceptance has stopped on that file since.
- Every acceptance stop since ADR-0287 is on `design/system/flai-cli.md`: S-0689 (2026-11-11), S-0761 (2026-11-13), and S-0797 (2026-11-17, with `docs/users/flai.md`). S-0797's stop is the second instance. Main and the branch each inserted words at different places in one long line, the `flai queue` row of the commands table and one sentence of the `--json` paragraph, removing nothing. `mergeLines` in `flai/internal/storygit/mdmerge.go` leaves a table row to the line merge (`proseLines`), and `mergeSentences` stops when both change one sentence. So both report "both sides change the same lines".
- The fix: where both sides change the same line, and the sentence merge cannot settle it or the line is a table row, merge that line again by its words. Cut each side's line into words, runs of spaces, and the line break, and run the same three-way merge on the tokens. Keep both changes when they touch different words. Still stop when both change or remove the same word, or when both insert at the same point between two words, since no order of the two insertions is right. Headings and fenced lines keep the line merge and its stop. The one change settles the stop in every rebase flai runs: `flai stream sync`, `flai task done`'s sync, `flai accept`, rebase and resolve, and the merge queue's batch all call `MergeMarkdown`.
- The merged line holds every word of main's line, so ADR-0297 already keeps it from reading as undone. The reproduction checks that.
- The trial merge (`git merge-tree`) still counts these files as a conflict, as ADR-0276 leaves it.

Touches, and where each came from. The story declared none. `flai touches suggest`, seeded with the merge code and its test, listed no co-change above its threshold. The rest follow S-0719's footprint, the same change one level down:

| Touch | Source | Why |
|---|---|---|
| `flai/internal/storygit/mdmerge.go` | layout | `mergeLines`, `proseLines`, and `mergeSentences`, where the stop is decided |
| `flai/internal/storygit/mdmerge_test.go` | layout | `TestMergeMarkdownSettles` and `TestMergeMarkdownLeavesWhatItCannotSettle` |
| `flai/cmd/accept_conflict_test.go` | layout | the acceptance reproduction, beside S-0719's of I-0346 |
| `flai/cmd/stream_sync_test.go` | layout | the sync reproduction, beside S-0719's |
| `design/adrs/0314-a-section-merge-merges-a-line-both-sides-change-by-its-words-keeping-both.md` | design | the ADR refining ADR-0287; `flai adr new` gives the number and slug, so the agent renames this touch when it differs |
| `design/adrs/README.md` | co-change | generated from the ADR files |
| `design/system/workflow.md` | design | § Branches and collisions (ADR-0019) says a same-sentence or table-row change stops |
| `design/system/agent-coordination.md` | design | its ADR-0175 paragraph |
| `design/system/flai-cli.md` | design | the `flai stream sync` and `flai accept` rows name the same-lines stop |
| `docs/users/flai.md` | design | the Shared documents bullet under sync |
| `design/issues/I-0387-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md` | declared by the goal | closed by the last task |
| `design/issues/summary.md` | co-change | written again when the issue closes |

No folder touch. No `after`: the story reads no open story's output. ADR-0287 and its code are on main already.

- Forecast 1h, delivery 2026-11-17T11:44:00Z. `flai forecast` gave 2m from size 2 (2 criteria, 0 touches) at 47 s a unit, before any touch was declared. That size does not count an ADR, new merge code, two real-git reproductions, and four documents. S-0719, the same shape, took 42m of agent time against a 1h forecast. This one adds table rows and the same-point refusal, so 1h. The delivery is flai's 10:46Z moved by the 58m added.
- Cost of delay 600.00 USD a week, `flai cod`'s figure from the operator's input (4h lost a cycle at 150 USD an hour), kept. I-0387 counts two occurrences, and S-0797's was this week.
