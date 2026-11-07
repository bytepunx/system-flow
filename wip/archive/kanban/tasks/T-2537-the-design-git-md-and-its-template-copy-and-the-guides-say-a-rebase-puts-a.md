---
id: T-2537
type: task
nature: remediation
title: The design, git.md and its template copy, and the guides say a rebase puts a story's changelog lines back under the top Unreleased
status: done
parent: S-0570
owner: alex
created: 2026-11-07T03:48:57Z
updated: 2026-11-07T04:54:41Z
transitions:
  - to: ready
    at: 2026-11-07T04:16:38Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:53:58Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:54:41Z
    by: agent-S-0570
stream: S-0570
tags: [docs, git]
touches: [design/system/workflow.md, design/conventions/git.md, template/root/design/conventions/git.md, docs/users/flai.md, docs/contributors/template.md, template/CHANGELOG.md]
after: [T-2535, T-2536]
usage:
  source: log
  seconds: 43
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 6
      output: 1954
      cache_read: 861636
      cache_write: 13016
      cost: 0.3099
---
# T-2537 The design, git.md and its template copy, and the guides say a rebase puts a story's changelog lines back under the top Unreleased

## Work

Waits for T-2535 and T-2536: it describes what they built, the sync's and accept's step and the batch's.

Say, linking T-2524's ADR, that after a rebase flai writes each `CHANGELOG.md` the branch changed again as main has it, with the branch's own lines under the top `## Unreleased`, and commits it on the branch:

- **`design/system/workflow.md`** § Hub files: what S-0550 changed, the changelogs bullet.
- **`design/conventions/git.md`**, the `CHANGELOG.md merge=union` rule. The rule is above the marker, in the baseline, so the template copy `template/root/design/conventions/git.md` changes in this task too, as CLAUDE.md asks.
- **`docs/users/flai.md`** § Story branches: what `flai stream sync` prints for a rewritten changelog.
- **`docs/contributors/template.md`** § Versioning.
- **`template/CHANGELOG.md`**: a line under `## Unreleased` for the template's git.md change.

## Done when

- Each file says it, and none says union alone keeps a story's line in place.
- `flai test` passes on these paths (the markdown lint).

## Notes

- Drafted by the planner. If git.md needs no change, because its rule still reads true, leave it, its template copy, and `template/CHANGELOG.md` as they are, and drop them from the touches.
