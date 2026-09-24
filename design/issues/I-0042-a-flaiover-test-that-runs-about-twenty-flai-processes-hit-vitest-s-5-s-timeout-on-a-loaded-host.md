---
id: I-0042
title: A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-09-24T09:14:39Z
last_reported: 2026-09-24T09:14:39Z
updated: 2026-09-24T09:14:39Z
---

# I-0042 A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host

## Description
A flaiover test that runs about twenty flai processes hit vitest's 5 s timeout on a loaded host

## Instances

### 2026-09-24T09:14:39Z
S-0116, 2026-09-24: writes.test.ts 'changes nothing until the shell turns settings on' took 4.6 to 5.1 s with a flai built from main and from story/S-0116 alike, load average 5.5 while another flai-started agent ran its own tests. It failed at the 5 s default on every run. Given 30 s in S-0116.

## Remediation
