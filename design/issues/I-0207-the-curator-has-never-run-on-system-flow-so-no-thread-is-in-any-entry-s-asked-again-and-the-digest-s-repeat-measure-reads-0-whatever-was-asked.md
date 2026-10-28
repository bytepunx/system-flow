---
id: I-0207
title: The curator has never run on system-flow, so no thread is in any entry's asked_again and the digest's repeat measure reads 0 whatever was asked
class: blocker
status: closed
count: 1
cost: 20m
first_reported: 2026-10-24T23:48:27Z
last_reported: 2026-10-24T23:48:27Z
updated: 2026-10-28T12:14:54Z
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
Closed 2026-10-28T12:14:54Z: S-0504 fixes both halves of the cause.

1. **Measure (ADR-0204).** `flai stats` reads the curator's activity document and reports `digest.curator` (`runs`, `last_run`) and `digest.weeks[].curator_runs`. `asked_again` is left out of any week no curator run has reached, so a project whose curator never ran no longer reads 0 repeats as a measurement. The Digest chart leaves those weeks a gap and notes when the curator has never run. The analyzer reads `digest.curator` before `asked_again`. `TestDigestLeavesAskedAgainOutWithoutACuratorRunOfI0207` reproduces the bug. On this repository, the branch's `flai stats` reads W40 to W44 as `asked again –` and `curator  runs 0 · never run`.
2. **Operation.** The operator turned `curate` on for system-flow (TH-0649). Its first run, started by the acceptances of S-0498 and S-0501, failed on Claude Code's usage limit and logged nothing; I-0212 records that. The operator chose on TH-0649 to let the next acceptance start the first run.
