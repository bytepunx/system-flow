---
id: T-2355
type: task
nature: feature
title: flai reads a role file's start prompt and renders it over a run's fields, and flai check reports one that does not parse or names an unknown field
status: backlog
parent: S-0554
owner: alex
created: 2026-10-30T20:48:42Z
updated: 2026-10-30T20:48:42Z
transitions: []
stream: S-0554
tags: [cli]
touches: [flai/internal/roles/start.go, flai/internal/roles/start_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go]
after: [T-2354]
---
# T-2355 flai reads a role file's start prompt and renders it over a run's fields, and flai check reports one that does not parse or names an unknown field

## Work

Parse the start prompt where the ADR puts it, in `flai/internal/roles/roles.go`, with a project's additions as the ADR says. Add `flai/internal/roles/start.go`, which renders it over a run's fields. A start prompt that does not parse, or names a field the ADR does not list, is `roles.file` from the validation in `roles.go`. A role file with no start prompt keeps today's behaviour for a defined role.

Second layer: it waits for the ADR. It shares no path with the conditions task.

## Done when

`start_test.go` renders each field the ADR lists, with an empty and a set value, and appends a project's additions. `roles_test.go` covers a start prompt that does not parse and one with an unknown field. The existing role tests pass unchanged.

## Notes

Drafted by the planner of E-0024 on 2026-10-30.
