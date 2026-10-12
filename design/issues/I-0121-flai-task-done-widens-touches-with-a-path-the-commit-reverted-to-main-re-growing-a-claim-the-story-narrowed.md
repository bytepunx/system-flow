---
id: I-0121
title: flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-10-08T08:25:05Z
last_reported: 2026-10-08T08:25:05Z
updated: 2026-10-12T16:52:44Z
---

# I-0121 flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed

## Description
flai task done widens touches with a path the commit reverted to main, re-growing a claim the story narrowed

## Instances

### 2026-10-08T08:25:05Z
Story: S-0232.
S-0232 agreed on MS-0014 to leave scripts/flai-snapshot.sh to S-0340. It removed the path from T-1223's and the story's touches, reverted the file to main's version, and closed T-1223 again with flai task done. task done added scripts/flai-snapshot.sh back to both touches, because the revert commit changed it, though the branch no longer differs from main there. The claim overlapped S-0340 again until the touches were narrowed by hand a second time. Widening could use the branch's net diff against main (git diff main...HEAD) rather than the paths of the last commit.

## Remediation

Story S-0370 remediates this issue, created from it at 2026-10-12T16:52:44Z.
