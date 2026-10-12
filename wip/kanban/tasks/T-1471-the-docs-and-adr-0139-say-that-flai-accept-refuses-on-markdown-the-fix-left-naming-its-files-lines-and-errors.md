---
id: T-1471
type: task
nature: improvement
title: The docs and ADR-0139 say that flai accept refuses on markdown the fix left, naming its files, lines, and errors
status: in-progress
parent: S-0367
owner: alex
created: 2026-10-12T17:07:47Z
updated: 2026-10-12T17:19:43Z
transitions:
  - to: ready
    at: 2026-10-12T17:19:43Z
    by: agent-S-0367
  - to: in-progress
    at: 2026-10-12T17:19:43Z
    by: agent-S-0367
stream: S-0367
tags: []
touches: [design/system/flai-cli.md, docs/users/flai.md, design/adrs/0139-a-test-tier-s-fix-command-rewrites-the-files-a-story-changes-before-flai-task.md, design/system/project-manifest.md, docs/operators/settings.md]
after: [T-1459]
---
# T-1471 The docs and ADR-0139 say that flai accept refuses on markdown the fix left, naming its files, lines, and errors

## Work

- `design/system/flai-cli.md` (`flai accept`), `docs/users/flai.md`, `design/system/project-manifest.md`, `docs/operators/settings.md`, and ADR-0139 (proposed, so it is still S-0367's to change) say that `flai accept` refuses when a tier with a `fix_command` still fails after the fix: at preflight over the main checkout's uncommitted files, and after the rebase over the files the branch changes. The refusal names the files, lines, and errors, and the story stays in review.
- Waits for T-1459, whose behaviour it describes.

## Done when

- [ ] The markdown lint passes on the changed docs, and they match what `flai accept` does.

## Notes
