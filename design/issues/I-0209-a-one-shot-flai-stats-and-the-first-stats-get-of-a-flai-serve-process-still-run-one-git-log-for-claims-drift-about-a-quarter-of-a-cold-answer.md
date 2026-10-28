---
id: I-0209
title: A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer
class: efficiency
status: closed
count: 1
cost: 650ms
first_reported: 2026-10-26T04:42:46Z
last_reported: 2026-10-26T04:42:46Z
updated: 2026-10-28T18:36:41Z
---

# I-0209 A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer

## Description
A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer

## Instances

### 2026-10-26T04:42:46Z
Story: S-0500.
S-0500 keeps storygit.Committed's answer per repository root for the life of the process, keyed by git show-ref --heads, so a warm stats.get in flai serve no longer runs git log (1,317 ms to 648 ms on this repository on 2026-10-26). A fresh process keeps nothing: every one-shot flai stats, every flai hostapi --timing, and the first stats.get after flai serve starts or after any branch tip moves still read the whole git log --name-only over main and every story branch. On 2026-10-26 under a load of 21 to 39 that was 655 ms of a 2,444 ms cold stats.get, and flai stats --json took 3,458 ms against 3,558 ms before. Keeping the answer on disk under .flai-cache with the same key, or reading only the commits added since the kept tips, would remove it; neither was in S-0500's plan. Figures in design/system/server-performance.md, section The stories' commits kept between calls (S-0500).

## Remediation

Story S-0531 remediates this issue, created from it at 2026-10-28T12:54:38Z.
Closed 2026-10-28T18:36:41Z: Fixed by S-0531: storygit.Committed also keeps its answer in .flai-cache/stories-commits.json in the main checkout, with the repository root and S-0500's key (the main branch, the wip folder, and git show-ref --heads), so a new process whose root and key match answers from the file and starts no git log; a missing, damaged, or other-key file is a miss that reads the log and rewrites the file atomically, and a failed write never fails the call. On 2026-10-28 on this repository a cold stats.get fell from 2,154 to 1,751 ms and flai stats --json from 2,248 to 1,591 ms, main's git log having taken 552 ms of each, with identical claims.drift (design/system/server-performance.md, The stories' commits kept between processes). The tests in flai/internal/storygit/committed_test.go and TestReadStartsNoLogInANewProcessWhileNothingChanged in flai/internal/statsread/statsread_test.go reproduce it. The first call after any local branch tip moved still reads the whole log.
