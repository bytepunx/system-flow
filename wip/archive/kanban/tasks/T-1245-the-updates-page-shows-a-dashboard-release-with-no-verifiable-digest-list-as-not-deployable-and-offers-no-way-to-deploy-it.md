---
id: T-1245
type: task
nature: feature
title: The Updates page shows a dashboard release with no verifiable digest list as not deployable and offers no way to deploy it
status: done
parent: S-0236
owner: alex
created: 2026-10-07T22:54:20Z
updated: 2026-10-12T18:58:12Z
transitions:
  - to: ready
    at: 2026-10-12T18:37:05Z
    by: agent-S-0236
  - to: in-progress
    at: 2026-10-12T18:37:06Z
    by: agent-S-0236
  - to: done
    at: 2026-10-12T18:58:12Z
    by: agent-S-0236
stream: S-0236
tags: [dashboard]
touches: [flaiover/src/lib/components/HostPanel.svelte, flaiover/src/lib/components/HostPanel.svelte.test.ts, docs/users/flaiover.md]
after: [T-1242]
usage:
  source: log
  seconds: 1266
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 12
      output: 6670
      cache_read: 906528
      cache_write: 38707
      cost: 0.5312
---
# T-1245 The Updates page shows a dashboard release with no verifiable digest list as not deployable and offers no way to deploy it

## Work

In `flaiover/src/lib/components/HostPanel.svelte`, add the fields T-1242 put in `flai dashboard versions --json`, `unverifiable?`, `signature?`, and `reason?`, to the `Release` type, as `HostProcesses.svelte` has them for flai releases. Show a release marked `unverifiable` as `not deployable: <reason>`, with a fallback reason when none is given, and give it no deploy button, as `HostProcesses.svelte` does for a flai release it cannot install. A release without the fields, from a flai older than this story, is shown as today.

Cover a deployable release, one that is not, and one without the fields in `HostPanel.svelte.test.ts`, as `HostProcesses.svelte.test.ts` does. Say in `docs/users/flaiover.md`, where it describes the Updates page's dashboard releases, that a release with no digest list this flai can verify, such as one from before signing, cannot be deployed from the page.

Waits for T-1242, which names the fields.

## Done when

- A release marked unverifiable shows that it is not deployable, and offers no deploy action.
- A deployable release and a release without the fields behave as before.
- `docs/users/flaiover.md` says so.
- `flai test` passes on the component and its test.

## Notes
