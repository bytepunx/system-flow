---
id: I-0207
title: The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked
class: blocker
status: open
count: 1
cost: 20m
first_reported: 2026-10-24T23:48:27Z
last_reported: 2026-10-24T23:48:27Z
updated: 2026-10-26T04:01:51Z
---

# I-0207 The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked

## Description
The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked

## Instances

### 2026-10-24T23:48:27Z
Story: S-0435.
S-0435 measured the digest's effect on 2026-10-24. flai stats reports asked_again 0 in every week, and digest.entries 382 active, 0 compacted, and 0 absorbed. None of these is a measurement. The 382 entries are S-0429's seed of 2026-10-23, and the host journal shows no curator run on this project: the curate host action was turned on only for signet, at 2026-10-24T22:41:18Z. A repeat counts only once the curator lists a thread in an entry's asked_again (ADR-0200), so the second half of S-0435's hypothesis cannot be tested here. Turn curate on for system-flow, let it run over the after window, and measure again once W44 and W45 are whole, from 2026-11-09: see design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md.

## Remediation

Story S-0504 remediates this issue, created from it at 2026-10-26T04:01:51Z.
