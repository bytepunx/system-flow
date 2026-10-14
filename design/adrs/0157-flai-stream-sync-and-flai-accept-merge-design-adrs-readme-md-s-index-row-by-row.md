---
id: ADR-0157
title: "flai stream sync and flai accept merge design/adrs/README.md's index row by row when a rebase stops on it, and the trial merge does not report it"
status: accepted
date: 2026-10-14
supersedes: []
superseded_by: []
refines: [ADR-0098, ADR-0126]
topics: [cli]
---

# ADR-0157 flai stream sync and flai accept merge design/adrs/README.md's index row by row when a rebase stops on it, and the trial merge does not report it

## Context

`flai adr new` adds a row for the new ADR to the index in `design/adrs/README.md`. Most stories add an ADR, and the shared paths include `design/adrs`, so several stories in progress at once each append a row to the same table. Git cannot merge two appends at the same place. [I-0099](../issues/I-0099-two-stories-that-each-add-an-adr-conflict-in-design-adrs-readme-md-whose-index-both-append-a-row-to-and-flai-opens-a-thread-for-each-sync-that-finds-it.md) records the cost:

- `flai stream sync`'s trial merge found every pair of three stories conflicting on this file alone, and opened five threads (TH-0220, TH-0221, TH-0222, TH-0226, TH-0228), resolving and reopening some as later syncs went. Each said to narrow or wait, and neither applied.
- The story accepted second resolved the index by hand at its sync, keeping both rows.

The resolution needs no judgement: every row belongs, and each row is one ADR's. [ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md) settles the issue summary by writing it again, and [ADR-0126](0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md) settles an issue file by merging its instances. The index is not wholly generated: its rows carry hand-written status notes such as "superseded in part by 0016", and the text above the table is prose.

## Decision

**When a rebase that `flai stream sync`, `flai task done`, or `flai accept` runs stops on `design/adrs/README.md`, flai merges the index itself from the stop's three versions, row by row, `git add`s it, and continues the rebase.** This refines ADR-0098 and ADR-0126: the index is one more of the files the rebase merges, as issue files are.

- The text above the table, the table's header, and the text below it each take the side that changed them.
- Each row is keyed by the ADR file it links, and takes the side that changed it. Rows both sides added are both kept, a note one side added to a row is kept, and a row one side removed is removed.
- The rows come out in ADR number order.
- When both sides changed the same part or the same row differently, flai merges nothing, and the stop is left for the agent with the path listed, as before.
- A stop whose every path is the index, a merged issue file, or a generated file is continued.
- The trial merge between two open story branches leaves the index out of a pair's conflicts, as it leaves out the summary and the issue files. A pair whose only conflict is the index counts as clean, so flai sends the two stories no message about it (ADR-0121).

## Consequences

- Two stories that each add an ADR no longer stop each other's sync or acceptance on the index, and no agent reads a message about it.
- Two stories that supersede the same ADR both change its row's status, differently, so the rebase still stops on the index. The trial merge does not report that row, but it still reports the pair, since both also change the superseded ADR's `superseded_by`.
- A hand edit to the text above the table on both sides stops the rebase as before.
- The index stays hand-editable: nothing is regenerated, and `flai check` still warns of a row with no file or a file with no row.

## Alternatives considered

- **Regenerate the index whole from the ADR files' front matter**, as I-0099 recommended and as ADR-0098 does for the summary. The status notes such as "superseded in part by 0016" are written by hand and are not all in front matter, so they would be lost, and the text above the table would have to be held apart from the table.
- **A git merge driver for the file in `.gitattributes`.** Each clone must name the driver in its `git config`, which flai does not install, and `git merge-tree` runs no driver, so the trial merge would still report the file (ADR-0098).
- **Leave it to the agents.** That is what I-0099 records, at about 5m and a thread each time.
