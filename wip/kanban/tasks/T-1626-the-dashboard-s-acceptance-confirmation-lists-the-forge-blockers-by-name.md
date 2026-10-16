---
id: T-1626
type: task
nature: feature
title: The dashboard's acceptance confirmation lists the forge blockers by name
status: backlog
parent: S-0395
owner: alex
created: 2026-10-16T04:38:18Z
updated: 2026-10-16T04:38:18Z
transitions: []
stream: S-0395
tags: [dashboard]
touches: ["flaiover/src/routes/review/[id]/+page.svelte"]
after: [T-1623]
---
# T-1626 The dashboard's acceptance confirmation lists the forge blockers by name

## Work

- Check whether the review page's acceptance confirmation lists the preview's blockers generically. Where it does not show T-1623's kinds by name, add them: each check or review linking its URL from the kept state.
- Add a test beside the page's other tests. If the page has none, test the component that renders the blockers.
- Waits for T-1623: it shows that task's blockers.

## Done when

The confirmation shows each forge blocker by name with its link, and the flaiover tests pass.

## Notes
