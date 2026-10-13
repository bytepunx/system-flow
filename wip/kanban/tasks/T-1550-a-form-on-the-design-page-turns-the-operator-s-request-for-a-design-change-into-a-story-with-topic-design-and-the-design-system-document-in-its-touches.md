---
id: T-1550
type: task
nature: feature
title: A form on the Design page turns the operator's request for a design change into a story with topic design and the design system document in its touches
status: backlog
parent: S-0384
owner: alex
created: 2026-10-13T17:48:07Z
updated: 2026-10-13T17:48:07Z
transitions: []
stream: S-0384
tags: [dashboard]
touches: [flaiover/src/lib/components/design/ProposeChange.svelte, flaiover/src/lib/components/design/ProposeChange.svelte.test.ts]
---
# T-1550 A form on the Design page turns the operator's request for a design change into a story with topic design and the design system document in its touches

## Work

`ProposeChange.svelte` takes a title, the operator's words, and an epic chosen from the open ones, and creates the story through the dashboard's existing story creation (`POST /api/items`, as `/new` does), with topic `design`, the goal from the words, a criterion that the change previews on the Design page, and the design system document in its touches; it links the new story when made and shows flai's refusal otherwise. Component tests cover the request sent, success, and a refusal. It waits for nothing.

## Done when

- `ProposeChange.svelte.test.ts` passes, and the story it sends passes `flai check --strict` when made by hand with the same fields.

## Notes

- Layer 1. The operator makes the story, so it is not a draft; it goes to ready by the operator's own move.
