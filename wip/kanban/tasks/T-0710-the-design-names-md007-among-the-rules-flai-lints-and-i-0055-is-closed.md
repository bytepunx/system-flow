---
id: T-0710
type: task
nature: remediation
title: The design names MD007 among the rules flai lints and I-0055 is closed
status: backlog
parent: S-0240
owner: arobson
created: 2026-10-02T16:50:12Z
updated: 2026-10-02T16:50:12Z
transitions: []
stream: S-0240
tags: []
touches: [design/system/flai-cli.md, design/issues]
after: [T-0708]
---
# T-0710 The design names MD007 among the rules flai lints and I-0055 is closed

## Work

Say in `design/system/flai-cli.md` that mdlint checks MD007 (ADR-0061 is accepted and stays as it is; its consequences foresee a rule added when agents start to break one), write the trace into the story's notes, and close I-0055 with `flai issue close`, naming the missing rule as the cause. Waits for T-0708, whose rule it describes; runs with T-0709, which touches no path of its own.

## Done when

- `flai-cli.md` names MD007 among the rules mdlint implements.
- The story's notes name the step the trace found.
- I-0055 is closed.

## Notes
