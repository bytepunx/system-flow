---
id: T-0844
type: task
nature: improvement
title: The close-out scripts, here and in the template, scope flai check to the story and record what lies outside it
status: done
parent: S-0249
owner: alex
created: 2026-10-05T00:18:19Z
updated: 2026-10-05T00:57:15Z
transitions:
  - to: ready
    at: 2026-10-05T00:51:44Z
    by: agent-S-0249
  - to: in-progress
    at: 2026-10-05T00:51:44Z
    by: agent-S-0249
  - to: done
    at: 2026-10-05T00:57:15Z
    by: agent-S-0249
stream: S-0249
tags: [flai, template, scripts]
touches: [scripts/close-out.sh, scripts/check.sh, template/root/scripts/close-out.sh, template/root/scripts/check.sh, flai/internal/check/check_test.go, flai/internal/storygit, flai/cmd/check.go, scripts/README.md, template/root/scripts/README.md]
after: [T-0841]
usage:
  source: log
  seconds: 331
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 25966
      cache_read: 3386887
      cache_write: 102018
      cost: 1.7963
---
# T-0844 The close-out scripts, here and in the template, scope flai check to the story and record what lies outside it

## Work

`scripts/close-out.sh S-nnnn` reaches `flai check --strict` in two ways. When the story touches `flai`, it gets there through `scripts/flai-test.sh`, then `scripts/smoke.sh`, then `scripts/check.sh`. Otherwise it calls `scripts/check.sh` directly. Both paths stopped on findings outside the story in I-0057.

Have the close-out export the story, for instance as `CLOSE_OUT_STORY`. Have `scripts/check.sh` add `--story "$CLOSE_OUT_STORY" --record-issues` when it is set, so both paths are scoped and `smoke.sh` needs no change. Run outside a close-out, as CI and `make` do, `check.sh` stays unscoped.

Make the same change to `template/root/scripts/close-out.sh` and `template/root/scripts/check.sh`. T-0845 records the template change in `template/CHANGELOG.md`.

In the integration tier, `TestMonorepoIsClean` fails on any error in the main checkout. That stopped S-0200 and S-0203 on main's `issues.duplicate-id`. When `CLOSE_OUT_STORY` is set, have it run scoped to that story, so that errors outside the story are logged and do not fail.

Issue files that `--record-issues` writes land in the story's worktree before the close-out's commit step, which commits them. Check that the commit step picks them up.

It waits for T-0841, whose flag `check.sh` passes, and so, through T-0841, for T-0840. It runs beside T-0845: the two share no path.

## Done when

- In a story worktree whose branch changes nothing that main's findings name, `scripts/close-out.sh S-nnnn -m "…"` runs past the check step. Main's own findings are printed as outside the story, and the issue each was recorded in is printed. Today those findings are `story.unaccepted` on S-0173 and `threads.archived` on TH-0094.
- `scripts/check.sh`, run with no close-out, still fails on those findings, as CI does.
- `scripts/template-test.sh` passes, and a project rendered from the template has the scoped `check.sh`.
- `scripts/flai-test.sh` passes.

## Notes
