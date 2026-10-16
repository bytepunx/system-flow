---
id: T-1624
type: task
nature: feature
title: An ADR records that under pull-request acceptance every acceptance, the orchestrator's included, waits for the required checks and reviews
status: backlog
parent: S-0395
owner: alex
created: 2026-10-16T04:38:08Z
updated: 2026-10-16T04:38:08Z
transitions: []
stream: S-0395
tags: [cli]
touches: [design/adrs]
---
# T-1624 An ADR records that under pull-request acceptance every acceptance, the orchestrator's included, waits for the required checks and reviews

## Work

- Write the ADR with `flai adr new`, refining ADR-0093 and S-0394's ADR. It records:
  - the blockers
  - that they come from the kept state at the branch's head
  - that the orchestrator obeys them and never approves a review, as TH-0463 answers question 3
  - that the operator's acceptance obeys them too, unlike an open thread, which only the orchestrator's acceptance waits for
- First in its layer: it waits for no task.

## Done when

The ADR is accepted, and `flai check --strict` passes on it.

## Notes
