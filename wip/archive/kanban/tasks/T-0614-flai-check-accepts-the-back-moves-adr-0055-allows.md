---
id: T-0614
type: task
nature: feature
title: flai check accepts the back moves ADR-0055 allows
status: done
parent: S-0173
owner: arobson
created: 2026-10-01T07:42:45Z
updated: 2026-10-01T07:46:34Z
transitions:
  - to: ready
    at: 2026-10-01T07:42:51Z
    by: agent-S-0173
  - to: in-progress
    at: 2026-10-01T07:46:34Z
    by: agent-S-0173
  - to: done
    at: 2026-10-01T07:46:34Z
    by: agent-S-0173
stream: S-0173
tags: []
touches: [flai/internal/check]
usage:
  source: log
  seconds: 0
  models: []
---

# T-0614 flai check accepts the back moves ADR-0055 allows

## Work

S-0167 let `flai move` take an item back a column (ADR-0055: ready to backlog, in-progress to ready, cancelled to backlog) but left `flai check`'s own table of allowed transitions as it was, so a history flai wrote fails `item.sequence`. S-0173's own restart on this host (in-progress to ready) made `flai check --strict` fail on main. check takes the allowed transitions from workitem, so the two cannot drift again.

## Done when

A history with each ADR-0055 back move passes `item.sequence`, done followed by anything still fails, and `flai check --strict` passes on this repository; `make flai-test` passes.

## Notes
