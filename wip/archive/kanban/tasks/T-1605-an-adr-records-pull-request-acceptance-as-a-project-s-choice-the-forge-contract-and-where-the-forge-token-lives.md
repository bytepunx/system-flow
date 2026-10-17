---
id: T-1605
type: task
nature: feature
title: An ADR records pull-request acceptance as a project's choice, the forge contract, and where the forge token lives
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:36Z
updated: 2026-10-17T07:03:20Z
transitions:
  - to: ready
    at: 2026-10-17T07:03:18Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T07:03:18Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:03:20Z
    by: agent-S-0392
stream: S-0392
tags: [cli]
touches: [design/adrs]
usage:
  source: log
  seconds: 2
  estimated: true
  priced_by: harness
  models: []
---
# T-1605 An ADR records pull-request acceptance as a project's choice, the forge contract, and where the forge token lives

## Work

- Write the ADR with `flai adr new`. It records:
  - `acceptance.mode` as the project's choice, with local acceptance the default
  - the forge contract, and GitHub as its first adapter, over REST
  - the forge token in the variable the host names, as a provider's key is (ADR-0129), never in the manifest, a log, an item, or a start record
  - why REST rather than `gh`
- Refine ADR-0032 and ADR-0129 in its front matter.
- Leave what the later stories decide to their own ADRs: the push and the pull request on the move to review, threads from checks, the acceptance gate, and the wait in metrics.
- First in its layer: it waits for no task. The decision is written before the code it describes is closed.

## Done when

The ADR is accepted, names the decisions above, and `flai check --strict` passes on it.

## Notes
