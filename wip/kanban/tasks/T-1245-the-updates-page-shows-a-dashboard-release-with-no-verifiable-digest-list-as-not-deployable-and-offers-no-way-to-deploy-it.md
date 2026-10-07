---
id: T-1245
type: task
nature: feature
title: The Updates page shows a dashboard release with no verifiable digest list as not deployable and offers no way to deploy it
status: backlog
parent: S-0236
owner: alex
created: 2026-10-07T22:54:20Z
updated: 2026-10-07T22:54:20Z
transitions: []
stream: S-0236
tags: [dashboard]
touches: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, docs/users/flaiover.md]
after: [T-1242]
---
# T-1245 The Updates page shows a dashboard release with no verifiable digest list as not deployable and offers no way to deploy it

## Work

In `flaiover/src/lib/components/HostPanel.svelte`, add the field T-1242 put in `flai dashboard versions --json` to the `Release` type. Show a release that is not deployable as such, with the reason, such as "no signed digest list", and give it no deploy button. A release the field is missing from, from a flai older than this story, is shown as today.

Cover a deployable release, one that is not, and one without the field in `HostPanel.svelte.test.ts`. Say in `docs/users/flaiover.md`, where it describes the Updates page's dashboard releases, that a release with no digest list this flai can verify, such as one from before signing, cannot be deployed from the page.

Waits for T-1242, which names the field.

## Done when

- A release marked not deployable shows that it is not, and offers no deploy action.
- A deployable release and a release without the field behave as before.
- `docs/users/flaiover.md` says so.
- `flai test` passes on the component and its test.

## Notes
