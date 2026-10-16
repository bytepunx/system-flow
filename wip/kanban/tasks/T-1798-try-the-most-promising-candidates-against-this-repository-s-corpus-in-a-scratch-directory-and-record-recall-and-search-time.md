---
id: T-1798
type: task
nature: research
title: Try the most promising candidates against this repository's corpus in a scratch directory and record recall and search time
status: backlog
parent: S-0426
owner: alex
created: 2026-10-16T05:58:17Z
updated: 2026-10-16T05:58:17Z
transitions: []
stream: S-0426
tags: []
touches: [design/system/agent-memory.md]
after: [T-1797]
---
# T-1798 Try the most promising candidates against this repository's corpus in a scratch directory and record recall and search time

## Work

Pick at most three approaches from `## Open source candidates` and `## State of the art` that need no paid service and no account: at least one hybrid of BM25 and a local embedding model, and flai's own BM25 as the baseline. Run each in a scratch directory outside the repository, over a copy of `design/` and `docs/`, and never install anything on the host beyond that directory.

Write a query set of at least 20 questions an agent asked of this repository, taken from thread titles and stories' goals, each with the sections that answer it. Measure for each approach: recall at 5 and at 10 against that set, query time, index build time, refresh time after one file changes, and size on disk and in memory.

Add `## Trial` to `design/system/agent-memory.md`: the query set or where it came from, what was run, the results in one table, and what each result means for flai. Mark anything not tried as not tried, with why. If the trial needs a paid service or a host install to go on, stop and ask on a thread on S-0426 rather than proceed.

Waits for T-1797: it chooses what to try from the two sections before it, and writes the next section of the same document.

## Done when

- `## Trial` has the query set's source, the commands run, and one table of recall, query time, build and refresh time, and footprint per approach, with flai's BM25 as the baseline.
- Nothing was installed outside the scratch directory, and nothing in the repository but the document changed.

## Notes
