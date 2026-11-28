---
id: T-3636
type: task
nature: improvement
title: flai accept leaves out of its commits a shared-set file the story does not change whose lint the fix cannot fix, names it, and no longer refuses for it
status: backlog
parent: S-0853
owner: alex
created: 2026-11-27T23:58:39Z
updated: 2026-11-27T23:58:39Z
transitions: []
stream: S-0853
tags: [flai]
touches: [flai/cmd/accept_markdown.go, flai/cmd/accept.go, flai/cmd/accept_fix_test.go, flai/cmd/accept_uncommitted_test.go, flai/cmd/accept_orchestrator_dirty_test.go, flai/cmd/accept_batch_dirty.go, flai/cmd/accept_batch_dirty_test.go, flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/hostapi/writes_test.go, docs/users/flai-reference.md]
after: [T-3631]
---
# T-3636 flai accept leaves out of its commits a shared-set file the story does not change whose lint the fix cannot fix, names it, and no longer refuses for it

## Work

Build rule 1 of T-3631's ADR. Today `flai accept` lints the main checkout's whole shared set at preflight, through `previewedToLint` and `refuseUnfixedUncommitted` in `flai/cmd/accept_markdown.go`, called from `flai/cmd/accept.go`. It refuses for any file the fix leaves failing, so one digest entry stops every acceptance (I-0411).

- Split what the fix left by whether the story's change, `main...story/S-nnnn`, touches the file. Refuse only for the files it touches, with today's text. Leave the rest out of the shared-set commit (`commitSharedBeforeMerge`) and the acceptance commit (`commitAcceptance`), uncommitted. Log each at warn, and name them in the acceptance's output and `--json`.
- The same holds for `flai accept --batch`, through `leftUncommitted` in `flai/cmd/accept_batch_dirty.go`, and for `--yes`.
- In `flai/internal/preview/accept.go`, list the files the acceptance will leave out for lint under a field of their own, so the dry run and the dashboard's review page show them. The preview does not lint today. Run the check there only if it costs no more than the preflight does, and say in the narrative which you chose.
- Correct `flai accept --help`, which says the fix runs over "what the main checkout has not committed": it runs over the shared set. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Update the refusal text the tests assert in `flai/cmd/accept_fix_test.go` and `flai/internal/hostapi/writes_test.go` where it changes.

It waits for T-3631: it builds the rule that ADR records.

## Done when

- A test in `flai/cmd/accept_fix_test.go` reproduces I-0411: an uncommitted `design/digest/` entry with an MD028 finding that the story does not change. The acceptance lands, leaves the entry uncommitted, and names it.
- A test shows a shared-set file the story's change touches still refuses, naming its file, line, and error.
- A test in `flai/cmd/accept_batch_dirty_test.go` shows the batch landing treats the file the same way.
- A test in `flai/internal/preview/accept_test.go` shows the preview's new field.
- `flai test` passes on the changed paths.

## Notes
