---
id: T-1250
type: task
nature: feature
title: The dashboard shows a project whose flai it refused as unsigned where it shows a flai that lacks a required method
status: backlog
parent: S-0237
owner: alex
created: 2026-10-07T23:00:46Z
updated: 2026-10-07T23:00:46Z
transitions: []
stream: S-0237
tags: [dashboard]
touches: [flaiover/src/routes/api/agent/+server.ts, flaiover/src/routes/api/agent/agent.test.ts, flaiover/src/lib/hostflai.svelte.ts, flaiover/src/lib/components/HostFlaiBanner.svelte, flaiover/src/lib/components/HostFlaiBanner.svelte.test.ts]
after: [T-1248]
---
# T-1250 The dashboard shows a project whose flai it refused as unsigned where it shows a flai that lacks a required method

## Work

The first half of criterion 3. It waits for T-1248, which keeps the refusal per project in the registry; this task shows it.

- `flaiover/src/routes/api/agent/+server.ts` already turns `status.missing` into an `error` naming the flai version and what to do. Do the same for a project whose flai was refused: "this dashboard refused flai <version> on the host: unsigned", with the reason and what to do. Say that a release of flai is needed, or that `dashboard.allow_unsigned` will allow a development build once it exists (S-0239). Do this even though no connection is up.
- `flaiover/src/lib/hostflai.svelte.ts` carries the refusal in `HostFlaiStatus` so that the header badge and the banner do not read it as only "not connected".
- `flaiover/src/lib/components/HostFlaiBanner.svelte` shows the refusal on every page of that project, as it shows a flai that lacks methods.
- Tests: `flaiover/src/routes/api/agent/agent.test.ts` for the refused project's status; `HostFlaiBanner.svelte.test.ts` for the banner's text.

## Done when

- A project whose flai the hub refused shows the reason and the flai's version in the banner and the header badge, in place of "not connected".
- The tests pass with `flai test` on the changed paths.

## Notes

Drafted by the planner. `HostFlai.svelte`, the header badge, reads `status.error` already; if it needs its own words for the refusal, the story's agent widens the touches.
