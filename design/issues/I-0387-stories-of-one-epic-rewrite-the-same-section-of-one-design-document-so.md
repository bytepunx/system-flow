---
id: I-0387
title: Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle
class: efficiency
status: open
count: 2
cost: 18m
first_reported: 2026-11-16T06:16:15Z
last_reported: 2026-11-17T02:03:50Z
updated: 2026-11-17T02:26:33Z
usage:
  strategic:
    - kind: analyzer
      seconds: 99
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 47
          output: 6584
          cache_read: 349360
          cache_write: 34306
          cost: 1.1031
        - model: claude-haiku-5-5
          input: 36
          output: 14843
          cache_read: 2398931
          cache_write: 47318
          cost: 0.1557
---

# I-0387 Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle

## Description
Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle

## Instances

### 2026-11-16T06:16:15Z
Report: design/analysis/2026-11-16-all.md.
Analysis 2026-11-16: E-0033's stories each appended to the same section of design/system/server-performance.md; 12 acceptance send-backs in W46, each a sync, a second close-out (p50 474 s), and a second acceptance.

### 2026-11-17T02:03:50Z
Story: S-0797.
S-0797's acceptance (rebase and resolve job S-0797-20261117T020209Z) was sent back on TH-1263. The rebase stopped in design/system/flai-cli.md and docs/users/flai.md, and the section merge could not settle either file: "section flai CLI > Commands: both sides change the same lines". Each conflict was one long line: the `flai queue [--json]` table row in flai-cli.md (13 KB) and the `--json` paragraph of `flai queue` in flai.md. Main, through S-0799 and S-0805, and S-0797 had each inserted words at different places in that one line, with nothing removed. The sentence merge of ADR-0287 did not apply to a table row. The agent settled it by applying its own character-level insertions onto main's line at their unique anchors. A merge that applies two sides' disjoint insertions within one line, or within one table cell, would have landed it with no send-back.

## Impact
coordination.weeks[] W46 conflicts.met_accept 12 (W44 4, W45 4), met_sync 22. .flai-cache/conflicts.jsonl since 2026-11-09: design/system/server-performance.md in 5 of 12 acceptance stops, 6 of 22 sync stops, 20 of 28 trial conflicts; 19 commits to it since 2026-11-07 by 12 stories. Send-back threads TH-0920 and TH-0922 on S-0603 both name it within forty minutes on 2026-11-09. waiting.longest[]: S-0598 108,744 s, S-0601 107,713 s, S-0599 105,595 s in review after their send-backs. ADR-0276 keeps ADR-0175's limit: two sides that change the same lines of one section still stop.

- time_lost_per_cycle: 4h

## Remediation

Found by the analysis in [2026-11-16-all.md](../analysis/2026-11-16-all.md).

Story S-0817 remediates this issue, created from it at 2026-11-17T02:26:33Z.
