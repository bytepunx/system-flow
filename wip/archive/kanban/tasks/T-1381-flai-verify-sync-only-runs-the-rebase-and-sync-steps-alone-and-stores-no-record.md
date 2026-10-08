---
id: T-1381
type: task
nature: improvement
title: flai verify --sync-only runs the rebase and sync steps alone and stores no record
status: done
parent: S-0347
owner: alex
created: 2026-10-08T08:44:19Z
updated: 2026-10-08T10:41:26Z
transitions:
  - to: ready
    at: 2026-10-08T10:36:29Z
    by: agent-S-0347
  - to: in-progress
    at: 2026-10-08T10:36:29Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:41:26Z
    by: agent-S-0347
stream: S-0347
tags: [cli, docs]
touches: [flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1377]
usage:
  source: log
  seconds: 297
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 16153
      cache_read: 2639870
      cache_write: 84999
      cost: 1.4075
---
# T-1381 flai verify --sync-only runs the rebase and sync steps alone and stores no record

## Work

The close-out's last check runs after its commit, so it must ask again, at the new head, whether the branch contains main under the rule T-1377 built. Give it a flai command instead of `git merge-base --is-ancestor`, so the rule lives in one place:

- Add `--sync-only` to `flai verify` in `flai/cmd/verify.go`: it runs the `rebase` and `sync` steps alone, prints them and any note as a full run does, stores no record (the last full run's record stays what `--last` and the review page show), and exits 0, 1, or 2 as a full run does. Refuse it with `--last` and with `--record-issues`.
- Test in `verify_test.go`: a branch that lacks only a `wip`-only commit passes; one that lacks a commit outside `wip` exits 1 naming `flai stream sync`; the stored record is unchanged after either.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference`, and describe `--sync-only` and the `sync` step's note in `docs/users/flai.md` under "Verify a story before review" and in its table of steps.

Waits for T-1377, whose `sync` step it runs.

## Done when

- `flai verify S-nnnn --sync-only` behaves as above, and the tests in `verify_test.go` show it.
- `docs/users/flai.md` and `docs/users/flai-reference.md` describe the flag and the note.
- `flai test flai/cmd` passes.

## Notes

Drafted by the planner for S-0347. S-0341 adds `--fresh` to the same command; keep the two flags apart.
