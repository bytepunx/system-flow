---
id: I-0262
title: flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-02T21:38:10Z
last_reported: 2026-11-02T21:38:10Z
updated: 2026-11-07T04:07:05Z
---

# I-0262 flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it

## Description
flai serve does not start a defined role that takes items on a schedule entry under each-ready-story, though flai check passes it

## Instances

### 2026-11-02T21:38:10Z
Story: S-0554.
S-0554 implemented each-ready-story for the planner's launcher (ADR-0250 §4) and made flai check accept {schedule: <setting>, when: each-ready-story} on any role that takes items. The defined roles' schedule launcher, lookAt in flai/internal/serve/role_schedule.go, still skips a role that takes items, so such a defined role is logged "scheduled role not started" each time its schedule comes round and runs on no ready story. Before S-0554 it was skipped the same way; the check now lets the entry through. Closing it means one run per ready story with per-item held-start bookkeeping in role_schedule.go, or flai check reporting the entry on a defined role until then.

## Remediation

Story S-0571 remediates this issue, created from it at 2026-11-02T21:51:45Z.
Closed 2026-11-07T04:07:05Z: S-0571: flai serve's role scheduler (flai/internal/serve/role_schedule.go) now starts a defined role whose item lists story once for each ready story, in pull order, when a {schedule: <setting>, when: each-ready-story} entry comes round, keeping a budget-held start per story (T-2511, reproduced by TestAnEachReadyStoryScheduleStartsTheRoleOnEveryReadyStory); and flai check reports that entry as roles.file on a role whose item does not list story (T-2512).
