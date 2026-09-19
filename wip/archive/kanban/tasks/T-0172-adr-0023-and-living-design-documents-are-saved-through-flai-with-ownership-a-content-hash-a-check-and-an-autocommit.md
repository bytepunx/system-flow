---
id: T-0172
type: task
nature: feature
title: "ADR-0023 and living design: documents are saved through flai, with ownership, a content hash, a check, and an autocommit"
status: done
parent: S-0040
owner: alex
created: 2026-09-19T02:11:29Z
updated: 2026-09-19T02:13:23Z
transitions:
  - to: ready
    at: 2026-09-19T02:12:19Z
    by: alex
  - to: in-progress
    at: 2026-09-19T02:12:19Z
    by: alex
  - to: done
    at: 2026-09-19T02:13:23Z
    by: alex
stream: S-0040
tags: []
touches: [design/adrs, design/system]
---

# T-0172 ADR-0023 and living design: documents are saved through flai, with ownership, a content hash, a check, and an autocommit

## Work
Write ADR-0023 refining ADR-0016 (the dashboard delegates writes to flai): documents edited in the dashboard are saved by a flai command, which decides what may be edited, detects a concurrent change by a content hash, validates with the check, and commits. Record the rules. Ownership: work items, narratives, and `board.md` have flai-owned front matter, so only the body may change; design and docs files are editable whole, front matter validated by the check; generated files (`wip/agents/index.md`, `design/issues/summary.md`), threads, issue front matter, and everything under `wip/archive` are not editable here. The hash is SHA-256 of the file as the editor loaded it. Refusal restores the file. The commit is path-limited, on the main checkout's branch, authored by the git identity in the environment, with a trailer naming the dashboard; `dashboard.autocommit: false` in `system-flow.yaml` leaves the change uncommitted. Update `design/system/flaiover-dashboard.md` (routes and API), `flai-cli.md`, and `project-manifest.md`.

## Done when
- ADR-0023 is accepted and indexed
- The three design documents describe the command, the API, and the manifest key
- `flai check --strict` is clean

## Notes
