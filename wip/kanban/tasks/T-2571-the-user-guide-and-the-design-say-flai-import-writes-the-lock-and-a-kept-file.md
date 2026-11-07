---
id: T-2571
type: task
nature: remediation
title: The user guide and the design say flai import writes the lock, and a kept file counts as the project's own
status: backlog
parent: S-0590
owner: alex
created: 2026-11-07T10:02:43Z
updated: 2026-11-07T10:02:43Z
transitions: []
stream: S-0590
tags: [docs, import]
touches: [docs/users/flai.md, design/system/template.md, design/system/flai-cli.md]
after: [T-2570]
---
# T-2571 The user guide and the design say flai import writes the lock, and a kept file counts as the project's own

## Work

Three documents still say that `flai import` writes no lock. Each describes the behaviour T-2570 changes, so this task waits for T-2570 and writes what it built.

- `docs/users/flai.md` § import: the sentence "`import` writes no lock" becomes "`import` writes `system-flow.lock.yaml` as `flai new` does". Say that a file the import kept, because the repository already had it, counts as yours. The next `flai upgrade` reports it as a conflict and does not replace it. In the `--commit` and `--formula` paragraphs, say that the commit holds the lock.
- `design/system/template.md`, the paragraph on variables: an imported project records its values in the lock it writes, no longer "at its first `flai upgrade` or `flai upgrade --relock`". Give the reason for rendered hashes, as T-2570's code comment gives it, and name I-0302.
- `design/system/flai-cli.md` § Import flow: step 6 writes the lock beside `system-flow.yaml`. The `--commit` bullet that lists the staged paths names the lock. Also check the `flai import` row of the command table.

Bump `updated` on each document you change.

## Done when

- [ ] No document under `docs/` or `design/system/` says that `flai import` writes no lock
- [ ] `docs/users/flai.md` and `design/system/template.md` say how the lock treats a file the import kept
- [ ] `flai test docs/users/flai.md design/system/template.md design/system/flai-cli.md` passes

## Notes
