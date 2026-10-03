---
id: T-0725
type: task
nature: improvement
title: The work-management convention and close-out.sh say the per-task cycle, here and in the template
status: done
parent: S-0197
owner: arobson
created: 2026-10-03T00:40:10Z
updated: 2026-10-03T00:59:21Z
transitions:
  - to: ready
    at: 2026-10-03T00:40:34Z
    by: agent-S-0197
  - to: in-progress
    at: 2026-10-03T00:52:36Z
    by: agent-S-0197
  - to: done
    at: 2026-10-03T00:59:21Z
    by: agent-S-0197
stream: S-0197
tags: []
touches: [design/conventions/git.md, design/conventions/work-management.md, design/conventions/delegation.md, template/root/design/conventions/git.md, template/root/design/conventions/work-management.md, template/root/design/conventions/delegation.md, scripts/close-out.sh, template/root/scripts/close-out.sh, template/CHANGELOG.md, template/template.yaml]
usage:
  source: log
  seconds: 405
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 97
      output: 37539
      cache_read: 3287781
      cache_write: 130574
      cost: 2.2133
---
# T-0725 The work-management convention and close-out.sh say the per-task cycle, here and in the template

## Work

In `design/conventions/work-management.md` and the template's copy, say the per-task cycle (commit, `flai stream sync`, resolve, test, commit fixes; the order is settled on TH-0072) and the sync before review, and in `git.md` both copies if the designer allows the baseline edit. In `scripts/close-out.sh` and the template's copy, say it in the header and stop with "run flai stream sync" when the branch does not contain the main branch, checking before the tests when the worktree is clean and after the commit otherwise. Add a template changelog entry. Waits for nothing: no path in common with the other tasks.

## Done when

- [ ] Both copies of the work-management convention describe the per-task cycle and the sync before review
- [ ] Both close-out scripts refuse a branch behind the main branch, saying to sync, and say the cycle in their header
- [ ] `template/CHANGELOG.md` records it

## Notes
