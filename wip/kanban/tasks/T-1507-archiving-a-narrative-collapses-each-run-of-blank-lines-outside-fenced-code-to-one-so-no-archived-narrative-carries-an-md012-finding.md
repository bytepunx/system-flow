---
id: T-1507
type: task
nature: remediation
title: Archiving a narrative collapses each run of blank lines outside fenced code to one, so no archived narrative carries an MD012 finding
status: backlog
parent: S-0376
owner: alex
created: 2026-10-13T16:44:09Z
updated: 2026-10-13T16:44:09Z
transitions: []
stream: S-0376
tags: [flai, narrative]
touches: [flai/internal/workitem/archive.go, flai/internal/workitem/archive_test.go]
---
# T-1507 Archiving a narrative collapses each run of blank lines outside fenced code to one, so no archived narrative carries an MD012 finding

## Work

- An archived narrative has no story to own it, so a markdown finding on it is outside every story's close-out and is recorded in an issue at each one (ADR-0123 keeps it a note). Only a hand commit on main clears it. I-0139's was cleared by `6526d727`, which broke the file's front matter, and `8a30de84`, which repaired it.
- `## Open questions` is edited by hand ([agent-narrative.md](../../../design/system/agent-narrative.md) § Narrative structure). A hand edit can leave two blank lines there after the story's own close-out has passed, and T-1506 does not reach that.
- `Archive` in `flai/internal/workitem/archive.go` moves each narrative with `moveFile`, an `os.Rename`, and reads nothing. When the source is a narrative, read it, collapse each run of two or more blank lines outside a fenced code block to one, write it to the destination, and remove the source. Keep the rename for every other file, and keep the refusal when the destination exists.
- Leave every other byte as it was, front matter included.
- Add `flai/internal/workitem/archive_test.go`: a narrative with two blank lines under `## Open questions`, archived, reads `## Open questions\n\n## Log`, and two blank lines inside a fenced code block are kept.
- It waits for no task: T-1506 changes other files, so the two run together.

## Done when

- [ ] A narrative archived with a run of blank lines outside fenced code has one blank line there, checked by a test that fails before the change.
- [ ] Blank lines inside fenced code and every other line are unchanged.
- [ ] `flai test` on the changed paths passes.

## Notes
