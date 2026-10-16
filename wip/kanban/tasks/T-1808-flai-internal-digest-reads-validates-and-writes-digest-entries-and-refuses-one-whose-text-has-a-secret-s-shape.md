---
id: T-1808
type: task
nature: feature
title: flai/internal/digest reads, validates, and writes digest entries, and refuses one whose text has a secret's shape
status: backlog
parent: S-0428
owner: alex
created: 2026-10-16T06:03:05Z
updated: 2026-10-16T06:03:05Z
transitions: []
stream: S-0428
tags: [cli]
touches: [flai/internal/digest/digest.go, flai/internal/digest/digest_test.go, flai/internal/digest/write.go, flai/internal/digest/write_test.go, flai/internal/digest/secret.go, flai/internal/digest/secret_test.go]
after: [T-1807]
---
# T-1808 flai/internal/digest reads, validates, and writes digest entries, and refuses one whose text has a secret's shape

## Work

Add the package `flai/internal/digest`, on the pattern of `flai/internal/analysis`:

- `digest.go`: `Folder`, `Dir(m)`, the `Entry` type with the ADR's fields, `List` and `Get` over the folder, and `Validate`, which returns each problem with its field's line: a missing field, an unknown kind or state, a source that names no work item, thread, or issue (archived items included), an `absorbed_by` that names no file.
- `write.go`: `New`, which allocates the next `D-nnnn`, slugs the title, writes the file in the checkout it runs in, validates it, and removes it again on a failure; `Edit`, which changes the body, topics, paths, sources, and `asked_again` with the same check.
- `secret.go`: `Secrets(text)`, the shapes of keys and tokens (a PEM block, an `sk-`, `ghp_`, `github_pat_`, `xox`, or `AKIA` prefix, a long run of base64 or hex after `key`, `token`, or `secret` and `=` or `:`), each with its line; `New` and `Edit` refuse an entry that has one.

Waits for T-1807, whose ADR fixes the fields and the rules.

## Done when

The package's tests cover a valid entry, each validation failure, the next ID after a gap, a refused write that leaves no file, and each secret shape with a near miss that passes; `flai test flai/internal/digest` passes.

## Notes

`New` writes where it runs, so S-0429's seed run in the story's worktree commits its entries on the branch, and the curator, in the main checkout, writes there.
