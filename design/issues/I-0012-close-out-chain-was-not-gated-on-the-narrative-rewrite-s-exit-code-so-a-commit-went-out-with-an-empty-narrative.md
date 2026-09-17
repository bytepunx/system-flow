---
id: I-0012
title: Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative
class: defect
status: open
count: 1
cost: 5m
first_reported: 2026-09-17T06:42:13Z
last_reported: 2026-09-17T06:42:13Z
updated: 2026-09-17T06:42:13Z
---

# I-0012 Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative

## Description
Close-out chain was not gated on the narrative rewrite's exit code, so a commit went out with an empty narrative

## Instances

### 2026-09-17T06:42:13Z
A python edit step asserted and failed, but the following flai moves, check, and commit ran anyway because the heredoc step was not joined with &&. Caught by reading the narrative after the push; fixed in a follow-up commit. Same class as the earlier grep-masked test failures: every step in a chain must be gated on the previous exit code, including inline scripts.

## Remediation
