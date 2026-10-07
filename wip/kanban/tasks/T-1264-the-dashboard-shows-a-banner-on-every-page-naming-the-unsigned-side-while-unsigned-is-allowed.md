---
id: T-1264
type: task
nature: feature
title: The dashboard shows a banner on every page naming the unsigned side while unsigned is allowed
status: backlog
parent: S-0239
owner: alex
created: 2026-10-07T23:09:38Z
updated: 2026-10-07T23:09:38Z
transitions: []
stream: S-0239
tags: [dashboard]
touches: [flaiover/src/lib/components/UnsignedBanner.svelte, flaiover/src/lib/components/UnsignedBanner.svelte.test.ts, flaiover/src/routes/+layout.svelte]
after: [T-1263]
---
# T-1264 The dashboard shows a banner on every page naming the unsigned side while unsigned is allowed

## Work

The banner of criterion 3 and its test in criterion 5. It waits for T-1263, whose project list carries the unsigned side the banner reads.

- `flaiover/src/lib/components/UnsignedBanner.svelte`, new: for the project picked, when its connection was allowed unsigned, a banner saying "unsigned allowed" and naming the side, "this dashboard is a development build", "the flai on the host is a development build", or both, with a link to the operator documentation's section on the setting. Nothing for a signed pair. It is not dismissible: allowed is never silent (ADR-0070).
- `flaiover/src/routes/+layout.svelte`: show it on every page but `/login`, beside `HostFlaiBanner`.
- Tests in `UnsignedBanner.svelte.test.ts`: each side named, both sides, and nothing for a signed pair or a project with no connection.

## Done when

- Every page but the login shows the banner, naming the unsigned side, while the project's connection was allowed unsigned, and none for a signed pair.
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. The component's name is the planner's guess; the story keeps `flaiover/src/lib/components` as a folder touch for it.
