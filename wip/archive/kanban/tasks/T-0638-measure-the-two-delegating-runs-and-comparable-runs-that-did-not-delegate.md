---
id: T-0638
type: task
nature: research
title: Measure the two delegating runs and comparable runs that did not delegate
status: done
parent: S-0188
owner: arobson
created: 2026-10-01T08:31:38Z
updated: 2026-10-01T09:18:51Z
transitions:
  - to: ready
    at: 2026-10-01T08:31:52Z
    by: agent-S-0188
  - to: in-progress
    at: 2026-10-01T09:01:26Z
    by: agent-S-0188
  - to: done
    at: 2026-10-01T09:18:51Z
    by: agent-S-0188
stream: S-0188
tags: []
touches: [design/system/agent-context.md]
usage:
  source: log
  seconds: 1045
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 13669
      cache_read: 2163328
      cache_write: 48896
      cost: 1.0437
---
# T-0638 Measure the two delegating runs and comparable runs that did not delegate

## Work

Read each run's log as `agent-context.md § Sub-agents › Measured` reads S-0175 and S-0118 (ADR-0051: assistant events deduplicated by message ID, the `result` for cost and turns, `parent_tool_use_id` for a sub-agent's calls): cost, cache reads, model calls, and turns for the story's agent and its sub-agents, and wall time. Pick comparable runs that did not delegate by nature, size, and components touched.

## Done when

Each delegating run has its numbers and at least one comparable non-delegating run, with what makes them comparable.

## Notes
- Measured at 09:18Z from the kept logs with `/tmp/s0188/measure.py` (ADR-0051's reading; sub-agent calls by `parent_tool_use_id`; sub-agent cost apportioned by its share of input and cache tokens). The numbers are in `agent-context.md § Sub-agents › Measured`.
- S-0184 (improvement, 39 files, +1092 −85, flai and flaiover) set against S-0174 (remediation, 30 files, +1038 −51, flai and flaiover). S-0185 (remediation, 27 files, +433 −82, flai) set against S-0180 (12 files, +140 −28) and S-0173 (30 files, +746 −134). All of them are claude-opus-5-5 at high effort, on this host, on 2026-10-01.
- Both agents still ran the whole Go suite and targeted tests themselves, and the verifier ran them again. S-0184's verifier found two convention gaps and S-0185's found five review points, and both agents fixed them before review.
