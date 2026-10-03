---
id: I-0057
title: flai check --strict stops a story's close-out on wip/ findings outside the story
class: efficiency
status: open
count: 10
cost: 7m
first_reported: 2026-10-02T16:13:36Z
last_reported: 2026-10-03T21:37:49Z
updated: 2026-10-03T21:37:49Z
---

# I-0057 flai check --strict stops a story's close-out on wip/ findings outside the story

## Description
flai check --strict stops a story's close-out on wip/ findings outside the story

## Instances

### 2026-10-02T16:13:36Z
2026-10-02: S-0191's close-out stopped at flai check --strict on one warning, threads.archived on TH-0067 (answered, its story S-0231 archived), which the story did not change and its agent may not resolve: the thread still asks the operator a question (point 5). S-0176 (TH-0032), S-0181 (14 warnings, none its own), and S-0187 (TH-0032) stopped the same way, and each went to review with the finding noted. TH-0056's answer: a finding outside the story is a note, and flai records it as an issue or bumps its count; until flai does, the story's agent records it here.

### 2026-10-03T06:49:36Z
Story: S-0199.
S-0199's close-out stopped at flai check --strict on one warning that is not its own: story.unaccepted on the archived S-0173, whose branch story/S-0173 still exists unmerged. S-0199 changed nothing there; the branch is the operator's to merge or delete. Gone to review with the finding noted.

### 2026-10-03T07:38:13Z
Story: S-0207.
scripts/close-out.sh S-0207 stopped at flai check --strict on 10 warnings, none from S-0207's changes: 9 wip.overlap between S-0207 and S-0200 and its task T-0752, which are in progress side by side and are settling the overlap on TH-0084, and story.unaccepted on S-0173's unmerged branch. The lint, the whole suite, and the template smoke had passed.

### 2026-10-03T08:15:16Z
Story: S-0201.
S-0201's close-out stopped at flai check --strict on story.unaccepted for S-0173's unmerged branch and wip.overlap with S-0200, none in S-0201's diff; main fails the same. The verifier ran the remaining steps (flaiover tests, markdown lint) by hand.

### 2026-10-03T17:53:04Z
Story: S-0200.
S-0200's close-out stopped at flai check --strict on two findings outside it: story.unaccepted on S-0173's unmerged branch, and wip.overlap with T-0753 of S-0201, which touches flai/internal/workitem/boardview.go inside S-0200's claim. Every other step passed.

### 2026-10-03T18:03:01Z
Story: S-0200.
S-0200's second close-out stopped in the integration tier: TestMonorepoIsClean failed on issues.duplicate-id, two I-0062 files on main from S-0207's acceptance, which story/S-0201 renames. Every finding was main's own. The verifier ran the narrative and ancestry checks by hand.

### 2026-10-03T18:29:38Z
Story: S-0203.
S-0203's close-out stopped in the integration tier: TestMonorepoIsClean failed on main's issues.duplicate-id (two I-0062 files), with story.unaccepted on S-0173, epic.lags-stories on E-0015, and threads.archived on TH-0087; main's check reports the same four. None is in S-0203's diff. The verifier ran the smoke, check, markdown lint, and install tests by hand; all passed but the check on main's findings.

### 2026-10-03T18:50:07Z
Story: S-0202.
S-0202's close-out stopped at flai check --strict on findings outside its diff: story.unaccepted on S-0173's unmerged branch, epic.lags-stories on E-0015, wip.overlap between S-0202 and S-0204 (in progress side by side, both describing their dashboard changes in design/system/flaiover-dashboard.md and docs/users/flaiover.md; flai stream sync reports they merge cleanly), and wip.overlap among S-0204, S-0206, and their tasks. The verifier ran the steps the script did not reach by hand: flaiover-test.sh passed (809 tests), the markdown lint passed, and the narrative and ancestry checks hold.

### 2026-10-03T19:42:28Z
Story: S-0204.
S-0204's close-out stopped at flai check --strict on warnings outside its diff: story.unaccepted for archived S-0173, epic.lags-stories on E-0015, and wip.overlap between S-0204 and S-0206 (parallel in-progress stories whose touches meet; settled on TH-0091..TH-0094). The verifier ran the remaining steps by hand; all passed.

### 2026-10-03T19:49:13Z
Story: S-0206.
S-0206's close-out stopped at flai check --strict on two warnings outside the story, which main's checkout reports too: story.unaccepted on archived S-0173 (branch never merged) and epic.lags-stories on E-0015. Every other step passed (lint, three test tiers, markdown lint, installer smoke). Following TH-0056, S-0206 went to review with them recorded in its narrative.

### 2026-10-03T21:37:49Z
Story: S-0205.
S-0205's close-out stopped at flai check --strict on three warnings outside its diff, which main reports too: story.unaccepted on archived S-0173 (branch never merged), epic.lags-stories on E-0015, and threads.archived on TH-0094 (answered, S-0206 archived). The verifier ran the steps the script did not reach by hand.

## Remediation
