---
id: I-0016
title: Release picks the delivered component from the first matching tag even when the story's commits never touched it
class: defect
status: closed
count: 3
cost: 7m
first_reported: 2026-09-18T17:08:11Z
last_reported: 2026-09-18T17:08:11Z
updated: 2026-09-19T10:34:13Z
---

# I-0016 Release picks the delivered component from the first matching tag even when the story's commits never touched it

## Description
Release picks the delivered component from the first matching tag even when the story's commits never touched it

## Instances

### 2026-09-18T17:08:11Z
S-0039 was pure flai work but carried tags [dashboard, cli]; the release gave flaiover a minor bump with zero files and flai only a patch. The same shape hit S-0035 and S-0037. Two causes: I tagged every E-0006 story dashboard,cli when creating the epic, and deliverTarget takes the first tag that names a component without looking at the touched files. Versions stand as cut.

Counted retroactively: S-0035 (flai work, flaiover took the minor).

Counted retroactively: S-0037 (mostly flai work, flaiover took the minor).

## Remediation
Closed 2026-09-19T10:34:13Z: S-0047: the release delivers only to a component the story's commits touched; a tag naming an untouched component never delivers. Versions already cut stand.
