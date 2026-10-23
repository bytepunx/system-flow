---
id: S-0429
type: story
nature: feature
title: flai digest harvest lists what the archive and the threads hold since a date, and a one-off seed fills the digest from the operator's answers and the closed issues so there is something to prime with
status: in-progress
parent: E-0025
owner: alex
created: 2026-10-16T05:58:28Z
updated: 2026-10-23T17:35:59Z
transitions:
  - to: ready
    at: 2026-10-23T17:35:36Z
    by: orchestrator
  - to: in-progress
    at: 2026-10-23T17:35:59Z
    by: agent-S-0429
tags: [cli]
topics: [cli, planning, analysis]
touches: [design/digest, design/digest/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/internal/digest/answer.go, flai/internal/digest/answer_test.go, flai/internal/digest/harvest.go, flai/internal/digest/harvest_test.go, flai/internal/digest/seed.go, flai/internal/digest/seed_test.go, flai/internal/digest/sources.go, flai/internal/digest/sources_test.go]
after: [S-0428]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 792
  estimated: true
  priced_by: estimate
  turns:
    - day: 2026-10-23
      test_runs: 1
      hand_edits: 1
      work: 25
  models:
    - model: claude-opus-5-5
      input: 144
      output: 713
      cache_read: 6259371
      cache_write: 421425
      cost: 2.254
  strategic:
    - kind: orchestrator
      seconds: 483
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 54
          output: 793
          cache_read: 5075714
          cache_write: 22828
          cost: 1.1215
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

- [ ] `flai digest harvest [--since <date>] [--kind decision|answer|recommendation|cause] [--json]` lists each candidate with its kind, its source (`S-nnnn` and the narrative, `TH-nnnn` and the entry, or `I-nnnn`), the topics and paths of the item it comes from (a thread's anchored item, an issue's story), its date, its text, and the entries that already cite the source; with no `--since` it lists everything.
- [ ] An answer is an entry by the operator (TH-0476 question 5) on a thread someone else opened, or a confirmation of a recommendation; an overruled recommendation is an operator's answer that follows a pending recommendation without confirming it, and its candidate carries both texts. Tests cover each, a thread resolved with no answer, and a thread whose only replies are agents'.
- [ ] `flai digest seed [--since <date>] [--dry-run] [--json]` writes one `answer` or `recommendation` entry per such candidate, titled from the thread's question and the answer's first sentence, and one `cause` entry per closed issue that names a cause and a fix, through the same path `flai digest new` takes; a source an entry already cites is skipped, so a second run writes nothing.
- [ ] The seed is run once on this repository; the story's notes record how many candidates of each kind the harvest listed, how many entries the seed wrote, and how many it skipped, and `flai check --strict` reports nothing on the entries.
- [ ] `docs/users/flai.md` with `flai-reference.md` regenerated and `design/system/flai-cli.md` describe both commands; the digest folder's README says the seed's entries are a starting point the curator merges.

## Tasks
- T-1813 One rule in flai/internal/digest turns a thread into its operator's answers and overruled recommendations, each a candidate entry with its topics and paths
- T-1814 flai/internal/digest reads the decisions of archived narratives and the causes and fixes of closed issues as candidate entries
- T-1815 digest.Harvest lists every candidate since a date with the entries that already cite its source, and digest.Seed writes the answers and causes not yet cited
- T-1816 flai digest harvest and flai digest seed print the candidates and write the seed, as text and JSON
- T-1817 The seed is run once on this repository, its entries committed and counted, and the docs describe harvest and seed

## Notes

- Drafted by the planner for E-0025, after S-0428, whose format and writer it uses.
- Decisions from narratives are listed but not seeded. There are hundreds (381 archived narratives), and most are already recorded in a design document or an ADR, which only a model can judge. The curator's first run takes them in (S-0431).
- The answer-to-entry rule here is the one S-0433 runs at the moment the operator answers. Keep it in one function in `flai/internal/digest` so both use it.

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
