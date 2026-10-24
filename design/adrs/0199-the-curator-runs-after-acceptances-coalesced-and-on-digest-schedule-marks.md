---
id: ADR-0199
title: "The curator runs after acceptances, coalesced, and on digest.schedule, marks absorbed the entries a document now records, and compacts the active entries no work item, narrative, or thread cited within digest.compact_after, deleting nothing"
status: accepted
date: 2026-10-24
supersedes: []
superseded_by: []
refines: [ADR-0196]
---

# ADR-0199 The curator runs after acceptances, coalesced, and on digest.schedule, marks absorbed the entries a document now records, and compacts the active entries no work item, narrative, or thread cited within digest.compact_after, deleting nothing

## Context

[ADR-0196](0196-the-curator-runs-behind-the-curate-host-action-and-takes-in-what-the-archive.md) gave the curator its run: on the operator's word only, it takes in what the archive and the threads gained as digest entries. [ADR-0191](0191-the-digest-is-one-markdown-file-per-entry-under-the-design-folder-s-digest.md) gave each entry a `status` of `active`, `absorbed`, or `compacted`, and a `cited` list that compaction reads, and said nothing in the digest is ever deleted. [ADR-0194](0194-a-story-s-pack-loads-the-digest-entries-its-topics-or-claim-select-ranked.md) primes an `active` entry whole, a `compacted` one as its title, and an `absorbed` one never. Nothing yet runs the curator without the operator asking, moves an entry out of `active`, or fills `cited`. So the digest only grows, and every fact a design document has since taken in is primed twice. TH-0476 settled a 30-day compaction window.

## Decision

**The curator runs on its own after acceptances and on `digest.schedule`, behind the `curate` host action; each run marks absorbed the entries a document now records, with `flai digest absorb`, and compacts, with `flai digest compact`, each active entry no work item, narrative, or thread has cited within `digest.compact_after`. Neither deletes anything.** This refines ADR-0196; the rest of it stands.

1. **After acceptances.** While `curate` is on for a project, `flai serve` starts the curator when it sees a story accepted, that is, moved to `done`, with the trigger `accepted S-nnnn`. An acceptance seen while a curator runs is remembered, not refused: when that run ends, one more run starts, its trigger naming every story accepted meanwhile, `accepted S-nnnn, S-nnnn`, however many there were. Acceptances seen at one look start one run. An acceptance made while `curate` is off starts nothing then or later.
2. **On a schedule.** `digest.schedule` is read as `analysis.schedule` is ([ADR-0099](0099-the-analyzer-runs-behind-the-analyze-host-action-and-writes-one-report-under.md), S-0223): a five-field cron expression in UTC or `daily`, empty meaning none. Each time it comes round with `curate` on, the curator starts with the trigger `schedule <spec>`, such as `schedule daily`. A start the checks refuse, as while a curator runs, is logged and not tried again until the schedule next comes round; times missed between two looks come round once; times that passed while `flai serve` was down, or while `curate` was off, do not come round. When acceptances and the schedule come round at one look, one run starts, with the acceptances' trigger.
3. **Each start is journalled** with its trigger, as the operator's (`asked`) is, and the trigger is the run's in `serve/agents.json` and its activity entry in `wip/agents/curator.md`. Every run takes in since the curator's last activity, as a run with no `--since` does.
4. **A citation** is an entry's ID, `D-nnnn` in any zero padding, written in a work item's file (an epic, story, or task, in `wip/kanban` or the archive), a narrative (in `wip/agents` or the archive), or a thread. The citer is the work item, the narrative's story, or the thread; its date is the citing file's `updated` front matter. An entry's latest citation is the latest such date. Digest entries citing one another do not count. `flai digest cited [--json]` lists each entry with its citers and its latest citation, and writes each entry's `cited` as its citers, in ID order, where that changed.
5. **Absorption.** `flai digest absorb <D-nnnn> --by <path> [--heading <heading>]` sets `status: absorbed` and `absorbed_by` to the path, from the project root, of the design document, convention, or ADR that now records the fact, or to `<path>#<anchor>` with a heading, which must be one of the document's. It refuses a path that names no file and a heading the document does not have, and writes nothing then. The body is kept. An entry absorbed again names the newer document. In every run the curator looks for absorption: for each active entry it searches the design with `doc_search` on the entry's title, reads what it finds, and absorbs the entry only when a document states its fact.
6. **Compaction.** `digest.compact_after` is a Go duration, `720h` (30 days) when unset, and `0` turns compaction off. `flai digest compact [--dry-run] [--json]` runs `cited` first, then turns each `active` entry whose latest citation, or its own `date` when it has none, is older than the window into `compacted`, its body emptied below its `# D-nnnn <title>` heading, and lists what it compacted. `--dry-run` lists what it would compact and writes nothing. An `absorbed` or `compacted` entry is left as it is, so a second run compacts nothing more. The curator runs `flai digest compact` once in every run, after its absorptions.
7. **Nothing is deleted.** Absorption and compaction change an entry's front matter and, for compaction, its body; no file is removed, and git keeps every earlier text. A compacted entry stays compacted when it is cited later; it still primes as its title, so the question it answered is known to have been asked.
8. **The guard** passes `flai digest absorb`, `cited`, and `compact` to the curator, beside `new`, `edit`, `harvest`, `list`, and `show`.

## Consequences

- The digest stays small without the operator asking: the facts the design took in leave the packs, and the facts no one relied on for a month shrink to a line.
- An acceptance with `curate` on costs a curator run; a burst of acceptances costs two, the one under way and one after it.
- Citing an entry's ID in a narrative's `## Decisions` or on a thread is what keeps it whole, as the conventions already tell agents (ADR-0194).
- A citing file's `updated` moves with any later edit of it, so a citation can read newer than it is. That errs toward keeping an entry whole, never toward compacting one that was relied on.
- Absorption is the curator's judgement, made with the design in front of it; `flai digest absorb` only checks that the document, and the heading, exist.
- `flai digest cited` and `compact` write in the checkout they run in, the main checkout for the curator, uncommitted, as its other entries are.

## Alternatives considered

- **A run per acceptance, queued.** A burst of acceptances would start as many curators, each reading what the one before already took in; one run after the burst reads it all.
- **Dating a citation by the thread entry or log line that names it.** More exact, but each kind of file would need its own parser; `updated` is on every one of them, and its error only keeps an entry longer.
- **Reading citations from git history.** Slower on every run, and a citation removed from a file would still count.
- **Deleting compacted entries after a second window.** Loses the record that the question was asked and its `asked_again` count (ADR-0191).
- **Reactivating a compacted entry when it is cited again.** A citation of a one-line entry is a citation of its title; the operator, or the curator with `flai digest edit`, can restore its text from git when it matters.
