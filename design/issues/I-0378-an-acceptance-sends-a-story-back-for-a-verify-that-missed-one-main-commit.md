---
id: I-0378
title: An acceptance sends a story back for a verify that missed one main commit landed while its close-out ran, though the close-out passed against main
class: efficiency
status: open
count: 1
cost: 15m
first_reported: 2026-11-15T07:05:58Z
last_reported: 2026-11-15T07:05:58Z
updated: 2026-11-15T07:05:58Z
---

# I-0378 An acceptance sends a story back for a verify that missed one main commit landed while its close-out ran, though the close-out passed against main

## Description
An acceptance sends a story back for a verify that missed one main commit landed while its close-out ran, though the close-out passed against main

## Instances

### 2026-11-15T07:05:58Z
Story: S-0738.
S-0738 was sent back on TH-1213 for a verify that did not cover main, closed out again (job S-0738-20261115T065141Z, passed at ee0e8a5 over main 721f679), and was sent back again three minutes after it reached review, on TH-1215, for one path main gained while that close-out ran: the operator's markdown fix to design/digest/D-0627 (f01a053e7). The close-out's last check passes over main commits that change only wip paths (ADR-0135), but a digest file is under design, so it did not sync. Each round costs a sync, a ten-minute close-out, and an acceptance attempt; the markdown tier alone was needed.

## Remediation
