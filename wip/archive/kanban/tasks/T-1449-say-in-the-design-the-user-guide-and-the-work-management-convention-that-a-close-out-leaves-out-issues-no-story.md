---
id: T-1449
type: task
nature: improvement
title: Say in the design, the user guide, and the work-management convention that a close-out leaves out issues.no-story
status: done
parent: S-0365
owner: alex
created: 2026-10-12T14:56:39Z
updated: 2026-10-12T15:04:31Z
transitions:
  - to: ready
    at: 2026-10-12T14:58:37Z
    by: agent-S-0365
  - to: in-progress
    at: 2026-10-12T14:59:32Z
    by: agent-S-0365
  - to: done
    at: 2026-10-12T15:04:31Z
    by: agent-S-0365
stream: S-0365
tags: [docs, conventions, template]
touches: [design/system/flai-cli.md, design/system/continuous-improvement.md, docs/users/flai.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-1447]
usage:
  source: log
  seconds: 299
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 19
      output: 4980
      cache_read: 1014897
      cache_write: 45284
      cost: 0.6083
---
# T-1449 Say in the design, the user guide, and the work-management convention that a close-out leaves out issues.no-story

## Work

Where the documents list what a close-out leaves out (`item.archive`, `board.wip-limit`, `narrative.state`, markdown on another open story's narrative) or does not record (`wip.overlap`), add `issues.no-story` outside the story, citing T-1447's ADR:

- `design/system/flai-cli.md`, the `flai check` row and its section on `--story` and `--record-issues`.
- `design/system/continuous-improvement.md`, the paragraph on the findings a close-out meets outside the story (S-0249, ADR-0085). Beside its `issues.no-story` bullet, say that the main checkout's check and the review page report it, and a story's close-out does not.
- `docs/users/flai.md`, § Check the repository, the `--story` and `--record-issues` paragraphs.
- `design/conventions/work-management.md`, the close-out rule. Make the same change to the baseline sentence in `template/root/design/conventions/work-management.md`, with an entry in `template/CHANGELOG.md`.

It waits for T-1447 for the ADR's number. It shares no path with T-1448, so the two run together.

## Done when

- Each document names `issues.no-story` among what a close-out leaves out, linking the ADR.
- The project's and the template's `work-management.md` baselines match, and `template/CHANGELOG.md` records the change.
- The markdown lint passes on the changed files.

## Notes
