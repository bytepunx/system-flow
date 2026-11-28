---
id: T-3639
type: task
nature: improvement
title: flai digest absorb and compact lint the entry they rewrite and refuse a rewrite that brings a markdown finding, as new and edit do
status: backlog
parent: S-0853
owner: alex
created: 2026-11-27T23:58:52Z
updated: 2026-11-27T23:58:59Z
transitions: []
stream: S-0853
tags: [flai]
touches: [flai/internal/digest/decay.go, flai/internal/digest/decay_test.go, flai/internal/digest/write.go]
after: [T-3631]
---
# T-3639 flai digest absorb and compact lint the entry they rewrite and refuse a rewrite that brings a markdown finding, as new and edit do

## Work

Build rule 3 of T-3631's ADR. `Absorb` and `Compact` in `flai/internal/digest/decay.go` write through `rewrite`, which runs `Validate` and `introduced` but no markdown lint. ADR-0318 made `New` and `Edit` lint through `lint` in `flai/internal/digest/write.go`, which calls `workitem.Repo.LintGuard`.

In `rewrite`, after the write and before the put-back, lint the file with the same before-and-after form `Edit` uses. A finding the rewrite brings puts the file back and refuses with a `*RefusedError` naming the rule `markdown.MDnnn`, its line, and markdownlint's message. A finding the entry had before does not refuse. A project with no markdownlint configuration is not linted. Move or share the helper in `write.go` if `rewrite` needs it in another form.

It waits for T-3631: it builds the rule that ADR records. It shares no path with T-3636 or T-3638 and runs beside them.

## Done when

- A test in `flai/internal/digest/decay_test.go` shows a compact whose heading fallback would bring a finding refused, with the file as it was.
- A test shows an absorb of an entry that already had a finding is not refused for it.
- `flai test` passes on the changed paths.

## Notes
