---
id: T-0393
type: task
nature: remediation
title: Docs say a repeat of a detached write is answered from a record that outlives flai serve
status: done
parent: S-0109
owner: alex
created: 2026-09-24T07:37:31Z
updated: 2026-09-24T07:42:08Z
transitions:
  - to: ready
    at: 2026-09-24T07:37:34Z
    by: agent-S-0109
  - to: in-progress
    at: 2026-09-24T07:41:34Z
    by: agent-S-0109
  - to: done
    at: 2026-09-24T07:42:08Z
    by: agent-S-0109
stream: S-0109
tags: []
touches: [design/system/flai-cli.md, docs/operators/index.md, design/system/flaiover-dashboard.md]
---
# T-0393 Docs say a repeat of a detached write is answered from a record that outlives flai serve

## Work

- `design/system/flai-cli.md`: the host actions row says how a detached write is recorded and how a repeat is answered.
- `docs/operators/index.md`: `serve/requests.json` beside the journal, and what it is for.
- `design/system/flaiover-dashboard.md`: `Repo.write()`'s retry is answered from that record for a detached write.

## Done when

- The three documents say what the code does, and `flai check --strict` is clean.

## Notes
