---
id: T-0710
type: task
nature: remediation
title: The design names MD007 among the rules flai lints and I-0055 is closed
status: done
parent: S-0240
owner: arobson
created: 2026-10-02T16:50:12Z
updated: 2026-10-02T16:57:38Z
transitions:
  - to: ready
    at: 2026-10-02T16:56:55Z
    by: agent-S-0240
  - to: in-progress
    at: 2026-10-02T16:56:55Z
    by: agent-S-0240
  - to: done
    at: 2026-10-02T16:57:38Z
    by: agent-S-0240
stream: S-0240
tags: []
touches: [design/system/flai-cli.md, design/issues]
after: [T-0708]
usage:
  source: log
  seconds: 43
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 29
      output: 8949
      cache_read: 1099742
      cache_write: 30289
      cost: 0.5915
---
# T-0710 The design names MD007 among the rules flai lints and I-0055 is closed

## Work

Say in `design/system/flai-cli.md` that mdlint checks MD007 (ADR-0061 is accepted and stays as it is; its consequences foresee a rule added when agents start to break one), write the trace into the story's notes, and close I-0055 with `flai issue close`, naming the missing rule as the cause. Waits for T-0708, whose rule it describes; runs with T-0709, which touches no path of its own.

## Done when

- `flai-cli.md` names MD007 among the rules mdlint implements.
- The story's notes name the step the trace found.
- I-0055 is closed.

## Notes
