---
id: I-0312
title: "A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path"
class: efficiency
status: open
count: 1
cost: 2s
first_reported: 2026-11-08T01:58:32Z
last_reported: 2026-11-08T01:58:32Z
updated: 2026-11-08T01:58:32Z
---

# I-0312 A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path

## Description
A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path

## Instances

### 2026-11-08T01:58:32Z
Story: S-0601.
S-0601 ran check.Run's rules at once and kept the props node read. A warm `flai check --strict` on this repository went from 7.9 to 11.7 s wall (4.9 to 5.3 s CPU) to 3.3 to 4.9 s wall (3.8 to 3.9 s CPU), at a load average of 84 to 103 on 24 cores. That is still over the story's 2 s. Timed in process at a load of 103 to 114, the critical path is: `Repo.List(true)` with S-0600's kept items, 0.66 to 1.10 s; `workItems`, which runs alone because it fills byID, about 0.45 to 0.8 s; then the slowest rule of those run at once, `digest`, 0.76 to 1.77 s, which reads and scans every digest entry; then `conflictMarkers`, 0.45 to 1.10 s. Each check, `flai verify`, close-out check step, and `flai task done` pays it. Possible remedies: run `oneItem` over the items on the pool inside `workItems`, and keep the digest rule's per-file results under .flai-cache keyed by stat, as S-0600 keeps the lint findings.

## Remediation
