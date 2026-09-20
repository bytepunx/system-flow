---
id: T-0283
type: task
nature: feature
title: ADR-0031 supersedes ADR-0022, 0026, and 0027, says what of ADR-0018 holds, and I-0022's risks are reviewed
status: done
parent: S-0077
owner: alex
created: 2026-09-20T12:39:12Z
updated: 2026-09-20T12:45:52Z
transitions:
  - to: ready
    at: 2026-09-20T12:44:54Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:44:54Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:45:52Z
    by: system-flow
stream: S-0077
tags: []
---
# T-0283 ADR-0031 supersedes ADR-0022, 0026, and 0027, says what of ADR-0018 holds, and I-0022's risks are reviewed

## Work
One ADR: the container holds nothing of the project. What is superseded and why; what of ADR-0018 (the token, how it reaches the container) still holds; the decision on the container's user. I-0022's successor risks closed or restated in the issue.

## Done when
- flai check --strict clean; the three ADRs carry superseded_by

## Notes
