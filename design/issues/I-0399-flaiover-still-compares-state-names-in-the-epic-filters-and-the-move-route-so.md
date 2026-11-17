---
id: I-0399
title: flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-17T08:25:17Z
last_reported: 2026-11-17T08:25:17Z
updated: 2026-11-17T12:50:34Z
---

# I-0399 flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there

## Description
flaiover still compares state names in the epic filters and the move route, so an operator's own accepted or cancelled lane is not seen as closed there

## Instances

### 2026-11-17T08:25:17Z
Story: S-0758.
S-0758 moved the fifteen views its criteria list, and their callers, to read a state's kind from workflow.ts. A search after that found state names still compared outside its list: the open-epic filter `status !== 'done' && status !== 'cancelled'` in flaiover/src/routes/design/+page.svelte, lib/components/DeploySettings.svelte, and lib/components/NewItemForm.svelte; `body.to === 'cancelled'` (the dry-run preview) and `body.to === 'done'` (the acceptance's longer timeout) in routes/api/items/[id]/move/+server.ts; and the move to `ready` that lib/components/design/ProposeChange.svelte posts. With a workflow whose accepted or cancelled lane has another name, those epic lists keep closed epics, and the move route treats a cancel or an acceptance as a plain move. Each should read the flow from board.get as the other views now do.

## Remediation

Story S-0838 remediates this issue, created from it at 2026-11-17T08:32:49Z.
Closed 2026-11-17T12:50:34Z: S-0838: each site now reads the board's flow and asks $lib/workflow what a state is. The open-epic lists of the Design page, NewItemForm, and DeploySettings use isOpen; ProposeChange moves to the first queue lane isPullable names; and the move route uses isCancelled for the dry-run preview (refusing a dry run to any other lane) and isAccepted for the acceptance's timeout. renamedFlow in lanes.fixture.ts, with ready, done, and cancelled renamed, reproduces the cause in kinds.test.ts and in each view's test.
