---
id: T-1424
type: task
nature: improvement
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out board.wip-limit
status: backlog
parent: S-0348
owner: alex
created: 2026-10-08T08:59:07Z
updated: 2026-10-08T08:59:07Z
transitions: []
stream: S-0348
tags: [docs, conventions, template]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-1422]
---
# T-1424 Say in the design, the user guide, and the work-management convention that a close-out leaves out board.wip-limit

## Work

Where the documents list what a close-out leaves out (`item.archive`, `narrative.state`, markdown on another open story's narrative) or does not record (`wip.overlap`), add `board.wip-limit`, citing the ADR of T-1422:

- `design/system/flai-cli.md`, the `flai check` row and its section on `--story` and `--record-issues`.
- `design/system/continuous-improvement.md`, the paragraph on the findings a close-out meets outside the story (S-0249, ADR-0085).
- `docs/users/flai.md`, § Check the repository, the `--story` and `--record-issues` paragraphs.
- `design/conventions/work-management.md`, the close-out rule, and the same baseline sentence in `template/root/design/conventions/work-management.md`, with an entry in `template/CHANGELOG.md`.

It waits for T-1422 for the ADR's number. It shares no path with T-1423, so the two run together.

## Done when

- Each document names `board.wip-limit` among what a close-out leaves out, linking the ADR.
- The project's and the template's `work-management.md` baselines match, and `template/CHANGELOG.md` records the change.
- The markdown lint passes on the changed files.

## Notes
