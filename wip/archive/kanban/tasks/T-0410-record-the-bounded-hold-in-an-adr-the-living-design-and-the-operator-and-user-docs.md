---
id: T-0410
type: task
nature: remediation
title: Record the bounded hold in an ADR, the living design, and the operator and user docs
status: done
parent: S-0114
owner: alex
created: 2026-09-24T08:34:04Z
updated: 2026-09-24T08:44:05Z
transitions:
  - to: ready
    at: 2026-09-24T08:39:45Z
    by: agent-S-0114
  - to: in-progress
    at: 2026-09-24T08:39:46Z
    by: agent-S-0114
  - to: done
    at: 2026-09-24T08:44:05Z
    by: agent-S-0114
stream: S-0114
tags: []
touches: [design/adrs, design/system, docs/users, docs/operators]
---
# T-0410 Record the bounded hold in an ADR, the living design, and the operator and user docs

## Work

An ADR refining ADR-0038 and ADR-0041: someone attending holds a ready story back for the attended window, then flai serve starts it. Update the `flai serve agent` and `flai mcp` rows in `design/system/flai-cli.md`, the attending and waiting paragraphs in `docs/operators/index.md`, the `agent.attended_minutes` row and the `flai mcp` notes in `docs/users/flai.md`, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

## Done when

- [x] The ADR is accepted and indexed
- [x] Design and docs say what the launcher and `flai mcp --agent` now do
- [x] `flai check --strict` passes

## Notes
