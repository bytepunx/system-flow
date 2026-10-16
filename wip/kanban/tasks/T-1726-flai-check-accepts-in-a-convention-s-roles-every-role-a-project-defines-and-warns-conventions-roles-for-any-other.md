---
id: T-1726
type: task
nature: feature
title: flai check accepts in a convention's roles every role a project defines and warns conventions.roles for any other
status: backlog
parent: S-0415
owner: alex
created: 2026-10-16T05:43:42Z
updated: 2026-10-16T05:43:42Z
transitions: []
stream: S-0415
tags: [cli]
touches: [flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go]
---
# T-1726 flai check accepts in a convention's roles every role a project defines and warns conventions.roles for any other

## Work

`conventions.Roles` and `StrategicRoles` are fixed lists today. Make the `conventions.roles` rule take the roles the project defines, from `roles.Load`, beside `story`, `explore`, and `verify`. Keep the warning's text for an unknown role and list the known ones in it. Test a convention that lists a defined role, and one that lists a role no file defines.

First layer: shares no path with T-1725.

## Done when

`go test ./internal/conventions/...` passes. `flai check --strict` here reports nothing new.

## Notes
