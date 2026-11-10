---
id: I-0338
title: flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync
class: defect
status: open
count: 1
cost: 3m
first_reported: 2026-11-10T09:13:03Z
last_reported: 2026-11-10T09:13:03Z
updated: 2026-11-10T10:10:59Z
---

# I-0338 flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync

## Description
flai task done's sync reports undoes_main for paths the story never changed when main gains a commit during the sync

## Instances

### 2026-11-10T09:13:03Z
Story: S-0682.
S-0682's T-2899 close (task_done, 2026-11-10T09:11:52Z): the sync step answered undoes_main for design/issues/I-0259-….md and wip/kanban/stories/S-0706-….md, with old_base ea0dce42b, though story/S-0682 changes neither path (git diff main...HEAD lists only the story's ten files). main had gained commit fe3b4e5 (create S-0706, I-0259's story at 09:11:37Z) while the close ran, and the branch's merge base was still 0bcb906. A second flai stream sync S-0682 a minute later answered undoes_main empty and left no record under .flai-cache/undone. Cost: checking the diff and the record to be sure acceptance was not blocked.

## Remediation

Story S-0712 remediates this issue, created from it at 2026-11-10T10:10:59Z.
