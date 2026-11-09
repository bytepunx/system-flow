---
id: I-0296
title: flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again
class: defect
status: closed
count: 1
first_reported: 2026-11-07T08:17:37Z
last_reported: 2026-11-07T08:17:37Z
updated: 2026-11-09T17:46:52Z
usage:
  strategic:
    - kind: analyzer
      seconds: 41
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 32
          output: 3102
          cache_read: 238316
          cache_write: 16129
          cost: 0.5376
        - model: claude-opus-5-5
          input: 14
          output: 2840
          cache_read: 759050
          cache_write: 13854
          cost: 0.278
---

# I-0296 flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again

## Description
flai serve records a Claude Code usage-limit refusal as a failed permission_prompt check and never checks the version again

## Instances

### 2026-11-07T08:17:37Z
Report: design/analysis/2026-11-07-all.md.
Analysis 2026-11-07: a 429 is not a prompt failure, yet protected-path writes under 2.1.311 were refused until the version moved on; the only remedy named is a hand edit of the host's state file.

## Impact
~/.flai/serve/claude-checks.json: 2.1.311 failed 2026-10-26T21:55:54Z with 'You've hit your weekly limit'; BEEFGOOF (the shim /home/alex/.local/bin/rmhclod) failed 2026-11-03 on a shell syntax error; ADR-0106 § 2: a version with a record is not checked again, whatever its outcome; TH-0669 and TH-0819 open on system-flow.yaml

## Remediation

Found by the analysis in [2026-11-07-all.md](../analysis/2026-11-07-all.md).

Story S-0629 remediates this issue, created from it at 2026-11-07T23:40:41Z.
Closed 2026-11-09T17:46:52Z: S-0629 fixed it. ADR-0270 refines ADR-0106 § 2: a permission_prompt check whose claude -p ends on the harness's usage limit, read with the ADR-0214 usage reader from the check's stream-json, records no outcome and opens no thread, and flai serve checks the version again after the reset, or an hour later when no reset time is named, on the host and in an agent image (flai/internal/serve/claudecheck.go). TestAClaudeCodeCheckRefusedOnTheUsageLimitIsCheckedAgainAfterTheReset replays the 2026-10-26 weekly-limit refusal; TestAUsageLimitWithNoResetTimeHoldsTheCheckAnHour and TestACheckInAnImageRefusedOnTheUsageLimitIsCheckedAgainAfterTheReset cover the rest. The 2.1.311 failed record already in ~/.flai/serve/claude-checks.json stays until the operator removes its entry.
