---
id: S-0429
type: story
nature: feature
title: flai digest harvest lists what the archive and the threads hold since a date, and a one-off seed fills the digest from the operator's answers and the closed issues so there is something to prime with
status: done
parent: E-0025
owner: alex
created: 2026-10-16T05:58:28Z
updated: 2026-10-23T19:11:14Z
transitions:
  - to: ready
    at: 2026-10-23T17:35:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T17:35:59Z
    by: agent-S-0429
  - to: review
    at: 2026-10-23T18:29:16Z
    by: agent-S-0429
  - to: in-progress
    at: 2026-10-23T18:40:32Z
    by: flai
  - to: review
    at: 2026-10-23T18:59:24Z
    by: agent-S-0429
  - to: done
    at: 2026-10-23T19:11:14Z
    by: agent-S-0429
tags: [cli]
topics: [cli, planning, analysis]
touches: [design/digest, design/digest/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/internal/digest/answer.go, flai/internal/digest/answer_test.go, flai/internal/digest/harvest.go, flai/internal/digest/harvest_test.go, flai/internal/digest/seed.go, flai/internal/digest/seed_test.go, flai/internal/digest/sources.go, flai/internal/digest/sources_test.go, docs/operators/settings.md, design/issues/I-0194-the-mcp-task-done-checks-a-story-with-the-installed-flai-so-a-story-that-widens-a-check-s-format-fails-its-own-task-close.md, design/issues/summary.md, flai/cmd/stream_sync_test.go, design/issues/I-0200-testsynclistspathschangedoutsidetheclaim-fails-on-main-since-s-0468-its-fixture-gives-a-task-the-folder-touch-design-which-adr-0188-now-refuses.md]
after: [S-0428]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 4375
  priced_by: harness
  turns:
    - day: 2026-10-23
      ceremony: 1
      test_runs: 1
      hand_edits: 5
      work: 84
  models:
    - model: claude-opus-5-5
      input: 348
      output: 158365
      cache_read: 25960328
      cache_write: 840717
      cost: 13.4134
  strategic:
    - kind: orchestrator
      seconds: 2167
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 151
          output: 2362
          cache_read: 27874908
          cache_write: 54030
          cost: 6.143
        - model: claude-sonnet-5-5
          input: 22
          output: 510
          cache_read: 348738
          cache_write: 69532
          cost: 0.2885
cost_of_delay:
  value: 76.92
  by: planner-S-0429
  at: 2026-10-23T16:44:02Z
forecast:
  duration: 50m
  delivery: 2026-10-24T06:15:00Z
  basis: "flai's median of 107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band gives 36m for size 20 (5 criteria, 15 touches); raised to 50m because the seed run over 610 threads and 168 closed issues, and fixing what flai check reports on the entries it writes, is work the size does not count; delivery is flai's play-out, 21st in the pull order behind S-0428 and the E-0023 and E-0024 stories, moved by the raise."
  by: planner-S-0429
  at: 2026-10-23T16:42:43Z
finalized:
  by: orchestrator
  at: 2026-10-23T16:44:53Z
---
# S-0429 flai digest harvest lists what the archive and the threads hold since a date, and a one-off seed fills the digest from the operator's answers and the closed issues so there is something to prime with

## Goal

The digest starts full enough to prime with. `flai digest harvest` lists, as candidates with their sources, topics, paths, and dates, what was learned since a date: each bullet of an archived narrative's `## Decisions`, each answer the operator gave on a thread, each recommendation they confirmed or overruled, and each closed issue's cause and fix. It is deterministic and writes nothing, so the curator (S-0431) reads the same list on every run. `flai digest seed` writes an entry for each operator answer and each closed issue the harvest lists, skipping those an entry already cites; it is run once on this repository in this story, and the entries it writes are committed.

## Acceptance criteria

- [x] `flai digest harvest [--since <date>] [--kind decision|answer|recommendation|cause] [--json]` lists each candidate with its kind, its source (`S-nnnn` and the narrative, `TH-nnnn` and the entry, or `I-nnnn`), the topics and paths of the item it comes from (a thread's anchored item, an issue's story), its date, its text, and the entries that already cite the source; with no `--since` it lists everything.
- [x] An answer is an entry by the operator (TH-0476 question 5) on a thread someone else opened, or a confirmation of a recommendation; an overruled recommendation is an operator's answer that follows a pending recommendation without confirming it, and its candidate carries both texts. Tests cover each, a thread resolved with no answer, and a thread whose only replies are agents'.
- [x] `flai digest seed [--since <date>] [--dry-run] [--json]` writes one `answer` or `recommendation` entry per such candidate, titled from the thread's question and the answer's first sentence, and one `cause` entry per closed issue that names a cause and a fix, through the same path `flai digest new` takes; a source an entry already cites is skipped, so a second run writes nothing.
- [x] The seed is run once on this repository; the story's notes record how many candidates of each kind the harvest listed, how many entries the seed wrote, and how many it skipped, and `flai check --strict` reports nothing on the entries.
- [x] `docs/users/flai.md` with `flai-reference.md` regenerated and `design/system/flai-cli.md` describe both commands; the digest folder's README says the seed's entries are a starting point the curator merges.

## Tasks
- T-1813 One rule in flai/internal/digest turns a thread into its operator's answers and overruled recommendations, each a candidate entry with its topics and paths
- T-1814 flai/internal/digest reads the decisions of archived narratives and the causes and fixes of closed issues as candidate entries
- T-1815 digest.Harvest lists every candidate since a date with the entries that already cite its source, and digest.Seed writes the answers and causes not yet cited
- T-1816 flai digest harvest and flai digest seed print the candidates and write the seed, as text and JSON
- T-1817 The seed is run once on this repository, its entries committed and counted, and the docs describe harvest and seed
- T-2089 TestSyncListsPathsChangedOutsideTheClaim's fixture names the file its task touches, not the folder ADR-0188 refuses, so the merge queue's integration tier passes

## Notes

- Drafted by the planner for E-0025, after S-0428, whose format and writer it uses.
- Decisions from narratives are listed but not seeded. There are hundreds (381 archived narratives), and most are already recorded in a design document or an ADR, which only a model can judge. The curator's first run takes them in (S-0431).
- The answer-to-entry rule here is the one S-0433 runs at the moment the operator answers. Keep it in one function in `flai/internal/digest` so both use it.

### Seed run

Run once on this repository on 2026-10-23 in the story's worktree, with `scripts/flai.sh digest harvest --json` then `scripts/flai.sh digest seed --json`:

| Kind | Harvest listed | Seed wrote | Seed skipped |
|------|----------------|------------|--------------|
| `decision` | 2629 | not seeded | not seeded |
| `answer` | 198 | 198 | 0 |
| `recommendation` | 12 | 12 | 0 |
| `cause` | 172 | 172 | 0 |

- The 382 entries are D-0001 to D-0382. No write was refused, so no entry had a secret's shape.
- A second `flai digest seed` wrote nothing: it skipped all 382 candidates as cited.
- `flai check --strict` reports nothing on the entries, and the markdown lint passes on `design/digest`.
- Many answers are terse (`yes`, `A`, `done`, `take recommendation`): each entry quotes the question beside its answer, and the curator (S-0431) merges and retitles them.

### Planning

Revisited by planner-S-0429 on 2026-10-23, with S-0428 in progress (its digest package not yet landed) and E-0025's cost of delay input set.

- **Touches.** The 15 declared touches stand, file by file; `flai touches suggest` added none.
  - **Declared, from layout:** `answer.go`, `sources.go`, `harvest.go`, and `seed.go` in `flai/internal/digest`, each with its test, one pair per task of the first three layers; `flai/cmd/digest.go` and its test, which S-0428 creates and T-1816 extends.
  - **Declared, from design:** `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/digest/README.md`, named by criterion 5.
  - **Co-change, not added:** the list is led by `docs/operators/index.md` and `settings.md` (23%, 22%), `design/system/flaiover-dashboard.md` (20%), `docs/users/flaiover.md` (18%), and `design/adrs/README.md` (16%), which travel with any `flai/cmd` change; this story adds no setting, no dashboard change, and no ADR (S-0428 writes ADR-0191). `flai/cmd/root.go` (6%) is not needed: S-0428 registers `digest` there, and the two subcommands are added inside `digest.go`. `flai/internal/threads`, `workitem`, and `issues` are read, not changed: `threads.List`, `Thread.Entries`, `workitem.NarrativeSection`, `workitem.ReadNarrative`, `Item.CompletedAt`, and `issues.List` give what the sources need, so the contingency on `narrative.go` in T-1814 falls away.
  - **Folder touch kept:** `design/digest`, on the story and on T-1817, because the seed's entry file names (`D-nnnn-<slug>.md`) are known only once it runs. Only S-0428 writes there before this story, and S-0433 after it, so the hold reaches no ready story while S-0429 is in progress.
- **Forecast.** flai gives 36m (median 107 s per unit of size over 85 done feature stories on claude-opus-5-5 in the large band, size 20). Set to 50m: the seed runs over 610 threads (591 resolved) and 168 closed issues and may write a few hundred entries, and fixing what `flai check --strict` reports on them with `flai digest edit` is work the size does not count. Delivery is flai's play-out of 2026-10-24T06:01Z moved by the raise; flai replays it as the order ahead changes.
- **Cost of delay.** 76.92 USD a week, from `flai cod`: the story's share by forecast (50m of 6h30m) of E-0025's 600 USD a week, which the operator set on TH-0476 as 4h lost per 168h cycle at 150 USD an hour. No input of the story's own; none is needed while it shares the epic's.
- **Sizing facts behind the tasks' notes:** 610 threads under `wip/threads` (none archived; 15 open), 442 archived narratives, 190 issues of which 168 are closed. Issue files carry no cause or fix heading: the cause is the `## Description`, the fix the reason of the `Closed <ts>: <reason>` line, with `## Remediation` as fallback.
- 2026-10-23T18:40:32Z: moved to in-progress: merge queue batch 20261023-183339 fails with S-0429 at 5693210c55b6: integration failed: flai/cmd/stream_sync_test.go:263 TestSyncListsPathsChangedOutsideTheClaim: {"ts":"2026-10-23T18:40:01.477Z","level":"FATAL","msg":"command failed","component":"cmd","err":"T-0001's touch design is a folder holding 1 file: a task takes no folder touch (ADR-0188); name the files the task changes, or leave the folder to the story's touches with a ### Planning note"}
