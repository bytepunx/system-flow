---
title: Simplifying agent-produced prose
updated: 2026-10-16
status: active
topics: [cli, template]
---

# Simplifying agent-produced prose

The finding of S-0413. The operator asked for research on making the prose flai's agents write simpler and shorter, so that the operator and other agents read less and read it faster. The operator named two means: a natural language linter, such as Vale or proselint, and a small model hosted locally and fitted with LoRA or a similar method, which rewrites verbose agent text until it passes that linter. This document says what prose flai's agents write today and what checks it, then covers prose linters, a local rewriting model, and the alternatives, and ends with a recommendation. Nothing here is built.

Every figure about this repository comes from a command run on 2026-10-16. Design documents and ADRs are linked; code is named by its path.

## Today

Measured on 2026-10-16 over the committed tree at `53eb8e40`, with a Python script kept outside the repository. A word is a run of letters and digits outside front matter, fenced code, and HTML comments. An entry is a `###` section under a narrative's `## Log` or under a thread's or message's `## Entries`. Live and archived files are counted together.

| Kind | Path | Reader | Files | Words | Median words per file | Entries | Median words per entry |
|------|------|--------|-------|-------|-----------------------|---------|------------------------|
| Story narratives | `wip/agents/S-*.md`, `wip/archive/agents/` | The agent that resumes the story; the operator | 381 | 236,110 | 556 | 1,989 log entries | 20 |
| Strategic agents' activity logs | `wip/agents/planner.md`, `orchestrator.md`, `analyzer.md` | The operator; that agent's next session | 3 | 35,031 | 4,236 | 645 | 49 |
| Threads | `wip/threads/` | The operator and agents | 468 | 173,488 | 317 | 1,012 agent entries | 105 |
| Messages between agents | `wip/messages/` | The two stories' agents; the operator on escalation | 126 | 47,026 | 308 | 554 agent entries | 64 |
| Epics | `wip/kanban/epics/`, `wip/archive/kanban/epics/` | The operator and agents | 25 | 11,338 | 304 | | |
| Stories | `wip/kanban/stories/`, `wip/archive/kanban/stories/` | The operator and agents | 413 | 219,461 | 524 | | |
| Tasks | `wip/kanban/tasks/`, `wip/archive/kanban/tasks/` | Agents; the operator | 1,719 | 266,613 | 133 | | |
| Issues | `design/issues/I-*.md` | The operator, agents, the analyzer | 146 | 43,509 | 204 | | |
| ADRs | `design/adrs/` | The operator and agents | 162 | 107,983 | 639 | | |
| Design documents | `design/system/` | The operator and agents | 26 | 206,366 | 3,490 | | |
| Analyzer reports | `design/analysis/` | The operator | 1 | 1,983 | 1,983 | | |
| Commit messages | `git log`, commits with a Claude trailer | The operator | 1,955 | 117,132 | 59 | 1,411 with a body | 58 per body |
| Final report to the operator | Not stored in the repository | The operator | | | | | |

The kinds above hold 1,466,040 words. A small part is not agent prose: the operator's thread entries, flai's own notices in threads and messages, the operator's brief, and the goals the operator writes for stories.

- A narrative's largest section is `## Decisions`: 89,314 words across the 381 narratives, against 63,672 in `## Log` and 46,242 in `## Current state`, each counted with its heading.
- Threads hold 1,558 entries. Agents wrote 1,012 of them, 151,764 words. The operator wrote 401, 2,674 words, with a median of 1 word: 161 of them are `Resolved.` alone. flai wrote 145, 9,506 words.
- Messages hold 635 entries; flai wrote 81 of them.
- The repository has 2,954 commits. 1,955 carry a Claude trailer. Of the other 999, 278 carry flaiover's trailer and 721 carry none.
- The design documents include the operator's [brief](brief.md), 865 words, and this folder's README.
- The final report goes to the operator in the harness's session, not to a file. The harness keeps it in its session log, which flai reads for usage only (`flai/internal/usage`).

### Readability

Measured over the same tree with the same script. Sentences are split at `.`, `!`, or `?` followed by a space and a capital or a digit, within paragraph and list lines. Headings, tables, and code are left out, and an inline code span counts as one word. The Flesch-Kincaid grade uses a heuristic syllable count, so it compares rows rather than standing alone.

| Prose | Sentences | Mean words per sentence | Median words per sentence | Sentences over 25 words | Flesch-Kincaid grade |
|-------|-----------|-------------------------|---------------------------|-------------------------|----------------------|
| Narrative log entries | 4,150 | 13.7 | 11 | 13% | 6.8 |
| Narratives, whole | 12,746 | 16.9 | 15 | 20% | 7.5 |
| Strategic agents' log entries | 3,363 | 9.8 | 5 | 10% | 6.2 |
| Thread entries by agents | 10,277 | 12.9 | 12 | 6% | 5.5 |
| Message entries by agents | 1,857 | 13.7 | 12 | 12% | 5.2 |
| Epics | 484 | 22.4 | 18 | 33% | 10.5 |
| Stories | 10,631 | 18.0 | 16 | 20% | 8.1 |
| Tasks | 14,200 | 16.5 | 14 | 16% | 6.6 |
| Issues | 2,178 | 18.5 | 15 | 26% | 8.5 |
| ADRs | 4,868 | 20.7 | 19 | 28% | 8.9 |
| Design documents | 4,573 | 22.2 | 19 | 32% | 9.2 |
| Analyzer reports | 69 | 21.3 | 21 | 33% | 8.0 |
| Commit messages | 11,013 | 10.6 | 11 | 1% | 5.5 |

Sentences are shortest in commits, threads, messages, and logs. They are longest in the documents agents write for later reading: design documents, epics, ADRs, the analyzer's report, and issues. There a quarter to a third of sentences run past 25 words.

### What checks this prose today

| Check | What it checks | Enforced |
|-------|----------------|----------|
| [Communication convention](../conventions/communication.md) | Reports to the operator: lead with the outcome, one idea per sentence, short sentences with verbs, no restatement of the request, no praise | No. Agents read it in every prime pack. |
| [Documentation convention](../conventions/documentation.md) | Documents: write plainly, tables for parallel facts, short sentences, no filler, no marketing, no duplicated content | No. Agents read it in every prime pack. |
| flai's markdown lint, [ADR-0061](../adrs/0061-flai-lints-the-markdown-it-writes-in-wip-with-its-own-implementation-of-the.md) | 23 markdownlint rules on the work items, threads, and narratives flai writes in `wip/`: headings, lists, whitespace, bare URLs, code fences, emphasis, tables | Yes. Package `flai/internal/mdlint`: `flai check` reports a finding, and a command that writes refuses one. |
| The `markdown` test tier in [system-flow.yaml](../../system-flow.yaml) | Every markdown file, through `scripts/lint-md.sh`, which runs markdownlint-cli2 0.20.0 with `.markdownlint.yaml` | Yes, in `flai test` and the close-out. Line length (MD013) is off. |
| Title limits | A conversation's title is cut to 72 runes (`flai/internal/messages`); the host API refuses an epic or story title over 200 characters (`flai/internal/hostapi`) | Yes, on titles only. |
| The prime pack's size budget, [ADR-0049](../adrs/0049-a-story-s-context-pack-fits-a-size-budget-what-the-story-names-loads-whole-what.md) | What an agent reads at story start: 81,920 bytes by default (`flai/internal/context`) | Yes, on what is read, not what is written. |

Nothing checks the prose itself. On 2026-10-16 a search of `flai/`, `flaiover/src/`, `template/`, `scripts/`, the manifest, and the repository's configuration found no check of sentence length, wordiness, passive voice, or readability. It found no prose linter such as Vale, proselint, write-good, or textlint. The two conventions state the style; agents follow them or not, and no tool measures the result.

The prime budget is the one place flai limits what an agent reads, and conventions alone exceed it. `flai prime --story S-0413 --json`, run on 2026-10-16, reports a pack of 128,217 bytes in 1,114 lines against the 81,920-byte budget, exceeded by conventions. The 14 conventions it loads whole come to 81,868 bytes, and their README to 3,063 more. The largest are [work-management.md](../conventions/work-management.md), 21,819 bytes, and [delegation.md](../conventions/delegation.md), 13,167 bytes.
