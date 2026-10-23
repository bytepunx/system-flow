---
id: ADR-0194
title: "A story's pack loads the digest entries its topics or claim select, ranked against the story within an eighth of the budget after the briefs, never an absorbed one, and the orchestrator's the answers and recommendations ranked against each thread awaiting the operator"
status: accepted
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0049, ADR-0075]
---

# ADR-0194 A story's pack loads the digest entries its topics or claim select, ranked against the story within an eighth of the budget after the briefs, never an absorbed one, and the orchestrator's the answers and recommendations ranked against each thread awaiting the operator

## Context

[ADR-0191](0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md) gave the digest a home: one entry per fact a past story learned, with the topics and paths it bears on and a status of `active`, `compacted`, or `absorbed`. S-0429 seeded it with the operator's answers and the closed issues' causes, 383 entries on this repository. None of them reaches an agent yet: a story's pack ([ADR-0049](0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md)) and a strategic agent's ([ADR-0075](0075-the-planner-the-orchestrator-and-the-analyzer-prime-by-role-plan-orchestrate-or.md)) load conventions, design, tech, and ADRs only, so an agent asks the operator again what a past story already settled.

Most entries carry broad topics such as `flai`, `cli`, `go`, and `all`, so topics alone would select most of the digest for any story of flai. Selection has to be wide and loading narrow: ranked against the story, and capped.

## Decision

**A story's pack gains a digest step, between the briefs and the ranked sections: the digest entries its topics or its claim select, ranked by BM25 against the story's title, goal, and criteria, loaded best first within the step's share of the budget, an eighth of it, and within the budget itself when the pack is within it before the step. Each is marked as the digest's. An absorbed entry never loads. The orchestrator's pack gains the same step over the `answer` and `recommendation` entries, ranked against each thread awaiting the operator.** This refines ADR-0049's pack and ADR-0075's strategic packs; the rest of both stands.

### Which entries a story's pack selects

- An entry is selected when its status is `active` or `compacted`, and one of its topics other than `all` is one of the story's topics, or one of its paths meets the story's claim: the story's touches and its tasks', each path equal to the entry's, a folder holding it, or held by it.
- An `absorbed` entry is never selected: the document its `absorbed_by` names says it now, and the pack reaches that document by its own steps.
- The planner's and the analyzer's packs take no digest step; they find entries with `doc_search`, which indexes the design folder.

### Where entries rank, and how much they take

- The step runs after the briefs and before the ranked sections, so the digest is never crowded out by sections and never crowds them out of more than its share.
- Entries are ranked by BM25 against the query the ranked step uses, over each entry's title and body. An entry whose paths meet the claim ranks before one selected by topics alone; within each, by score, then the newest date first. An entry with no word in common with the query is not loaded unless its paths meet the claim.
- Entries load in that order while the digest items and their heading together stay within an eighth of the budget, 10 KB of 80 KB, and, on a pack within its budget before the step, the pack within the budget; one that does not fit is skipped and the next is tried. The share is a constant beside `DefaultBudget`, not a manifest key.
- When the pack is over its budget before the step, as `Exceeded` says, the step still loads its entries within its share, as every brief is kept on a pack over its budget (TH-0032); the ranked step loads nothing then. A pack over its budget grows by at most the share, and a pack within its budget never goes over it because of the step. The operator chose this on TH-0624 (2026-10-23), since every pack on this repository was over its budget before the step, from its conventions at 80 KB and from its briefs at 240 KB, and none would otherwise have loaded an entry.

### How an entry loads

- Each is an item with `step: digest`, its ID `D-nnnn`, its path, its title, its size, and its reason: what selected it, `topics: <the topics it shares>` or `path: <the claimed path it meets>`, then `rank n`, such as `path: flai/cmd/prime.go; rank 1`, in text and in `--json`.
- An `active` entry loads its title and body. A `compacted` entry loads as a brief: its one-line title, ID, path, and reason, with no body, as ADR-0191's compaction leaves none.
- In text the entries print together under a `digest` heading, after the briefs, saying that they are what past stories decided and the operator answered, not the design: an agent that relies on one cites its `D-nnnn` in its narrative or thread, checks it against the design before changing what it describes, and asks the operator only what no entry answers.
- Entries go into parts ([ADR-0104](0104-the-mcp-prime-tool-returns-a-context-pack-in-parts-of-at-most-40-000-bytes-of.md)) as any other item does.

### A sub-agent's pack

A sub-agent's pack ([ADR-0059](0059-a-story-s-agent-hands-search-test-runs-and-verification-to-an-explorer-and-a.md)) selects and ranks entries as the story's does, but loads each, active or compacted, as a brief only, after the design and ADR briefs, while the budget has room; those with no room are counted with the briefs left out.

### The orchestrator's pack

`flai prime --role orchestrate` selects the `active` and `compacted` entries of kind `answer` or `recommendation`, and ranks them, by BM25 over each entry's title and body, against each thread awaiting the operator (ADR-0181), by its title and first entry. It takes, round by round, each thread's best entry not yet taken, threads in ID order, within the step's share, an eighth of the budget, and the budget as a story's pack does; an entry with no word in common with a thread is not taken for it. The reason names the thread, `thread TH-nnnn; rank n`. With no thread awaiting the operator the step loads nothing. An entry's path is a valid `source` on `thread_reply`, so a recommendation cites the answer the operator gave last time.

## Consequences

- What a past story decided reaches the next story that touches the same paths or shares words with its goal, without anyone writing it into a design document first.
- Every pack spends up to an eighth of its budget on the digest, taken from what the ranked sections had, or added to a pack already over its budget. The conventions, what is named, and the briefs are never cut for it.
- Citing an entry's ID is what keeps it from being compacted (S-0432), so the conventions tell agents to cite what they rely on.
- A compacted entry still reaches a story as one line, so the question it answered is known to have been asked.

## Alternatives considered

- **Entries as candidates of the ranked step, one ranking with the design sections.** Entries are short and rank high on the words of a story's goal; when there are many they would take the whole ranked share, and the design sections the ranked step exists for would not load. A share of their own keeps both.
- **Select by topics alone.** The seed's topics are broad: nearly every flai story would select most of the digest, in no useful order. The claim and the rank are what make an entry relevant.
- **A manifest key for the share.** No project has asked for one; a constant is changed with a story when one does.
- **Rank the orchestrator's entries against all its threads at once.** A thread asked again shares its words with one answer; pooled with the other threads' words, that answer can fall below entries that match many threads weakly. Ranking per thread puts it first for the thread that asks again.
