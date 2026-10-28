---
id: I-0209
title: A one-shot flai stats, and the first stats.get of a flai serve process, still run one git log for claims drift, about a quarter of a cold answer
class: efficiency
status: open
count: 1
cost: 650ms
first_reported: 2026-10-26T04:42:46Z
last_reported: 2026-10-26T04:42:46Z
updated: 2026-10-28T12:54:38Z
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
