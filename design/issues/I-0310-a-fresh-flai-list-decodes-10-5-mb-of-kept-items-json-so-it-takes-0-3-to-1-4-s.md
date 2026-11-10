---
id: I-0310
title: A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names
class: efficiency
status: open
count: 1
cost: 1s
first_reported: 2026-11-08T01:17:08Z
last_reported: 2026-11-08T01:17:08Z
updated: 2026-11-10T08:34:53Z
---

# I-0310 A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names

## Description
A fresh flai list decodes 10.5 MB of kept items JSON, so it takes 0.3 to 1.4 s, not the 0.3 s S-0600's title names

## Instances

### 2026-11-08T01:17:08Z
Story: S-0600.
S-0600's T-2688 measured a fresh Repo.List(true) with the kept items files present at 0.94 s median, 0.31 s at the lowest load, against 4.54 s on main (load average 34 to 66 on 24 cores). Most of it is decoding .flai-cache/items-*.json: 10.5 MB of JSON, 6.0 MB of it the archived tasks and 3.8 MB the archived stories. Ways forward: a cheaper encoding (gob), decoding a kept item only when its key matches, or not listing the archive where only its epics' stories are needed (flai move's Follow).

## Remediation

Story S-0676 remediates this issue, created from it at 2026-11-10T08:34:53Z.
