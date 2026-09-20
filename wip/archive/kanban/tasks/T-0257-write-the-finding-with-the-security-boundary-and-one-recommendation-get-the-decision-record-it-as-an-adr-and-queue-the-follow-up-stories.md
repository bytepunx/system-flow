---
id: T-0257
type: task
nature: research
title: Write the finding with the security boundary and one recommendation, get the decision, record it as an ADR, and queue the follow-up stories
status: done
parent: S-0071
owner: alex
created: 2026-09-20T07:00:45Z
updated: 2026-09-20T07:27:10Z
transitions:
  - to: ready
    at: 2026-09-20T07:12:32Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:12:32Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:27:10Z
    by: system-flow
stream: S-0071
tags: []
---
# T-0257 Write the finding with the security boundary and one recommendation, get the decision, record it as an ADR, and queue the follow-up stories

## Work
State for each candidate what a token holder and a compromised container can make the host do; propose the boundary (named actions with typed arguments, enabled by the operator, confirmations or a second credential where needed) and how it sits with ADR-0018, ADR-0026, ADR-0027. How the host side runs and what the operator sees. One recommendation in design/system, indexed. Put it to the operator, record the decision with flai adr new, create the follow-up stories under E-0003 in build order.

## Done when
- The finding is in design/system with one recommendation
- The ADR is recorded
- The follow-up stories exist, or the notes say why not
- flai check --strict and the markdown lint pass

## Notes
