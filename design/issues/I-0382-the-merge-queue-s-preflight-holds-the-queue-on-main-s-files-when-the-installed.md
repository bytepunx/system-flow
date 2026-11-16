---
id: I-0382
title: The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed
class: blocker
status: closed
count: 1
cost: 3h
first_reported: 2026-11-16T00:41:30Z
last_reported: 2026-11-16T00:41:30Z
updated: 2026-11-16T03:54:40Z
---

# I-0382 The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed

## Description
The merge queue's preflight holds the queue on main's files when the installed flai is older than the flai main's files are written for, and the hold cannot clear until a release is installed

## Instances

### 2026-11-16T00:41:30Z
S-0748 was accepted at 21:46Z on 2026-11-15. It wrote the orchestrator and planner role files with when: {to: <state>} together with the code that reads that form. flai serve ran the installed flai 1.99.0, built at 21:33Z without S-0748, so the preflight's flai check reported roles.file: front matter does not parse on design/workflows/roles/orchestrator.md:15 and planner.md:28 and held the queue (TH-1231), telling the operator to fix the files in the main checkout. The files were right: a flai built from the tree reported no error on the same checkout. The publish story that would have released the fix, S-0793, sat in the held queue, and the orchestrator's and planner's own role files did not parse under the installed flai. The hold cleared by itself once 1.100.0 was installed and flai serve restarted.

## Impact
On 2026-11-15 the queue built no batch from 21:46Z until flai 1.100.0 was installed at about 00:22Z, with S-0793, S-0790, and S-0755 enqueued

- time_lost_per_cycle: 3h

## Remediation

Proposed, for the story to settle in an ADR: the preflight tells a finding the installed flai cannot judge from one in the files. When main's own flai sub-project is ahead of the running flai's commit, the preflight either runs its `flai check` with a flai built from main, as `scripts/flai.sh` does, or holds with words that name the cause, the installed flai's version against main's unreleased commits, and the remedy, publish and install the release, in place of telling the operator to fix files that are right. Consider also holding an acceptance that changes a format the running flai reads (role files, the workflow file, the manifest) until a flai that reads it is installed.

Story S-0797 remediates this issue, created from it at 2026-11-16T00:41:48Z.
Closed 2026-11-16T03:54:40Z: S-0797 (ADR-0311): when the preflight's in-process flai check finds an error and the running flai is behind main's flai sub-project (release.FlaiBehind), the preflight checks main again with the manifest's flai.command and holds only on what a flai not behind main finds, so a format change the installed flai cannot read no longer holds this repository's queue. TestThePreflightChecksMainAgainWithMainsFlaiWhenTheRunningFlaiIsBehind reproduces I-0382. Where the second check cannot settle it, the hold records flai_behind, and the thread and flai queue name the installed flai and say to publish and install the release rather than fix the files.
