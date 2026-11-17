---
id: I-0374
title: flai guard counts a sub-agent that ended on an API error as still running, and refuses the story's agent wait_for_events for its close-out
class: defect
status: closed
count: 1
cost: 10m
first_reported: 2026-11-15T00:30:38Z
last_reported: 2026-11-15T00:30:38Z
updated: 2026-11-17T07:40:30Z
---

# I-0374 flai guard counts a sub-agent that ended on an API error as still running, and refuses the story's agent wait_for_events for its close-out

## Description
flai guard counts a sub-agent that ended on an API error as still running, and refuses the story's agent wait_for_events for its close-out

## Instances

### 2026-11-15T00:30:38Z
Story: S-0738.
S-0738's T-3202 task sub-agent (general-purpose ac8d9be303813c862), run in the foreground, ended with 'API Error: Unable to connect to API: SSL certificate hostname mismatch'. A second sub-agent ran and returned. Then wait_for_events for the close-out's verified event was refused: 'the story's agent cannot call wait_for_events while its sub-agents run (general-purpose ac8d9be303813c862)'. TaskStop finds no such task, so the guard's record of running sub-agents keeps one that ended on an error, apparently because no SubagentStop reached it. The agent waited on the close-out by polling the job log with a background Bash loop instead.

## Remediation
Closed 2026-11-17T07:40:30Z: Duplicate of I-0376, closed by S-0819 with the same fix: flai guard drops from a session's record, at the story agent's wait_for_events, each sub-agent whose foreground launch has its result in the session transcript (ADR-0315); S-0738's case, a result that is only an error string, is linked through the sub-agent's meta file.
