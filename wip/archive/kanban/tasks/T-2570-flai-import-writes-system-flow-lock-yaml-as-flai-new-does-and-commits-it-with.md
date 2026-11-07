---
id: T-2570
type: task
nature: remediation
title: flai import writes system-flow.lock.yaml as flai new does and commits it with the import, so the import formula's check finds no layout.lock
status: done
parent: S-0590
owner: alex
created: 2026-11-07T10:02:21Z
updated: 2026-11-07T10:14:06Z
transitions:
  - to: ready
    at: 2026-11-07T10:04:37Z
    by: agent-S-0590
  - to: in-progress
    at: 2026-11-07T10:04:39Z
    by: agent-S-0590
  - to: done
    at: 2026-11-07T10:14:06Z
    by: agent-S-0590
stream: S-0590
tags: [flai, import]
touches: [flai/cmd/import.go, flai/cmd/import_commit.go, flai/cmd/import_test.go, flai/cmd/import_commit_test.go, docs/users/flai-reference.md]
usage:
  source: log
  seconds: 567
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 41
      output: 5555
      cache_read: 4927546
      cache_write: 68673
      cost: 1.6243
---
# T-2570 flai import writes system-flow.lock.yaml as flai new does and commits it with the import, so the import formula's check finds no layout.lock

## Work

I-0302's instance, S-0564: the board's Import ran `flai import --formula`. It recorded template 1.1.24 in `system-flow.yaml` and wrote no `system-flow.lock.yaml`. The instance's first step, `flai check --strict`, then exited 1 on `layout.lock` (`flai/internal/check/check.go`). The cause is in `runImport` in `flai/cmd/import.go`: step 2 calls `template.Render`, but no `lock.Save` follows it. `flai new` in `flai/cmd/new.go` saves the lock right after its render.

Write the test first, and see it fail against the current code. In `flai/cmd/import_commit_test.go`, run `flai import --formula` on a git repository. Then check three things:

- `system-flow.lock.yaml` exists.
- The import commit holds it, so `git status --porcelain` is clean after the commit.
- `flai check --strict` in the imported project reports no `layout.lock`.

Then fix it:

- After the render in `runImport`, save the lock as `flai new` does. Use `lock.Template` with the source's repo and ref, the manifest's version, and the time applied. Take `Files` from `res.Hashes`, `Topics` from `upgrade.TopicsOf(root, res.Written)`, and `Vars` from `lockVars(vars)`. Write nothing under `--dry-run`, or on any path that returns before the render.
- `res.Hashes` holds the rendered hash of every template path, including the files the import kept because the repository already had them. So a kept file that differs from the template counts as the project's own. The next `flai upgrade` reports it as a conflict and does not replace it. That is the safe reading. Do not hash the current files, as `upgrade.Relock` does: that would let an upgrade overwrite a repository's own file as if it were unchanged baseline. Write this reasoning in a code comment.
- In `commitImportPaths` in `flai/cmd/import_commit.go`, add `system-flow.lock.yaml` to the paths. Both `--commit` and `--formula` then commit it.
- In `TestImportAppliesTheNewestRelease` in `flai/cmd/import_test.go`, change the comment that says an import writes no lock. Add a check that the lock records the version applied.
- In `flai import`'s `Long` help, say that it writes `system-flow.lock.yaml` beside `system-flow.yaml`. Run `make flai-reference` so that `docs/users/flai-reference.md` matches the help.

This task waits for no other.

## Done when

- [ ] A test in `flai/cmd/import_commit_test.go` fails without the fix and passes with it: after `flai import --formula`, the lock is committed and `flai check --strict` reports no `layout.lock`
- [ ] `flai import --dry-run` writes no lock
- [ ] `docs/users/flai-reference.md` matches the import help
- [ ] `flai test flai/cmd` passes

## Notes
