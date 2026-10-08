---
id: T-1379
type: task
nature: improvement
title: The definition of done and the close-out rule say the branch contains main but for commits that change only wip the branch does not change
status: done
parent: S-0347
owner: alex
created: 2026-10-08T08:44:10Z
updated: 2026-10-08T10:36:05Z
transitions:
  - to: ready
    at: 2026-10-08T10:31:19Z
    by: agent-S-0347
  - to: in-progress
    at: 2026-10-08T10:31:19Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:36:05Z
    by: agent-S-0347
stream: S-0347
tags: [template, docs]
touches: [design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md]
after: [T-1375]
usage:
  source: log
  seconds: 286
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 15
      output: 5089
      cache_read: 831681
      cache_write: 26779
      cost: 0.4434
---
# T-1379 The definition of done and the close-out rule say the branch contains main but for commits that change only wip the branch does not change

## Work

`work-management.md` says the branch is synced "so that it contains the main branch", and that the close-out "checks again that the branch contains the main branch". After T-1375's ADR, a branch that lacks only commits changing `wip/` paths it does not change is done too.

- Change those two rules in `template/root/design/conventions/work-management.md`, the baseline, and copy the same words into `design/conventions/work-management.md` above its marker, as a template change lands here in the same story. Link the ADR.
- Add an entry to `template/CHANGELOG.md` under the unreleased version for the convention and for the template's `scripts/close-out.sh`, whose last check T-1382 changes, naming the flai version that first has `flai verify --sync-only`.

Waits for T-1375, whose ADR it cites. Runs beside T-1377: they share no path.

## Done when

- Both copies of `work-management.md` carry the same baseline text, naming the `wip` exception and linking the ADR.
- `template/CHANGELOG.md` has the entry.
- `flai test template` and the markdown lint pass on the changed files.

## Notes

Drafted by the planner for S-0347.
