---
id: I-0376
title: flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-15T01:05:45Z
last_reported: 2026-11-15T01:05:45Z
updated: 2026-11-17T02:47:36Z
---

# I-0376 flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out

## Description
flai guard counts a sub-agent that ended on an API error as still running and refuses wait_for_events for the close-out

## Instances

### 2026-11-15T01:05:45Z
Story: S-0739.
S-0739's T-3203 sub-agent (general-purpose, foreground) ended with an API error (SSL certificate hostname mismatch) and returned a partial result; flai guard still counted it as running, so it refused wait_for_events for the close-out job's verified event ('the story's agent cannot call wait_for_events while its sub-agents run (general-purpose ae3487c9a0041b35d)'), and TaskStop knew no such task. The agent had to poll the job's log with Bash instead.

## Remediation

Story S-0819 remediates this issue, created from it at 2026-11-17T02:47:36Z.
