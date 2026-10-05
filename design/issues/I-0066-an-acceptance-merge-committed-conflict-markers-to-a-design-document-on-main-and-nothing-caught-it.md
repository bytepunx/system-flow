---
id: I-0066
title: An acceptance merge committed conflict markers to a design document on main, and nothing caught it
class: defect
status: closed
count: 1
cost: 5m
first_reported: 2026-10-03T18:34:51Z
last_reported: 2026-10-03T18:34:51Z
updated: 2026-10-05T03:20:37Z
---

# I-0066 An acceptance merge committed conflict markers to a design document on main, and nothing caught it

## Description
An acceptance merge committed conflict markers to a design document on main, and nothing caught it

## Instances

### 2026-10-03T18:34:51Z
Story: S-0204.
design/system/work-hierarchy.md on main carries <<<<<<< HEAD / ======= / >>>>>>> f0f5443 (S-0201) around the 'Who changes what' rule after S-0200 and S-0201 were accepted; flai check --strict does not look for conflict markers in markdown. Resolved in S-0204's work-hierarchy edit by keeping both sides.

## Remediation
Closed 2026-10-05T03:20:37Z: S-0253: flai check reports markdown.conflict-marker, an error, on each conflict marker line in the markdown under design, docs, wip, and the root, and flai accept refuses, merging nothing, a story branch any of whose added or changed files carries a marker, naming each path and line
