---
id: T-0248
type: task
nature: research
title: Write the finding with the security of each option and one recommendation, get the decision, record it, and queue the flai integration story
status: cancelled
parent: S-0069
owner: alex
created: 2026-09-20T04:07:28Z
updated: 2026-09-20T06:15:05Z
transitions:
  - to: cancelled
    at: 2026-09-20T06:15:05Z
    by: system-flow
stream: S-0069
tags: []
---
# T-0248 Write the finding with the security of each option and one recommendation, get the decision, record it, and queue the flai integration story

## Work
State the security of each option plainly: flairport has no authentication of its own and holds every project's token, so the edge check is the only lock. Cover what happens when the tunnel is up and the edge policy is missing or wrong, whether the origin can refuse requests that did not come through the edge (a signed header, a token, binding to loopback), how access is revoked, and whether authstar could be the identity check behind a plain tunnel. Write the finding under `design/system` with one recommendation and why each other option lost. Put the recommendation to the operator; record their decision as an ADR refining ADR-0018 and ADR-0024 when it changes what is exposed or how, and as a note in the living design otherwise. If one or more options are suitable, create the follow-up story under E-0007 for flai CLI integration with criteria drawn from the finding; if none is, say why in the story's notes.

## Done when
- The finding is in `design/system`, linked from its index, with one recommendation
- The decision is recorded as an ADR or a design note
- The follow-up story exists under E-0007, or the notes say why not
- `flai check --strict` and the markdown lint pass

## Notes
- 2026-09-20T06:15:05Z: moved to cancelled: S-0069 cancelled: I am thinking of taking a very different route
