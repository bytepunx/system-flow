---
id: T-2190
type: task
nature: remediation
title: flai keeps one list of the five ADR statuses and flai check reports adr.status on a status outside it
status: backlog
parent: S-0512
owner: alex
created: 2026-10-28T08:27:39Z
updated: 2026-10-28T08:27:39Z
transitions: []
stream: S-0512
tags: [cli]
touches: [flai/internal/adr/adr.go, flai/internal/check/check.go, flai/internal/check/check_test.go]
---
# T-2190 flai keeps one list of the five ADR statuses and flai check reports adr.status on a status outside it

## Work

Add one exported list of every status an ADR may carry, `proposed`, `accepted`, `superseded`, `deprecated`, and `rejected`, in a package that `check`, `docedit`, and `adr` can all import without a cycle: `adr` imports `check` and `docedit`, so the list cannot live in `adr`; `check` is the first candidate, with `adr.Statuses` (what `flai adr new` takes) kept in `adr.go` and documented as a subset of it, or built from it.

In `check.go`, where the `adr.*` findings are made (`adr.id`, `adr.front-matter`, `adr.decision`), add the warning `adr.status` on an ADR whose `status` is not in the list, on the line of the `status` key, naming the status found and the five allowed. The template `0000-template.md` carries `status: template`: leave it out as the other ADR findings do, or pass over the status `template` on that file alone, whichever the existing `adrName` handling makes simpler.

In `check_test.go`, add a test that reproduces I-0208: a repository with an ADR whose status is `bogus` gets `adr.status`, and one whose status is `rejected` gets none.

This task waits for nothing.

## Done when

- `flai check` on a repository with an ADR of status `bogus` reports `adr.status` with the status found and the five allowed; an ADR of status `rejected` reports nothing.
- One exported list holds the five statuses, and `flai adr new --status` still refuses anything but `proposed` and `accepted`.
- The Go tests of `flai/internal/check` pass, and `flai test` on the changed paths passes.

## Notes

Drafted by planner-S-0512. The import graph was read on 2026-10-28: `adr` imports `check`, `docedit`, `execx`, `storygit`, `template`, `topics`, and `workitem`; `docedit` imports `workitem` and `topics`; `check` imports neither `adr` nor `docedit`.
