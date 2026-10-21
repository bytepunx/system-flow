---
id: T-1975
type: task
nature: remediation
title: scripts/check.sh, scripts/lint-md.sh, and the integration checks scope to MERGE_QUEUE_STORIES in a batch, recording no issue
status: backlog
parent: S-0464
owner: alex
created: 2026-10-21T07:30:32Z
updated: 2026-10-21T07:30:32Z
transitions: []
stream: S-0464
tags: [flai, template]
touches: [scripts/check.sh, scripts/lint-md.sh, scripts/README.md, template/root/scripts/check.sh, template/root/scripts/README.md, template/CHANGELOG.md, flai/internal/check/check_test.go, flai/internal/mdlint/repo_test.go]
after: [T-1973]
---
# T-1975 scripts/check.sh, scripts/lint-md.sh, and the integration checks scope to MERGE_QUEUE_STORIES in a batch, recording no issue

## Work

With T-1974 the batch names its stories in `MERGE_QUEUE_STORIES`; this task makes the checks inside smoke and integration read it, so a batch fails only on what its stories bring.

- `scripts/check.sh`: with `CLOSE_OUT_STORY` set, as now (`--story` and `--record-issues`); else with `MERGE_QUEUE_STORIES` set, `--story "$MERGE_QUEUE_STORIES"` and no `--record-issues`; else unscoped, as in CI and `make`. Pass the list whole, never split an unquoted variable (`tooling.md`).
- `template/root/scripts/check.sh`: the same, since a project made from the template runs the same queue. Add an entry to `template/CHANGELOG.md`.
- `scripts/lint-md.sh`: with `MERGE_QUEUE_STORIES` set and no files, take the scoped branch it takes for `CLOSE_OUT_STORY`, linting `wip/` only where the batch changes it against the main branch, so a bad line on main's `wip/` fails no batch either (I-0117's rule).
- `TestMonorepoIsClean` in `flai/internal/check/check_test.go` and the repository lint in `flai/internal/mdlint/repo_test.go` read `MERGE_QUEUE_STORIES` when `CLOSE_OUT_STORY` is unset, and scope to the union of the stories' changes.
- Say so in `scripts/README.md` and `template/root/scripts/README.md`.

It waits for T-1973, since `scripts/check.sh` passes several stories to `flai check --story`. It does not wait for T-1974: the variable's name is fixed by the plan, and it is tested here by setting it by hand.

## Done when

- `env -u CLOSE_OUT_STORY MERGE_QUEUE_STORIES=<two open stories> scripts/check.sh` passes in a checkout with a warning outside both, such as an unarchived cancelled story, and `env -u CLOSE_OUT_STORY -u MERGE_QUEUE_STORIES scripts/check.sh` fails on it: I-0174's reproduction, now passing.
- That run records no issue: `git status` shows no change under `design/issues`.
- `CLOSE_OUT_STORY` set still runs `--story` and `--record-issues`, as before.
- `flai test` passes on the changed paths, the template tier among them.

## Notes
