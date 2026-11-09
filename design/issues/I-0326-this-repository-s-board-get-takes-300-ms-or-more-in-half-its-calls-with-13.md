---
id: I-0326
title: This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review
class: efficiency
status: open
count: 1
cost: 5m
first_reported: 2026-11-09T23:09:35Z
last_reported: 2026-11-09T23:09:35Z
updated: 2026-11-09T23:09:35Z
---

# I-0326 This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review

## Description
This repository's board.get takes 300 ms or more in half its calls with 13 boards at once, from the merge queue's git while stories are in review

## Instances

### 2026-11-09T23:09:35Z
Story: S-0636.
S-0636 measured 13 board.get calls at once on the operator's host on 2026-11-09 at 23:04 to 23:05 UTC, load 96 to 123 on 24 cores, after S-0602, S-0604, and S-0606 (main at f667efc2d). The other 12 projects answered under 300 ms in all 120 calls. This repository took 169 to 404 ms, and 5 of 10 calls took 300 ms or more (373, 318, 304, 350, 404 ms), against 244 to 657 ms before (8b29cf861). With 5 stories in review, merge.queue took 91 to 241 ms of each call: mergequeue.ComputeAt (flai/internal/mergequeue/eligible.go) runs mainDirt's git status --porcelain on the main checkout, 52 to 148 ms, and five git rev-parse calls one after another, 21 to 81 ms, on every board.get. The thread read S-0602 kept no longer shows. Candidate fix: keep the queue's git reads between calls while the main-refs watcher saw no move and the item snapshot is fresh, as S-0604 did for release.pending's show-ref, or run the rev-parse calls at once. See design/system/server-performance.md § 13 boards at once, measured again after S-0602 and S-0606 (S-0636), TH-0912, D-0535.

## Remediation
