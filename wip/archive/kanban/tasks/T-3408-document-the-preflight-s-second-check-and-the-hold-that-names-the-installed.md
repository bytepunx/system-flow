---
id: T-3408
type: task
nature: remediation
title: Document the preflight's second check and the hold that names the installed flai and the release to install
status: done
parent: S-0797
owner: alex
created: 2026-11-16T01:44:52Z
updated: 2026-11-16T03:34:03Z
transitions:
  - to: ready
    at: 2026-11-16T03:14:36Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-16T03:14:36Z
    by: agent-S-0797
  - to: done
    at: 2026-11-16T03:34:03Z
    by: agent-S-0797
stream: S-0797
tags: [docs, merge-queue]
touches: [design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [T-3403]
usage:
  source: log
  seconds: 1167
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 10599
      cache_read: 2422947
      cache_write: 89166
      cost: 1.2964
---
# T-3408 Document the preflight's second check and the hold that names the installed flai and the release to install

## Work

Write ADR-0309's rules where the preflight is described, linking the ADR, in the same story as the code (`documentation.md`).

- `design/system/flai-cli.md` § The merge queue pass: under **The preflight.**, the second check with `flai.command` when the running flai is behind main's flai sub-project. Under **The thread to the operator.**, the `flai_behind` cause and its words.
- `docs/users/flai.md` § When the main checkout holds the queue: how to tell a hold on files from a hold on an older flai, and the remedy for the second, which is to publish, install the release, and restart `flai serve`.
- `docs/operators/index.md`, the paragraph **When your main checkout holds the queue.**: the same, in a sentence or two, linking the users guide.

It waits for T-3403, whose decision it describes. It changes no file the code tasks change, so it runs beside them.

## Done when

- The three sections state the second check, the cause, and the remedy, each linking ADR-0309.
- None of them still tells the operator only to fix the named files.
- The markdown lint and `flai check --strict` pass on the three files.

## Notes
