---
id: TH-0476
title: "E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window"
anchor:
  path: wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md
  item: E-0025
status: open
participants: [planner-E-0025]
created: 2026-10-16T05:57:18Z
updated: 2026-10-19T04:51:24Z
---

# TH-0476 E-0025 needs a cost of delay input and four design answers: recommend 4h lost per cycle, one file per entry under design/digest, a curator role of its own, a 30-day compaction window

On wip/kanban/epics/E-0025-a-digest-of-the-archive-kept-by-an-agent-and-decayed-as-the-design-absorbs-it-primes-a-story-s-agent-with-what-past-stories-decided-and-what-the-operator-answered-so-the-same-question-is-not-asked-twice.md.

## Entries

### 2026-10-16T05:57:18Z planner-E-0025
I am drafting E-0025's stories now. Five inputs are yours. Each has my recommendation first. I draft on the recommended answers and change the stories if you answer otherwise.

## 1. Cost of delay input (blocks the stories' values)

`flai cod E-0025` refuses: the epic has no inputs.

- **Recommended:** `time_lost_per_cycle: 4h` on E-0025. At `planning.hour_rate` 150 and the default 168h cycle, that is 600 USD a week, shared among the stories by duration. It counts the operator's time on questions asked again and the rework of a story that guessed against an earlier answer. I-0151's 117h a cycle is mostly agents idling, not work, so I take only a share of it.
- Alternatives: `2h` (300 USD a week), or `8h` (1200 USD a week).

## 2. Where the digest lives and how an entry is written (the epic's first open question)

- **Recommended:** one markdown file per entry, `design/digest/D-nnnn-<slug>.md`, named in the manifest as `design/analysis` is. The title is the fact in one sentence. Front matter holds `id`, `kind` (`decision`, `answer`, `recommendation`, `cause`), `sources` (S-, TH-, I- IDs), `topics`, `paths`, `date`, `status` (`active`, `absorbed`, `compacted`), `absorbed_by`, and `asked_again` (the threads that asked what it already answered). `doc_search` and the dashboard's documents page then find and render each entry as they do an ADR, and prime selects an entry by its topics and paths as it selects an ADR.
- Alternatives: one file per topic, each entry a section (an entry on two topics has to pick one); or front matter on the archived items (scatters the digest across `wip/archive`, which `doc_search` does not index).

## 3. Who keeps the digest (the epic's second open question)

- **Recommended:** a role of its own, the **curator**: role `curate`, host action `curate`, activity document `wip/agents/curator.md`, definition `.claude/agents/curator.md`, settings `digest.agent` and `digest.schedule`. It is built on today's analyzer machinery, as the epic's Notes allow, without waiting for E-0024. It writes entries only through `flai digest` commands, and the guard holds it to them.
- Alternatives: the analyzer keeps it, with a `digest` focus (cheaper, but its activity mixes with analyses and its guard grows a second folder); or the curator waits for E-0024 and is a role file after S-0419 (less code, but the epic then waits behind E-0024's twelve stories).

## 4. The compaction window (the epic's third open question)

- **Recommended:** `digest.compact_after: 30d`. An entry is needed when a narrative, thread, or work item cites its ID. An active entry nobody has cited in 30 days is compacted to its one-line title, which is still primed, as a brief is. Prime records nothing, so the measure stays in the repository.
- Alternatives: `60d`; or count an entry as needed whenever a pack included it (needs a use log under `.flai-cache`, lost on another host).

## 5. Who counts as the operator for an answer that becomes an entry at once

- **Recommended:** the manifest's `owner` (alex), as `permission_prompt` takes it (ADR-0097). Their reply to a thread they did not open, or their confirmation of a recommendation, becomes an `answer` entry. A reply that overrules a pending recommendation becomes a `recommendation` entry that says what was overruled and why.
- Alternative: any author that is not an agent run (`planner-*`, `orchestrator`, `analyzer`, `agent-S-*`, `flai`).

Meanwhile I plan everything that needs no answer. Once question 1 is answered, I write the stories' values.

### 2026-10-19T04:51:24Z planner-E-0025
An update. E-0025's eight drafts, S-0428 to S-0435, and their 42 tasks, T-1807 to T-1848, are complete except for their cost-of-delay values, which wait on question 1. The drafts already follow my recommended answers to questions 2 to 5. The plan, its order, and its layers are on TH-0480.

Since I opened this thread:

- **S-0426 was accepted.** Its finding, in design/system/agent-memory.md, supports a digest primed by topics and paths.
- **E-0021 and E-0026 are done.** Nothing in E-0025 waits on them.
- **No E-0024 role story has started.** So the curator still builds on the analyzer machinery (question 3).

The quickest answer is "yes to all five". I would then set `time_lost_per_cycle: 4h` on E-0025, write each story's value from `flai cod`, and leave the drafts for you to finalize. If any answer differs, name only that one, and I will change the stories it reaches.
