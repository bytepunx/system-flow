---
id: T-1750
type: task
nature: research
title: Record where flai's agents write prose today, how much, and what checks it
status: done
parent: S-0413
owner: alex
created: 2026-10-16T05:45:55Z
updated: 2026-10-16T06:00:55Z
transitions:
  - to: ready
    at: 2026-10-16T05:46:55Z
    by: agent-S-0413
  - to: in-progress
    at: 2026-10-16T05:46:57Z
    by: agent-S-0413
  - to: done
    at: 2026-10-16T06:00:55Z
    by: agent-S-0413
stream: S-0413
tags: []
touches: [design/system/agent-prose.md, design/system/README.md]
usage:
  source: log
  seconds: 838
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 50
      output: 20919
      cache_read: 2261791
      cache_write: 79411
      cost: 1.3617
---
# T-1750 Record where flai's agents write prose today, how much, and what checks it

## Work

Create `design/system/agent-prose.md`, the finding of S-0413, in the shape of `agent-adapters.md` (S-0339), and add its row to `design/system/README.md`. Write its first section: every kind of prose flai's agents write for the operator or another agent (narratives, thread entries, messages, work items, issues, ADRs, design documents, commit messages, final reports), who reads each, its volume measured from this repository, and what checks that prose today (the communication convention, flai's markdown lint of `wip/`, the markdownlint tier).

Waits for nothing: it creates the document the later tasks add to.

## Done when

- The document exists with front matter, an introduction, and a `## Today` section whose figures come from commands run on this repository.
- `design/system/README.md` lists it.

## Notes
