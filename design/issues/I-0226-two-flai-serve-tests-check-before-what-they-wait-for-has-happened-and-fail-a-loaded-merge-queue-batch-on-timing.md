---
id: I-0226
title: Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing
class: defect
status: open
count: 1
cost: 25m
first_reported: 2026-10-28T23:45:11Z
last_reported: 2026-10-28T23:45:11Z
updated: 2026-10-28T23:45:11Z
---

# I-0226 Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing

## Description
Two flai serve tests check before what they wait for has happened, and fail a loaded merge queue batch on timing

## Instances

### 2026-10-28T23:45:11Z
Story: S-0507.
Merge queue batch 20261028-232929 sent S-0507 back on TH-0724, though S-0507 changes nothing in flai/internal/serve. TestAStorySentBackIsStartedAfreshOrBeforeARestart wanted exactly one journal entry once the stub agent ended, and got the thread flai serve opens for an agent that ended with its story in progress as a second. TestRequestMergeIsAnsweredByFlaiServesMergePass read the merge's log line after waiting for its journal entry alone, and found the log empty. Both pass alone and failed in no earlier batch; that batch's mcpserver suite took 679 s beside them. S-0507's T-2317 makes the first check for one start and the second wait for the log line.

## Remediation
