---
id: T-2119
type: task
nature: research
title: A draft ADR proposes custom dockerized agents on a flai baseline image, with the trade-offs against flai today and Docker Agent, its limitations, its security concerns, and the Kubernetes and signet directions
status: done
parent: S-0497
owner: alex
created: 2026-10-25T21:26:59Z
updated: 2026-10-25T21:35:33Z
transitions:
  - to: ready
    at: 2026-10-25T21:31:15Z
    by: system-flow
  - to: in-progress
    at: 2026-10-25T21:31:16Z
    by: system-flow
  - to: done
    at: 2026-10-25T21:35:32Z
    by: system-flow
stream: S-0497
tags: []
touches: [design/adrs/README.md, design/system/agent-adapters.md]
after: [T-2118]
usage:
  source: log
  seconds: 256
  estimated: true
  priced_by: harness
  models:
    - model: claude-fable-5-1
      input: 241
      output: 11010
      cache_read: 1387536
      cache_write: 63421
      cost: 2.1682
---
# T-2119 A draft ADR proposes custom dockerized agents on a flai baseline image, with the trade-offs against flai today and Docker Agent, its limitations, its security concerns, and the Kubernetes and signet directions

## Work

Write the draft ADR with `flai adr new --body-stdin --commit`, status proposed, refining ADR-0130 and ADR-0131 and citing the subsection T-2118 writes. Its context says why Docker Agent is not a clear win; its decision names the baseline image, the harness flai writes, what flai controls (mounts, tools, `flai mcp`, `flai guard`, the permission prompt, cost capture per harness, model, and agent); its consequences carry the trade-offs against flai today and Docker Agent, the limitations, the security concerns, and the Kubernetes and signet directions; its alternatives name flai today, Docker Agent, and a loop of flai's own. Waits for T-2118 because it cites the subsection by heading.

## Done when

- The ADR exists under `design/adrs` with its index row, status proposed, and `flai check --strict` passes.
- It carries trade-offs, limitations, security concerns, and the Kubernetes and signet directions, each in its own list.
- Acceptance criteria 1 and 2 of S-0497 are verifiable from the ADR and the subsection together.

## Notes
