---
id: T-0415
type: task
nature: remediation
title: "An ADR: flai serve starts a ready story's agent whenever the in-progress limit has room, whoever is attending"
status: done
parent: S-0116
owner: alex
created: 2026-09-24T09:00:31Z
updated: 2026-09-24T09:01:30Z
transitions:
  - to: ready
    at: 2026-09-24T09:00:42Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T09:00:43Z
    by: agent-S-0116
  - to: done
    at: 2026-09-24T09:01:30Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [design/adrs, design/system]
---
# T-0415 An ADR: flai serve starts a ready story's agent whenever the in-progress limit has room, whoever is attending

## Work

Write the ADR the designer asked for on TH-0010 with `flai adr new --supersedes`. It supersedes ADR-0042's hold, and the attended rule of ADR-0038 and ADR-0041 where it is still in force. It keeps ADR-0042's `flai mcp --agent` naming. It records the retirement of `attended_minutes` and the restart command and button.

## Done when

- The ADR is accepted in `design/adrs`, the index lists it, and the ADRs it supersedes carry `superseded_by`.
- `flai check --strict` is clean.

## Notes
