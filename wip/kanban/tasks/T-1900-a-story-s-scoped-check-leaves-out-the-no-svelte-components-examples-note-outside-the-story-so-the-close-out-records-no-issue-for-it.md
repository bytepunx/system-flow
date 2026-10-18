---
id: T-1900
type: task
nature: improvement
title: A story's scoped check leaves out the no-svelte components.examples note outside the story, so the close-out records no issue for it
status: in-progress
parent: S-0446
owner: alex
created: 2026-10-18T02:38:23Z
updated: 2026-10-18T03:01:02Z
transitions:
  - to: ready
    at: 2026-10-18T03:01:01Z
    by: agent-S-0446
  - to: in-progress
    at: 2026-10-18T03:01:02Z
    by: agent-S-0446
stream: S-0446
tags: [flai, check, close-out]
touches: [flai/internal/check/scope.go, flai/internal/check/scope_test.go, flai/cmd/check.go, flai/cmd/check_test.go, design/system/flai-cli.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, docs/users/flai.md, flai/cmd/check_stats_test.go]
after: [T-1899]
---
# T-1900 A story's scoped check leaves out the no-svelte components.examples note outside the story, so the close-out records no issue for it

## Work

This is the second half of the fix. With T-1899 alone, the close-out still records the one note outside the story, so I-0161 would be bumped on every close-out.

- In `ScopeToStory` in `flai/internal/check/scope.go`, leave out the advisory `components.examples` note T-1899 makes when it is outside the story. The note is outside when nothing the story claims lies under the examples folder it names. Drop it as `board.wip-limit` and `threads.archived` are dropped. A note on a folder the story does change stays in, as a note. A `components.examples` warning on a file stays recorded as it is today.
- Test it in `flai/internal/check/scope_test.go` both ways: outside it is left out, and inside it is kept.
- Reproduce I-0161 in `flai/cmd/check_test.go`. Use the fake node from `flai/internal/check/components_test.go`, with no svelte, and examples files the story does not touch. `flai check --story S-nnnn --record-issues` must record no issue and leave `design/issues/summary.md` as it was.
- Name the new left-out kind in these places, with the ADR T-1899 wrote:
  - the `--story` help text and `recordOutside`'s comment in `flai/cmd/check.go`
  - `ScopeToStory`'s comment
  - the `flai check` row in `design/system/flai-cli.md`
  - the `--story` and `--record-issues` paragraphs of `docs/users/flai.md`
  - the close-out rule in `design/conventions/work-management.md` and its template copy `template/root/design/conventions/work-management.md`, with an entry in `template/CHANGELOG.md`.

It waits for T-1899: the scope leaves out the note T-1899 makes, at the folder T-1899 puts it, and both tasks change `docs/users/flai.md`.

## Done when

A scoped check with no svelte in flaiover, on a story that claims nothing under `design/components/flaiover`, reports no `components.examples` finding, and `--record-issues` records no issue. The tests above pass with `flai test` on the changed paths. `flai check --strict` is clean.

## Notes

The fake node lives in a `_test.go` file of package `check`, so `flai/cmd` cannot import it. Copy the smallest part the test needs, or set `nodeRunner` through a test hook.
