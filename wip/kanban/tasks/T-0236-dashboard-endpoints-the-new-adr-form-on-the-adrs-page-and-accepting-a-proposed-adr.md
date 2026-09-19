---
id: T-0236
type: task
nature: feature
title: "Dashboard: endpoints, the new ADR form on the ADRs page, and accepting a proposed ADR"
status: in-progress
parent: S-0060
owner: alex
created: 2026-09-19T09:39:14Z
updated: 2026-09-19T09:47:15Z
transitions:
  - to: ready
    at: 2026-09-19T09:47:15Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T09:47:15Z
    by: system-flow
stream: S-0060
tags: []
---

# T-0236 Dashboard: endpoints, the new ADR form on the ADRs page, and accepting a proposed ADR

## Work
`POST /api/adrs` `{ title, status?, supersedes?, refines?, body }`, `GET /api/adrs/template`, and `POST /api/adrs/:id/accept`, all through flai, shaped like S-0059's item endpoints (values as `--flag=value`, the title after `--`, 422 with findings). A "new ADR" action on the ADRs page when the dashboard can write, opening a form: title, status (proposed by default), supersedes and refines chosen from the existing ADRs, and the markdown starting from the template's sections with the explorer's preview beside it; a refusal keeps the text. After creating, the operator lands on the new ADR in the explorer and the list shows it. A proposed ADR's row or page has an accept action. Component and endpoint tests.

## Done when
- The endpoints and the form are tested
- `make flaiover-test` and `make flaiover-build` pass

## Notes
