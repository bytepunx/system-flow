---
id: T-0171
type: task
nature: research
title: Write the finding with a recommendation, get the decision, record it, and queue what follows
status: done
parent: S-0052
owner: alex
created: 2026-09-19T01:55:25Z
updated: 2026-09-19T07:56:35Z
transitions:
  - to: ready
    at: 2026-09-19T07:56:34Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:56:34Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:56:35Z
    by: system-flow
stream: S-0052
tags: []
---

# T-0171 Write the finding with a recommendation, get the decision, record it, and queue what follows

## Work
Write the finding in `design/system` (or as a proposed ADR if the recommendation changes what the container is given): the options side by side, one recommendation, why each other option lost, and what was and was not tried. Present it to the operator and record the decision: an ADR refining ADR-0018 when the security posture or the container's inputs change, a living-design note otherwise, including when the decision is to change nothing. Queue the work that follows as stories under E-0006 with goal and acceptance criteria, or amend S-0041 if it belongs there. Update `docs/operators/index.md` only if the decision changes what operators are told today.

## Done when
- The finding exists, is linked from the story, and ends with one recommendation
- The operator's decision is recorded where the criteria say
- Follow-up stories exist, or the story notes say why none are needed
- `flai check --strict` is clean

## Notes
