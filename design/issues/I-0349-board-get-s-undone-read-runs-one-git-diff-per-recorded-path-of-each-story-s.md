---
id: I-0349
title: board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call
class: efficiency
status: open
count: 1
cost: 2m
first_reported: 2026-11-11T08:52:49Z
last_reported: 2026-11-11T08:52:49Z
updated: 2026-11-11T08:52:49Z
---

# I-0349 board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call

## Description
board.get's undone.read runs one git diff per recorded path of each story's undone record, one after another, on every call

## Instances

### 2026-11-11T08:52:49Z
Story: S-0698.
S-0698 measured this repository's board.get on 2026-11-11 at 08:45 UTC (design/system/server-performance.md, The merge queue's git read once (S-0698)). Once the merge queue's git was one for-each-ref beside one git status, undone.read was the largest part of the call: 107 to 126 ms of a 150 to 187 ms call with 13 boards at once, at a load of 6. S-0693, in progress, has an undone record of 20 paths, and boardUndone (flai/internal/hostapi/hostapi.go) runs storygit.RecordedUndone for it on every board.get: one git rev-parse --abbrev-ref HEAD and one git diff per recorded path, one after another, 20 diffs a call. In three more pairs at loads of 6.4 to 10.9, undone.read took 214 ms of a 391 ms call before and 118 of 328 ms after. Candidate fix: one git diff for all the record's paths, or keep the answer while the branch's head and main's are unchanged, as the trial merges are kept (S-0678).

## Remediation
