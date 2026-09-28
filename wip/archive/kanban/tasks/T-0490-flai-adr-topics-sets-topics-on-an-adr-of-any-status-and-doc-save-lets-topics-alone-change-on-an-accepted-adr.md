---
id: T-0490
type: task
nature: feature
title: flai adr topics sets topics on an ADR of any status, and doc save lets topics alone change on an accepted ADR
status: done
parent: S-0134
owner: alex
created: 2026-09-27T03:58:47Z
updated: 2026-09-27T04:01:36Z
transitions:
  - to: ready
    at: 2026-09-27T03:58:54Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-27T03:59:53Z
    by: agent-S-0134
  - to: done
    at: 2026-09-27T04:01:36Z
    by: agent-S-0134
stream: S-0134
tags: []
touches: [flai/internal/adr, flai/internal/docedit, flai/cmd]
---
# T-0490 flai adr topics sets topics on an ADR of any status, and doc save lets topics alone change on an accepted ADR

## Work

adr.SetTopics writes or replaces the topics key and nothing else, with the check-and-undo of New and Accept; `flai adr topics ADR-nnnn <topics…>` with --autocommit and --trailer. docedit.Save allows an accepted ADR's save when the only change is its topics key and refuses anything else with a reason that names flai adr topics.

## Done when

Tests show topics set on accepted and proposed ADRs with every other byte unchanged, an unknown topic refused, doc save accepting a topics-only change and refusing a body change; make test passes.

## Notes
