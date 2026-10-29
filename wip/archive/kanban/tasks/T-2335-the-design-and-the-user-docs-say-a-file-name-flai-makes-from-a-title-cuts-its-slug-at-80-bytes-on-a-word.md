---
id: T-2335
type: task
nature: improvement
title: The design and the user docs say a file name flai makes from a title cuts its slug at 80 bytes on a word
status: done
parent: S-0542
owner: alex
created: 2026-10-29T15:24:07Z
updated: 2026-10-29T15:31:49Z
transitions:
  - to: ready
    at: 2026-10-29T15:28:50Z
    by: agent-S-0542
  - to: in-progress
    at: 2026-10-29T15:28:51Z
    by: agent-S-0542
  - to: done
    at: 2026-10-29T15:31:49Z
    by: agent-S-0542
stream: S-0542
tags: [docs]
touches: [design/system/work-hierarchy.md, design/system/documentation-standard.md, docs/users/flai.md]
after: [T-2333]
usage:
  source: log
  seconds: 178
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 24
      output: 5943
      cache_read: 885366
      cache_write: 48691
      cost: 0.6141
---
# T-2335 The design and the user docs say a file name flai makes from a title cuts its slug at 80 bytes on a word

## Work

Say where each document names files by ID and slug that the slug is the title's, cut at 80 bytes on a word, so that a name and a tool's temporary file beside it stay under the 255-byte limit (I-0231):

- `design/system/work-hierarchy.md`, its `File name:` line;
- `design/system/documentation-standard.md` § Naming, the work item bullet, and bump its `updated`;
- `docs/users/flai.md`, where it says how `flai story new` and `flai issue new` name their files, and the `MS-0004-<slug>.md` sentence for conversations.

Waits for T-2333, so that it describes the cut T-2333 settles; it runs beside T-2334, sharing no file.

## Done when

- The three documents state the cut and its reason, and link no line number.
- The markdown lint passes on them.

## Notes
