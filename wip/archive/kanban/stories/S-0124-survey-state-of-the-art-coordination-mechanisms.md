---
id: S-0124
type: story
nature: research
title: Survey state of the art coordination mechanisms
status: done
parent: E-0009
owner: alex
created: 2026-09-26T07:17:08Z
updated: 2026-09-26T08:02:04Z
transitions:
  - to: ready
    at: 2026-09-26T07:17:53Z
    by: alex
  - to: in-progress
    at: 2026-09-26T07:18:07Z
    by: agent-S-0124
  - to: review
    at: 2026-09-26T08:00:06Z
    by: agent-S-0124
  - to: done
    at: 2026-09-26T08:02:04Z
    by: alex
tags: [cli]
touches: [flai/cmd]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0124 Survey state of the art coordination mechanisms

## Goal

Research state of the art options for agent orchestration and coordination such that flai can determine when agents can safely run their planned work in parallel vs. when they need to obey order and execute stories in serial.

The right solution will be able to explain to the user why a story is flagged yellow and pauses until the other stor(y|ies) before it are complete.

Use threads so we can discuss your findings and select a direction and document it in a new ADR.

## Acceptance criteria
- [x] A new ADR is written with decisions and input on direction by the operator
- [x] New stories are authored for this parent epic (0009) that will support the implementation of the mechanism(s) we choose

## Tasks
- T-0455 Survey coordination mechanisms and map what flai has today into a findings document
- T-0456 Discuss the findings with the designer in threads and settle a direction
- T-0457 Write the ADR for the chosen coordination mechanism
- T-0458 Author the stories under E-0009 that implement the chosen mechanism

## Notes

- Findings: `design/system/agent-coordination.md`. Direction chosen on TH-0018 (design 1) and TH-0019 (every policy recommendation), recorded in ADR-0046 (accepted, refines ADR-0019).
- Stories for E-0009, in pull order: S-0128 (hold in flai), S-0129 (yellow card), S-0130 (`after:`), S-0131 (trial merge and drift at sync), S-0132 (notice at acceptance). Left in backlog for the designer to move to ready.
