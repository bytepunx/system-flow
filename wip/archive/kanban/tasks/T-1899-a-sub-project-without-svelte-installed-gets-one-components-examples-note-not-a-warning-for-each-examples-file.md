---
id: T-1899
type: task
nature: improvement
title: A sub-project without svelte installed gets one components.examples note, not a warning for each examples file
status: done
parent: S-0446
owner: alex
created: 2026-10-18T02:37:51Z
updated: 2026-10-18T03:00:41Z
transitions:
  - to: ready
    at: 2026-10-18T02:57:23Z
    by: agent-S-0446
  - to: in-progress
    at: 2026-10-18T02:57:24Z
    by: agent-S-0446
  - to: done
    at: 2026-10-18T03:00:41Z
    by: agent-S-0446
stream: S-0446
tags: [flai, check, components]
touches: [flai/internal/check/check.go, flai/internal/check/components_test.go, design/adrs, design/adrs/README.md, design/system/component-library.md, docs/users/flai.md]
usage:
  source: log
  seconds: 197
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7932
      cache_read: 1490990
      cache_write: 68368
      cost: 0.8706
---
# T-1899 A sub-project without svelte installed gets one components.examples note, not a warning for each examples file

## Work

This task builds the first half of the fix the plan's thread on S-0446 proposes. Confirm it in the narrative's `## Decisions` before building it, or change it there and say why.

The cause, from both instances in I-0161 (S-0392 and S-0393): a story worktree has no `flaiover/node_modules`. `components.ReadProps` then returns `ErrSvelteMissing`. `exampleProps` in `flai/internal/check/check.go` reports it as a warning on every examples file, about 48 in all.

- In `exampleProps`, when `errors.Is(err, components.ErrSvelteMissing)`, report one advisory note per sub-project with `c.advise`. Say how many examples files were not checked, and keep the reader's error with its `npm install` advice. A host without node already gets a note the same way. Every other read error stays a warning on each file.
- Put the note at the sub-project's examples folder, such as `design/components/flaiover`, not at a file. The scoped leave-out in the next task can then tell it apart from a file's finding.
- In `flai/internal/check/components_test.go`, test that a sub-project without svelte gets one note at its folder, and that `--strict` passes over it. Test that a component whose props fail to read for another reason still gets its warning at line 1. Change `TestComponentExamplesWhosePropsCannotBeRead` to match.
- Write an ADR with `flai adr new` that refines ADR-0164 §5. It records two things: the one note for a sub-project without svelte, and the next task's decision that a story's scoped check leaves that note out outside the story, as ADR-0133 does for `board.wip-limit`. List it in `design/adrs/README.md`.
- Update `design/system/component-library.md` § on `components.examples`, and the note paragraph of `docs/users/flai.md` on `components.examples`. Add the new note to the `--strict` advisory list in `docs/users/flai.md`.

First layer: it waits for no task.

## Done when

A check with no svelte in a sub-project gives one advisory `components.examples` note at that sub-project's examples folder and no warning on its files. The tests above pass with `flai test` on the changed paths. The ADR is accepted and indexed, and `flai check --strict` reports nothing on it.

## Notes

The ADR's file name is not known until `flai adr new` makes it, so the task touches the shared `design/adrs` folder.
