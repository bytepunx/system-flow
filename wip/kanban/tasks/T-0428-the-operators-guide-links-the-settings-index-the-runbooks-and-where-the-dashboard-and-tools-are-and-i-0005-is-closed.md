---
id: T-0428
type: task
nature: improvement
title: The operators guide links the settings index, the runbooks, and where the dashboard and tools are, and I-0005 is closed
status: done
parent: S-0028
owner: alex
created: 2026-09-24T09:28:03Z
updated: 2026-09-24T09:42:27Z
transitions:
  - to: ready
    at: 2026-09-24T09:28:10Z
    by: agent-S-0028
  - to: in-progress
    at: 2026-09-24T09:41:05Z
    by: agent-S-0028
  - to: done
    at: 2026-09-24T09:42:27Z
    by: agent-S-0028
stream: S-0028
tags: []
touches: [docs/operators, design/conventions, design/issues]
---
# T-0428 The operators guide links the settings index, the runbooks, and where the dashboard and tools are, and I-0005 is closed

## Work
- `docs/operators/index.md` opens with where things are: the settings index, the runbooks, the dashboard's address and pages, `flai serve journal`, `flai host status`, the logs and state folders.
- Close I-0005 with `flai issue close`, pointing at S-0028; update the documentation convention's project addition that calls operator docs a draft.

## Done when
- `docs/operators/index.md` links `settings.md`, every runbook, and the dashboard and tools; I-0005 is closed with a pointer to S-0028; `flai check --strict` and the markdown lint pass.

## Notes
