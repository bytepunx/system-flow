---
id: S-0429
type: story
nature: feature
title: flai digest harvest lists what the archive and the threads hold since a date, and a one-off seed fills the digest from the operator's answers and the closed issues so there is something to prime with
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T05:58:28Z
updated: 2026-10-19T20:33:45Z
transitions: []
tags: [cli]
topics: [cli, planning, analysis]
touches: [design/digest, design/digest/README.md, design/system/flai-cli.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/digest.go, flai/cmd/digest_test.go, flai/internal/digest/answer.go, flai/internal/digest/answer_test.go, flai/internal/digest/harvest.go, flai/internal/digest/harvest_test.go, flai/internal/digest/seed.go, flai/internal/digest/seed_test.go, flai/internal/digest/sources.go, flai/internal/digest/sources_test.go]
after: [S-0428]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
forecast:
  duration: 45m
  delivery: 2026-10-20T09:33:00Z
  basis: "Its own forecast of 45m; 28th in the pull order with an in-progress limit of 5, behind S-0404, S-0405, S-0403, S-0402, S-0401, S-0400, S-0241, S-0406, S-0407, S-0408, S-0409, S-0410, S-0411, S-0412, S-0414, S-0415, S-0416, S-0417, S-0418, S-0419, S-0420, S-0421, S-0422, S-0423, S-0424, S-0425 and S-0428."
  by: flai
  at: 2026-10-19T20:33:45Z
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

- **Touches.** The story declared none. Every touch comes from its five tasks, file by file.
  - **From layout:** `answer.go`, `sources.go`, `harvest.go`, and `seed.go` in `flai/internal/digest`, each with its test; `flai/cmd/digest.go` and its test, which S-0428 makes.
  - **From design:** `design/system/flai-cli.md`, `docs/users/flai.md`, `docs/users/flai-reference.md`, and `design/digest/README.md` for criterion 5.
  - **Co-change, added:** `design/system/flai-cli.md` (36% of the commits that changed `threads.go`, `narrative.go`, or `cmd/thread.go`) and `docs/users/flai.md` (23%).
  - **Co-change, not added:** `flai/internal/threads/threads_test.go` and `flai/internal/workitem/narrative_test.go` (27%), since the digest reads threads and narratives without changing their packages; T-1814 adds `narrative.go` if reading an archived narrative needs it, and `flai task done` widens the touches then. `design/system/agent-narrative.md` (23%), since no narrative rule changes here.
  - **Folder touch kept:** `design/digest`, because the seed's entry file names are known only once it runs. Only this story writes there while it is in progress, and it holds a ready story only for that time.
- **Forecast.** Raised from flai's 32m to 45m. The seed runs over about 480 threads and 138 closed issues, and fixing what the check reports on its entries is work the size, 5 criteria and 15 touches, does not count. Delivery is flai's play-out moved by the raise.
- **Cost of delay.** Not set. Neither the story nor E-0025 has inputs, so `flai cod` refuses. The epic's input is asked on TH-0476, question 1, and the value follows its answer.
