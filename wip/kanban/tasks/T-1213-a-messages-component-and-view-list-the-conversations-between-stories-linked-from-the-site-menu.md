---
id: T-1213
type: task
nature: feature
title: A Messages component and view list the conversations between stories, linked from the site menu
status: backlog
parent: S-0336
owner: alex
created: 2026-10-07T20:17:00Z
updated: 2026-10-08T04:31:21Z
transitions: []
stream: S-0336
tags: [flaiover]
touches: [flaiover/src/lib/components/Messages.svelte, flaiover/src/lib/components/Messages.svelte.test.ts, flaiover/src/routes/messages/+page.svelte, flaiover/src/routes/messages/messages.svelte.test.ts, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts]
after: [T-1212]
---
# T-1213 A Messages component and view list the conversations between stories, linked from the site menu

## Work

Build the view. It waits for T-1212, whose route it fetches.

- `Messages.svelte` renders conversations: the two stories linked, the `about` paths, which side it awaits, the age, and the entries, open first.
- `src/routes/messages/+page.svelte` lists every open conversation, with a toggle for the closed ones, and refreshes on project events as the threads view does.
- Add Messages to `sitemenu.ts`.

## Done when

- Component and page tests cover open and closed conversations, an escalated one, and an empty project.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
