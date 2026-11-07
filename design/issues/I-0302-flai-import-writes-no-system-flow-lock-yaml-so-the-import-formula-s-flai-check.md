---
id: I-0302
title: flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-07T08:50:28Z
last_reported: 2026-11-07T08:50:28Z
updated: 2026-11-07T09:44:23Z
---

# I-0302 flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock

## Description
flai import writes no system-flow.lock.yaml, so the import formula's flai check --strict step fails on layout.lock

## Instances

### 2026-11-07T08:50:28Z
Story: S-0564.
S-0564's import run: the board's Import of /home/alex/git/flai-scratch-import (flai import --formula, flai 1.84.0, commit 1e7935b) recorded template 1.1.24 in system-flow.yaml but wrote no system-flow.lock.yaml. The instance S-0001's first step, flai check --strict, exited 1 on "layout.lock: template 1.1.24 is recorded but there is no system-flow.lock.yaml; run flai upgrade --relock". The agent blocked on TH-0001 for 6 minutes until the operator allowed flai upgrade --relock as a one-off (74 files relocked, commit f71d76d). Then check exited 0 and flai import test passed. So every import from the board fails its formula's first step, and the instance's branch changes system-flow.yaml and the lock, though the formula says its branch changes nothing. flai import --commit could relock as it imports.

## Remediation

Story S-0590 remediates this issue, created from it at 2026-11-07T09:44:23Z.
