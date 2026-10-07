---
id: T-1263
type: task
nature: feature
title: "The dashboard's connection list says \"unsigned allowed\" per project, with the side"
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:30Z
updated: 2026-10-07T23:09:30Z
transitions: []
stream: S-0239
tags: [dashboard]
touches: [flaiover/src/routes/api/projects/+server.ts, flaiover/src/routes/api/projects/projects.test.ts, flaiover/src/lib/project.svelte.ts, flaiover/src/lib/components/ProjectSwitcher.svelte, flaiover/src/lib/components/ProjectSwitcher.svelte.test.ts]
after: [T-1260]
---
# T-1263 The dashboard's connection list says "unsigned allowed" per project, with the side

## Work

The connection list of criterion 3. It waits for T-1260, which records the unsigned side per connection.

- `flaiover/src/routes/api/projects/+server.ts`: each project's glance carries the unsigned side or sides of its connection when it was allowed, beside S-0237's signed or refused, and nothing for a signed pair.
- `flaiover/src/lib/project.svelte.ts`: the client's project list keeps the field, so the switcher and T-1264's banner read it from one place.
- `flaiover/src/lib/components/ProjectSwitcher.svelte`: a connected project whose connection was allowed unsigned says "unsigned allowed" with the side, in place of "signed" or "unsigned".
- Tests: `projects.test.ts` for a signed pair, an unsigned flai allowed, and an unsigned dashboard allowed; `ProjectSwitcher.svelte.test.ts` for their labels.

## Done when

- `/api/projects` and the project switcher say "unsigned allowed" per project with the side, and say nothing of it for a signed pair.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. If S-0237 put the per-project connection state elsewhere (`flaiover/src/lib/settings.ts`, say), the story's agent widens the touches to it.
