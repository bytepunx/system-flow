---
id: I-0012
title: Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative
class: defect
status: open
count: 2
cost: 5m
first_reported: 2026-09-17T06:42:13Z
last_reported: 2026-09-18T05:52:24Z
updated: 2026-09-18T05:52:24Z
---

# I-0012 Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative

## Description
Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative

## Instances

### 2026-09-17T06:42:13Z
A python edit step asserted and failed, but the following flai moves, check, and commit ran anyway because the heredoc step was not joined with &&. Caught by reading the narrative after the push; fixed in a follow-up commit. Same class as the earlier grep-masked test failures: every step in a chain must be gated on the previous exit code, including inline scripts.

### 2026-09-18T05:52:24Z
Second occurrence: 'flai check --strict | tail -1' masked a non-zero exit, so a story was queued and pushed with a check error, and the acceptance commit before it failed CI for the same finding. Every gate must be checked on its own exit code, never through a pipe.

## Remediation
