---
id: I-0124
title: flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-10-08T09:00:31Z
last_reported: 2026-10-08T09:00:31Z
updated: 2026-10-12T16:52:47Z
---

# I-0124 flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

## Description
flai stream sync replays a story's revert of its own change over an identical change on main, undoing main's

## Instances

### 2026-10-08T09:00:31Z
Story: S-0232.
S-0232 changed scripts/flai-snapshot.sh's skip list to --skip=publish,sign, then, as agreed on MS-0014, reverted the file to the main of that time. S-0340, accepted later, made the same change to the same lines. When flai stream sync rebased story/S-0232 onto main, the first commit applied as nothing and the revert applied cleanly, undoing S-0340's --skip=publish,sign on the story branch. No conflict was reported. The only sign was the sync's note that the branch changed a path outside the story's touches. Restoring main's file took a fix commit. A sync could compare a path the story's touches leave out with main and warn when the branch's net change there removes lines main gained since the branch point.

## Remediation

Story S-0371 remediates this issue, created from it at 2026-10-12T16:52:47Z.
